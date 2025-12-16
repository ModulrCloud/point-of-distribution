package websocket

import (
	anchorBlocks "github.com/modulrcloud/modulr-anchors-core/block_pack"
	anchorsStructs "github.com/modulrcloud/modulr-anchors-core/structures"
	coreBlocks "github.com/modulrcloud/modulr-core/block_pack"
	coreStructs "github.com/modulrcloud/modulr-core/structures"
)

type incomingMsg struct {
	Route string `json:"route"`
}

type BlockLocator struct {
	EpochIndex int    `json:"epochIndex"`
	Creator    string `json:"creator"`
	Index      int    `json:"index"`
}

type BlockWithAfpRequest struct {
	Route   string `json:"route"`
	BlockId string `json:"blockID"`
}

type BlockWithAfpResponse struct {
	Block *coreBlocks.Block                        `json:"block"`
	Afp   *coreStructs.AggregatedFinalizationProof `json:"afp"`
}

type AnchorBlockWithAfpRequest struct {
	Route   string `json:"route"`
	BlockId string `json:"blockID"`
}

type AnchorBlockWithAfpResponse struct {
	Block *anchorBlocks.Block                         `json:"block"`
	Afp   *anchorsStructs.AggregatedFinalizationProof `json:"afp"`
}

type AcceptBlockWithAfpRequest struct {
	Route string                                  `json:"route"`
	Block coreBlocks.Block                        `json:"block"`
	Afp   coreStructs.AggregatedFinalizationProof `json:"afp"`
}

type AcceptAnchorBlockWithAfpRequest struct {
	Route string                                     `json:"route"`
	Block anchorBlocks.Block                         `json:"block"`
	Afp   anchorsStructs.AggregatedFinalizationProof `json:"afp"`
}

type statusResponse struct {
	Status string `json:"status"`
}
