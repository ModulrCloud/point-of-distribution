package websocket

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	anchorBlocks "github.com/modulrcloud/modulr-anchors-core/block_pack"
	anchorsStructs "github.com/modulrcloud/modulr-anchors-core/structures"
	coreBlocks "github.com/modulrcloud/modulr-core/block_pack"
	coreStructs "github.com/modulrcloud/modulr-core/structures"

	"github.com/modulrcloud/point-of-distribution/databases"

	"github.com/lxzan/gws"
)

type lockEntry struct {
	mu       sync.Mutex
	refCount int32
}

type lockManager struct {
	locks     sync.Map
	semaphore chan struct{}
}

func newLockManager(limit int) *lockManager {
	if limit <= 0 {
		limit = 100
	}
	return &lockManager{semaphore: make(chan struct{}, limit)}
}

func handleGetBlockWithAfp(req BlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	key := req.BlockId
	if key == "" || stores == nil || stores.CoreBlocksData == nil {
		return
	}
	var resp BlockWithAfpResponse
	nextBlockId := nextBlockId(key)
	if nextBlockId != "" {
		if afpBytes, err := stores.CoreBlocksData.Get([]byte("AFP:"+nextBlockId), nil); err == nil {
			var afp coreStructs.AggregatedFinalizationProof
			if err := json.Unmarshal(afpBytes, &afp); err == nil {
				resp.Afp = &afp
			}
		}
	}
	if blockBytes, err := stores.CoreBlocksData.Get([]byte(key), nil); err == nil {
		var block coreBlocks.Block
		if err := json.Unmarshal(blockBytes, &block); err == nil {
			resp.Block = &block
		}
	}
	if resp.Block != nil {
		if data, err := json.Marshal(resp); err == nil {
			connection.WriteMessage(gws.OpcodeText, data)
		}
	}
}

func handleGetAnchorBlockWithAfp(req AnchorBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	key := req.BlockId
	if key == "" || stores == nil || stores.AnchorsCoreBlocksData == nil {
		return
	}
	var resp AnchorBlockWithAfpResponse
	nextBlockId := nextBlockId(key)
	if nextBlockId != "" {
		if afpBytes, err := stores.AnchorsCoreBlocksData.Get([]byte("AFP:"+nextBlockId), nil); err == nil {
			var afp anchorsStructs.AggregatedFinalizationProof
			if err := json.Unmarshal(afpBytes, &afp); err == nil {
				resp.Afp = &afp
			}
		}
	}
	if blockBytes, err := stores.AnchorsCoreBlocksData.Get([]byte(key), nil); err == nil {
		var block anchorBlocks.Block
		if err := json.Unmarshal(blockBytes, &block); err == nil {
			resp.Block = &block
		}
	}
	if resp.Block != nil {
		if data, err := json.Marshal(resp); err == nil {
			connection.WriteMessage(gws.OpcodeText, data)
		}
	}
}

func handleAcceptAggregatedLeaderFinalizationProof(req WsAggregatedLeaderFinalizationProofStoreRequest, connection *gws.Conn, stores *databases.Stores) {
	key := composeAggregatedLeaderFinalizationProofKey(req.Proof.EpochIndex, req.Proof.Leader)
	if key == "" || stores == nil || stores.AggregatedLeaderFinalizationProofs == nil {
		return
	}

	if proofBytes, err := json.Marshal(req.Proof); err == nil {
		if err := stores.AggregatedLeaderFinalizationProofs.Put([]byte(key), proofBytes, nil); err == nil {
			acknowledge(connection)
		}
	}
}

func handleGetAggregatedLeaderFinalizationProof(req WsAggregatedLeaderFinalizationProofRequest, connection *gws.Conn, stores *databases.Stores) {
	key := composeAggregatedLeaderFinalizationProofKey(req.EpochIndex, req.Leader)
	if key == "" || stores == nil || stores.AggregatedLeaderFinalizationProofs == nil {
		return
	}

	var resp WsAggregatedLeaderFinalizationProofResponse
	if proofBytes, err := stores.AggregatedLeaderFinalizationProofs.Get([]byte(key), nil); err == nil {
		var proof coreStructs.AggregatedLeaderFinalizationProof
		if err := json.Unmarshal(proofBytes, &proof); err == nil {
			resp.Proof = &proof
		}
	}
	if data, err := json.Marshal(resp); err == nil {
		connection.WriteMessage(gws.OpcodeText, data)
	}
}

func handleAcceptBlockWithAfp(req AcceptBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores, locks *lockManager) {
	locator, blockKey := blockLocatorFromBlock(req.Block), ""
	if locator != nil {
		blockKey = composeBlockKey(*locator)
	}
	if blockKey == "" || stores == nil || stores.CoreBlocksData == nil {
		return
	}
	withBlockLock(blockKey, locks, func() {
		if !req.Block.VerifySignature() {
			return
		}
		if req.Block.Index > 0 {
			if req.Afp == nil || !validateAfpForBlock(req.Block.Index, req.Block.PrevHash, *req.Afp) {
				return
			}
		}
		if blockBytes, err := json.Marshal(req.Block); err == nil {
			if err := stores.CoreBlocksData.Put([]byte(blockKey), blockBytes, nil); err == nil {
				if req.Afp != nil {
					if afpBytes, err := json.Marshal(req.Afp); err == nil && req.Afp.BlockId != "" {
						_ = stores.CoreBlocksData.Put([]byte("AFP:"+req.Afp.BlockId), afpBytes, nil)
					}
				}
				acknowledge(connection)
			}
		}
	})
}

func handleAcceptAnchorBlockWithAfp(req AcceptAnchorBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores, locks *lockManager) {
	locator, blockKey := blockLocatorFromAnchorBlock(req.Block), ""
	if locator != nil {
		blockKey = composeBlockKey(*locator)
	}
	if blockKey == "" || stores == nil || stores.AnchorsCoreBlocksData == nil {
		return
	}
	withBlockLock(blockKey, locks, func() {
		if !req.Block.VerifySignature() {
			return
		}
		if req.Block.Index > 0 {
			if req.Afp == nil || !validateAnchorAfpForBlock(req.Block.Index, req.Block.PrevHash, *req.Afp) {
				return
			}
		}
		if blockBytes, err := json.Marshal(req.Block); err == nil {
			if err := stores.AnchorsCoreBlocksData.Put([]byte(blockKey), blockBytes, nil); err == nil {
				if req.Afp != nil {
					if afpBytes, err := json.Marshal(req.Afp); err == nil && req.Afp.BlockId != "" {
						_ = stores.AnchorsCoreBlocksData.Put([]byte("AFP:"+req.Afp.BlockId), afpBytes, nil)
					}
				}
				acknowledge(connection)
			}
		}
	})
}

func nextBlockId(blockId string) string {
	parts := strings.Split(blockId, ":")
	if len(parts) != 3 {
		return ""
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil {
		return ""
	}
	parts[2] = strconv.Itoa(idx + 1)
	return strings.Join(parts, ":")
}

func composeBlockKey(locator BlockLocator) string {
	if locator.Creator == "" || locator.Index < 0 {
		return ""
	}
	return strings.Join([]string{strconv.Itoa(locator.EpochIndex), locator.Creator, strconv.Itoa(locator.Index)}, ":")
}

func blockLocatorFromBlock(block coreBlocks.Block) *BlockLocator {
	epochIndex := extractEpochIndex(block.Epoch)
	if epochIndex < 0 {
		return nil
	}
	return &BlockLocator{EpochIndex: epochIndex, Creator: block.Creator, Index: block.Index}
}

func blockLocatorFromAnchorBlock(block anchorBlocks.Block) *BlockLocator {
	epochIndex := extractEpochIndex(block.Epoch)
	if epochIndex < 0 {
		return nil
	}
	return &BlockLocator{EpochIndex: epochIndex, Creator: block.Creator, Index: block.Index}
}

func extractEpochIndex(epoch string) int {
	if epoch == "" {
		return -1
	}
	parts := strings.Split(epoch, "#")
	last := parts[len(parts)-1]
	idx, err := strconv.Atoi(last)
	if err != nil {
		return -1
	}
	return idx
}

func composeAggregatedLeaderFinalizationProofKey(epochIndex int, leader string) string {
	if epochIndex < 0 || leader == "" {
		return ""
	}
	return fmt.Sprintf("ALFP:%d:%s", epochIndex, leader)
}

func validateAfpForBlock(blockIndex int, prevHash string, afp coreStructs.AggregatedFinalizationProof) bool {
	if blockIndex <= 0 {
		return true
	}
	if afp.BlockHash != prevHash {
		return false
	}
	expectedIndex := blockIndex - 1
	parts := strings.Split(afp.BlockId, ":")
	if len(parts) != 3 {
		return false
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil || idx != expectedIndex {
		return false
	}
	return true
}

func validateAnchorAfpForBlock(blockIndex int, prevHash string, afp anchorsStructs.AggregatedFinalizationProof) bool {
	if blockIndex <= 0 {
		return true
	}
	if afp.BlockHash != prevHash {
		return false
	}
	expectedIndex := blockIndex - 1
	parts := strings.Split(afp.BlockId, ":")
	if len(parts) != 3 {
		return false
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil || idx != expectedIndex {
		return false
	}
	return true
}

func (lm *lockManager) withLock(blockKey string, fn func()) {
	if lm == nil || blockKey == "" || fn == nil {
		return
	}
	entry := lm.acquireEntry(blockKey)
	entry.mu.Lock()
	fn()
	entry.mu.Unlock()
	lm.releaseEntry(blockKey, entry)
}

func withBlockLock(blockKey string, lm *lockManager, fn func()) {
	if lm == nil {
		return
	}
	lm.withLock(blockKey, fn)
}

func (lm *lockManager) acquireEntry(blockKey string) *lockEntry {
	if val, ok := lm.locks.Load(blockKey); ok {
		entry := val.(*lockEntry)
		atomic.AddInt32(&entry.refCount, 1)
		return entry
	}

	lm.acquireSlot()
	entry := &lockEntry{refCount: 1}
	actual, loaded := lm.locks.LoadOrStore(blockKey, entry)
	if loaded {
		lm.releaseSlot()
		entry = actual.(*lockEntry)
		atomic.AddInt32(&entry.refCount, 1)
	}

	return entry
}

func (lm *lockManager) releaseEntry(blockKey string, entry *lockEntry) {
	if entry == nil {
		return
	}
	if atomic.AddInt32(&entry.refCount, -1) == 0 {
		lm.locks.Delete(blockKey)
		lm.releaseSlot()
	}
}

func (lm *lockManager) acquireSlot() {
	if lm.semaphore == nil {
		return
	}
	lm.semaphore <- struct{}{}
}

func (lm *lockManager) releaseSlot() {
	if lm.semaphore == nil {
		return
	}
	<-lm.semaphore
}

func acknowledge(connection *gws.Conn) {
	if connection == nil {
		return
	}
	if resp, err := json.Marshal(statusResponse{Status: "OK"}); err == nil {
		connection.WriteMessage(gws.OpcodeText, resp)
	}
}
