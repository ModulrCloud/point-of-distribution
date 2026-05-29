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
	Block                 *external_structs.CoreBlock                   `json:"block"`
	Afp                   *external_structs.AggregatedFinalizationProof `json:"afp"`
	AggregatedHeightProof *external_structs.AggregatedHeightProof       `json:"aggregatedHeightProof,omitempty"`
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

type AggregatedHeightProofStoreRequest struct {
	Route string                                 `json:"route"`
	Proof external_structs.AggregatedHeightProof `json:"proof"`
}

type AggregatedHeightProofGetRequest struct {
	Route          string `json:"route"`
	AbsoluteHeight int    `json:"absoluteHeight"`
}

type AggregatedHeightProofGetResponse struct {
	Proof *external_structs.AggregatedHeightProof `json:"proof"`
}

type AggregatedEpochRotationProofStoreRequest struct {
	Route string                                        `json:"route"`
	Proof external_structs.AggregatedEpochRotationProof `json:"proof"`
}

type AggregatedEpochRotationProofGetRequest struct {
	Route   string `json:"route"`
	EpochId int    `json:"epochId"`
}

type AggregatedEpochRotationProofGetResponse struct {
	Proof *external_structs.AggregatedEpochRotationProof `json:"proof"`
}

type AggregatedEpochAnnouncementProofStoreRequest struct {
	Route string                                            `json:"route"`
	Proof external_structs.AggregatedEpochAnnouncementProof `json:"proof"`
}

type AggregatedEpochAnnouncementProofGetRequest struct {
	Route       string `json:"route"`
	NextEpochId int    `json:"nextEpochId"`
}

type AggregatedEpochAnnouncementProofGetResponse struct {
	Proof *external_structs.AggregatedEpochAnnouncementProof `json:"proof"`
}

type BlockByHeightRequest struct {
	Route          string `json:"route"`
	AbsoluteHeight int    `json:"absoluteHeight"`
}

type BlockByHeightResponse struct {
	Block                 *external_structs.CoreBlock             `json:"block"`
	AggregatedHeightProof *external_structs.AggregatedHeightProof `json:"aggregatedHeightProof"`
}

type AggregatedAnchorEpochAckProofStoreRequest struct {
	Route string                                         `json:"route"`
	Proof external_structs.AggregatedAnchorEpochAckProof `json:"proof"`
}

type AggregatedAnchorEpochAckProofGetRequest struct {
	Route   string `json:"route"`
	EpochId int    `json:"epochId"`
}

type AggregatedAnchorEpochAckProofGetResponse struct {
	Proof *external_structs.AggregatedAnchorEpochAckProof `json:"proof"`
}
