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
	AggregatedHeightProof *external_structs.AggregatedHeightProof       `json:"heightAttestation,omitempty"`
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

type HeightAttestationStoreRequest struct {
	Route string                                 `json:"route"`
	Proof external_structs.AggregatedHeightProof `json:"proof"`
}

type HeightAttestationGetRequest struct {
	Route          string `json:"route"`
	AbsoluteHeight int    `json:"absoluteHeight"`
}

type HeightAttestationGetResponse struct {
	Proof *external_structs.AggregatedHeightProof `json:"proof"`
}

type EpochDataAttestationStoreRequest struct {
	Route string                                        `json:"route"`
	Proof external_structs.AggregatedEpochRotationProof `json:"attestation"`
}

type EpochDataAttestationGetRequest struct {
	Route   string `json:"route"`
	EpochId int    `json:"epochId"`
}

type EpochDataAttestationGetResponse struct {
	Proof *external_structs.AggregatedEpochRotationProof `json:"attestation"`
}

type BlockByHeightRequest struct {
	Route          string `json:"route"`
	AbsoluteHeight int    `json:"absoluteHeight"`
}

type BlockByHeightResponse struct {
	Block                 *external_structs.CoreBlock             `json:"block"`
	AggregatedHeightProof *external_structs.AggregatedHeightProof `json:"heightAttestation"`
}

type AnchorEpochAckStoreRequest struct {
	Route string                                         `json:"route"`
	Proof external_structs.AggregatedAnchorEpochAckProof `json:"proof"`
}

type AnchorEpochAckGetRequest struct {
	Route   string `json:"route"`
	EpochId int    `json:"epochId"`
}

type AnchorEpochAckGetResponse struct {
	Proof *external_structs.AggregatedAnchorEpochAckProof `json:"proof"`
}
