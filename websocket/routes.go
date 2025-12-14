package websocket

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/lxzan/gws"
	anchorBlocks "github.com/modulrcloud/modulr-anchors-core/block_pack"
	anchorsStructs "github.com/modulrcloud/modulr-anchors-core/structures"
	coreBlocks "github.com/modulrcloud/modulr-core/block_pack"
	"github.com/modulrcloud/modulr-core/structures"
	"github.com/modulrcloud/point-of-distribution/databases"
)

func handleGetBlockWithAfp(req BlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	key := composeBlockKey(req.Locator)
	if key == "" || stores == nil || stores.CoreBlocksData == nil {
		return
	}
	if blockBytes, err := stores.CoreBlocksData.Get([]byte(key), nil); err == nil {
		var block coreBlocks.Block
		if err := json.Unmarshal(blockBytes, &block); err == nil {
			resp := BlockWithAfpResponse{Block: &block}
			if block.Index > 0 {
				prevKey := composeBlockKey(BlockLocator{EpochIndex: req.Locator.EpochIndex, Creator: block.Creator, Index: block.Index - 1})
				if afpBytes, err := stores.CoreBlocksData.Get([]byte("AFP:"+prevKey), nil); err == nil {
					var afp structures.AggregatedFinalizationProof
					if err := json.Unmarshal(afpBytes, &afp); err == nil {
						resp.Afp = &afp
					}
				}
			}
			if data, err := json.Marshal(resp); err == nil {
				connection.WriteMessage(gws.OpcodeText, data)
			}
		}
	}
}

func handleGetAnchorBlockWithAfp(req AnchorBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	key := composeBlockKey(req.Locator)
	if key == "" || stores == nil || stores.AnchorsCoreBlocksData == nil {
		return
	}
	if blockBytes, err := stores.AnchorsCoreBlocksData.Get([]byte(key), nil); err == nil {
		var block anchorBlocks.Block
		if err := json.Unmarshal(blockBytes, &block); err == nil {
			resp := AnchorBlockWithAfpResponse{Block: &block}
			if block.Index > 0 {
				prevKey := composeBlockKey(BlockLocator{EpochIndex: req.Locator.EpochIndex, Creator: block.Creator, Index: block.Index - 1})
				if afpBytes, err := stores.AnchorsCoreBlocksData.Get([]byte("AFP:"+prevKey), nil); err == nil {
					var afp anchorsStructs.AggregatedFinalizationProof
					if err := json.Unmarshal(afpBytes, &afp); err == nil {
						resp.Afp = &afp
					}
				}
			}
			if data, err := json.Marshal(resp); err == nil {
				connection.WriteMessage(gws.OpcodeText, data)
			}
		}
	}
}

func handleAcceptBlockWithAfp(req AcceptBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	locator, blockKey := blockLocatorFromBlock(req.Block), ""
	if locator != nil {
		blockKey = composeBlockKey(*locator)
	}
	if blockKey == "" || stores == nil || stores.CoreBlocksData == nil {
		return
	}
	if !req.Block.VerifySignature() || !validateAfpForBlock(req.Block.Index, req.Block.PrevHash, req.Afp) {
		return
	}
	if blockBytes, err := json.Marshal(req.Block); err == nil {
		if err := stores.CoreBlocksData.Put([]byte(blockKey), blockBytes, nil); err == nil {
			if afpBytes, err := json.Marshal(req.Afp); err == nil && req.Afp.BlockId != "" {
				_ = stores.CoreBlocksData.Put([]byte("AFP:"+req.Afp.BlockId), afpBytes, nil)
			}
			acknowledge(connection)
		}
	}
}

func handleAcceptAnchorBlockWithAfp(req AcceptAnchorBlockWithAfpRequest, connection *gws.Conn, stores *databases.Stores) {
	locator, blockKey := blockLocatorFromAnchorBlock(req.Block), ""
	if locator != nil {
		blockKey = composeBlockKey(*locator)
	}
	if blockKey == "" || stores == nil || stores.AnchorsCoreBlocksData == nil {
		return
	}
	if !req.Block.VerifySignature() || !validateAnchorAfpForBlock(req.Block.Index, req.Block.PrevHash, req.Afp) {
		return
	}
	if blockBytes, err := json.Marshal(req.Block); err == nil {
		if err := stores.AnchorsCoreBlocksData.Put([]byte(blockKey), blockBytes, nil); err == nil {
			if afpBytes, err := json.Marshal(req.Afp); err == nil && req.Afp.BlockId != "" {
				_ = stores.AnchorsCoreBlocksData.Put([]byte("AFP:"+req.Afp.BlockId), afpBytes, nil)
			}
			acknowledge(connection)
		}
	}
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

func validateAfpForBlock(blockIndex int, prevHash string, afp structures.AggregatedFinalizationProof) bool {
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

func acknowledge(connection *gws.Conn) {
	if connection == nil {
		return
	}
	if resp, err := json.Marshal(statusResponse{Status: "OK"}); err == nil {
		connection.WriteMessage(gws.OpcodeText, resp)
	}
}
