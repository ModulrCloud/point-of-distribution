package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/modulrcloud/point-of-distribution/config"
	"github.com/modulrcloud/point-of-distribution/databases"

	"github.com/lxzan/gws"
)

type handler struct {
	stores      *databases.Stores
	coreLocks   *lockManager
	anchorLocks *lockManager
}

func (h *handler) OnOpen(conn *gws.Conn) {}

func (h *handler) OnClose(conn *gws.Conn, err error) {}

func (h *handler) OnPing(conn *gws.Conn, payload []byte) {}

func (h *handler) OnPong(conn *gws.Conn, payload []byte) {}

func (h *handler) OnMessage(connection *gws.Conn, message *gws.Message) {
	defer message.Close()

	var incoming incomingMsg
	if err := json.Unmarshal(message.Bytes(), &incoming); err != nil {
		connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_json"}`))
		return
	}

	switch incoming.Route {
	case "get_block_with_afp":
		var req BlockWithAfpRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			handleGetBlockWithAfp(req, connection, h.stores)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_block_with_afp_request"}`))
		}
	case "get_anchor_block_with_afp":
		var req AnchorBlockWithAfpRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			handleGetAnchorBlockWithAfp(req, connection, h.stores)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_anchor_block_with_afp_request"}`))
		}
	case "accept_block_with_afp":
		var req AcceptBlockWithAfpRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			blockKey := ""
			if locator := blockLocatorFromBlock(req.Block); locator != nil {
				blockKey = composeBlockKey(*locator)
			}
			afpBlockId, afpBlockHash := "", ""
			if req.Afp != nil {
				afpBlockId, afpBlockHash = req.Afp.BlockId, req.Afp.BlockHash
			}
			log.Printf("[ws] accept_block_with_afp blockKey=%s epoch=%s creator=%s index=%d prevHash=%s hasAfp=%t afpBlockId=%s afpBlockHash=%s",
				blockKey, req.Block.Epoch, req.Block.Creator, req.Block.Index, req.Block.PrevHash, req.Afp != nil, afpBlockId, afpBlockHash,
			)
			handleAcceptBlockWithAfp(req, connection, h.stores, h.coreLocks)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_accept_block_with_afp_request"}`))
		}
	case "accept_anchor_block_with_afp":
		var req AcceptAnchorBlockWithAfpRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			blockKey := ""
			if locator := blockLocatorFromAnchorBlock(req.Block); locator != nil {
				blockKey = composeBlockKey(*locator)
			}
			afpBlockId, afpBlockHash := "", ""
			if req.Afp != nil {
				afpBlockId, afpBlockHash = req.Afp.BlockId, req.Afp.BlockHash
			}
			log.Printf("[ws] accept_anchor_block_with_afp blockKey=%s epoch=%s creator=%s index=%d prevHash=%s hasAfp=%t afpBlockId=%s afpBlockHash=%s",
				blockKey, req.Block.Epoch, req.Block.Creator, req.Block.Index, req.Block.PrevHash, req.Afp != nil, afpBlockId, afpBlockHash,
			)
			handleAcceptAnchorBlockWithAfp(req, connection, h.stores, h.anchorLocks)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_accept_anchor_block_with_afp_request"}`))
		}
	case "accept_aggregated_leader_finalization_proof":
		var req AggregatedLeaderFinalizationProofStoreRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			vsIndex, vsHash := req.Proof.VotingStat.Index, req.Proof.VotingStat.Hash
			afpBlockId, afpBlockHash := req.Proof.VotingStat.Afp.BlockId, req.Proof.VotingStat.Afp.BlockHash
			log.Printf("[ws] accept_aggregated_leader_finalization_proof epochIndex=%d leader=%s votingStatIndex=%d votingStatHash=%s afpBlockId=%s afpBlockHash=%s",
				req.Proof.EpochIndex, req.Proof.Leader, vsIndex, vsHash, afpBlockId, afpBlockHash,
			)
			handleAcceptAggregatedLeaderFinalizationProof(req, connection, h.stores)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_accept_aggregated_leader_finalization_proof_request"}`))
		}
	case "get_aggregated_leader_finalization_proof":
		var req AggregatedLeaderFinalizationProofRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			handleGetAggregatedLeaderFinalizationProof(req, connection, h.stores)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_get_aggregated_leader_finalization_proof_request"}`))
		}
	default:
		connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"unknown_type"}`))
	}
}

func CreateWebsocketServer(cfg config.Config, stores *databases.Stores) error {
	coreLocks := newLockManager(cfg.MaxConcurrentLocks)
	anchorLocks := newLockManager(cfg.MaxConcurrentLocks)
	upgrader := gws.NewUpgrader(&handler{stores: stores, coreLocks: coreLocks, anchorLocks: anchorLocks}, &gws.ServerOption{
		ParallelEnabled:   true,
		Recovery:          gws.Recovery,
		PermessageDeflate: gws.PermessageDeflate{Enabled: true},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r)
		if err != nil {
			return
		}

		go func() {
			conn.ReadLoop()
		}()
	})

	address := cfg.WSInterface + ":" + strconv.Itoa(cfg.WSPort)
	fmt.Printf("Websocket server is starting at ws://%s ...✅\n", address)
	return http.ListenAndServe(address, nil)
}
