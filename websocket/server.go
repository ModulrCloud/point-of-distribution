package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lxzan/gws"
	"github.com/modulrcloud/point-of-distribution/config"
	"github.com/modulrcloud/point-of-distribution/databases"
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
			handleAcceptBlockWithAfp(req, connection, h.stores, h.coreLocks)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_accept_block_with_afp_request"}`))
		}
	case "accept_anchor_block_with_afp":
		var req AcceptAnchorBlockWithAfpRequest
		if err := json.Unmarshal(message.Bytes(), &req); err == nil {
			handleAcceptAnchorBlockWithAfp(req, connection, h.stores, h.anchorLocks)
		} else {
			connection.WriteMessage(gws.OpcodeText, []byte(`{"error":"invalid_accept_anchor_block_with_afp_request"}`))
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
