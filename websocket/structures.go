package websocket

import "github.com/modulrcloud/point-of-distribution/external_structs"

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
	Block *external_structs.CoreBlock                   `json:"block"`
	Afp   *external_structs.AggregatedFinalizationProof `json:"afp"`
}

type AnchorBlockWithAfpRequest struct {
	Route   string `json:"route"`
	BlockId string `json:"blockID"`
}

type AnchorBlockWithAfpResponse struct {
	Block *external_structs.AnchorBlock                 `json:"block"`
	Afp   *external_structs.AggregatedFinalizationProof `json:"afp"`
}

type AggregatedLeaderFinalizationProofStoreRequest struct {
	Route string                                             `json:"route"`
	Proof external_structs.AggregatedLeaderFinalizationProof `json:"proof"`
}

type AggregatedLeaderFinalizationProofRequest struct {
	Route      string `json:"route"`
	EpochIndex int    `json:"epochIndex"`
	Leader     string `json:"leader"`
}

type AggregatedLeaderFinalizationProofResponse struct {
	Proof *external_structs.AggregatedLeaderFinalizationProof `json:"proof"`
}

type AcceptBlockWithAfpRequest struct {
	Route string                                        `json:"route"`
	Block external_structs.CoreBlock                    `json:"block"`
	Afp   *external_structs.AggregatedFinalizationProof `json:"afp"`
}

type AcceptAnchorBlockWithAfpRequest struct {
	Route string                                        `json:"route"`
	Block external_structs.AnchorBlock                  `json:"block"`
	Afp   *external_structs.AggregatedFinalizationProof `json:"afp"`
}

type statusResponse struct {
	Status string `json:"status"`
}
