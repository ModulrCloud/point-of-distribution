package main

import (
	"log"

	"github.com/modulrcloud/point-of-distribution/config"
	"github.com/modulrcloud/point-of-distribution/databases"
	"github.com/modulrcloud/point-of-distribution/websocket"
)

func main() {
	cfg := config.Load()
	stores, err := databases.Init(cfg.DataPath)
	if err != nil {
		log.Fatalf("failed to init databases: %v", err)
	}
	defer stores.Close()

	if err := websocket.CreateWebsocketServer(cfg, stores); err != nil {
		log.Fatalf("websocket server error: %v", err)
	}
}
