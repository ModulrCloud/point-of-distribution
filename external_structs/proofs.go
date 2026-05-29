package external_structs

import "encoding/json"

type VotingStat struct {
	Index int                         `json:"index"`
	Hash  string                      `json:"hash"`
	Afp   AggregatedFinalizationProof `json:"afp"`
}

type AggregatedFinalizationProof struct {
	PrevBlockHash string            `json:"prevBlockHash"`
	BlockId       string            `json:"blockId"`
	BlockHash     string            `json:"blockHash"`
	Proofs        map[string]string `json:"proofs"`
}

func (afp *AggregatedFinalizationProof) UnmarshalJSON(data []byte) error {

	type alias AggregatedFinalizationProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	*afp = AggregatedFinalizationProof(aux)

	return nil

}

func (afp AggregatedFinalizationProof) MarshalJSON() ([]byte, error) {

	type alias AggregatedFinalizationProof

	aux := alias(afp)

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	return json.Marshal(aux)
}

type AggregatedAnchorRotationProof struct {
	EpochIndex int               `json:"epochIndex"`
	Anchor     string            `json:"anchor"`
	VotingStat VotingStat        `json:"votingStat"`
	Signatures map[string]string `json:"signatures"`
}

func (aarp *AggregatedAnchorRotationProof) UnmarshalJSON(data []byte) error {

	type alias AggregatedAnchorRotationProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Signatures == nil {
		aux.Signatures = make(map[string]string)
	}

	*aarp = AggregatedAnchorRotationProof(aux)

	return nil

}

func (aarp AggregatedAnchorRotationProof) MarshalJSON() ([]byte, error) {

	type alias AggregatedAnchorRotationProof

	aux := alias(aarp)

	if aux.Signatures == nil {
		aux.Signatures = make(map[string]string)
	}

	return json.Marshal(aux)
}

type AggregatedLeaderFinalizationProof struct {
	EpochIndex int               `json:"epochIndex"`
	Leader     string            `json:"leader"`
	VotingStat VotingStat        `json:"votingStat"`
	Signatures map[string]string `json:"signatures"`
}

func (alfp *AggregatedLeaderFinalizationProof) UnmarshalJSON(data []byte) error {

	type alias AggregatedLeaderFinalizationProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Signatures == nil {
		aux.Signatures = make(map[string]string)
	}

	*alfp = AggregatedLeaderFinalizationProof(aux)

	return nil

}

func (alfp AggregatedLeaderFinalizationProof) MarshalJSON() ([]byte, error) {

	type alias AggregatedLeaderFinalizationProof

	aux := alias(alfp)

	if aux.Signatures == nil {
		aux.Signatures = make(map[string]string)
	}

	return json.Marshal(aux)
}

type AggregatedHeightProof struct {
	AbsoluteHeight int               `json:"absoluteHeight"`
	BlockId        string            `json:"blockId"`
	BlockHash      string            `json:"blockHash"`
	EpochId        int               `json:"epochId"`
	HeightInEpoch  int               `json:"heightInEpoch"`
	Proofs         map[string]string `json:"proofs"`
}

func (ha *AggregatedHeightProof) UnmarshalJSON(data []byte) error {

	type alias AggregatedHeightProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	*ha = AggregatedHeightProof(aux)

	return nil
}

func (ha AggregatedHeightProof) MarshalJSON() ([]byte, error) {

	type alias AggregatedHeightProof

	aux := alias(ha)

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	return json.Marshal(aux)
}

type NextEpochData struct {
	NextEpochHash               string              `json:"nextEpochHash"`
	NextEpochValidatorsRegistry []string            `json:"nextEpochValidatorsRegistry"`
	NextEpochQuorum             []string            `json:"nextEpochQuorum"`
	NextEpochLeadersSequence    []string            `json:"nextEpochLeadersSequence"`
	NextEpochStartTimestamp     uint64              `json:"nextEpochStartTimestamp"`
	DelayedTransactions         []map[string]string `json:"delayedTransactions"`
}

type AggregatedEpochAnnouncementProof struct {
	EpochId       int               `json:"epochId"`
	NextEpochId   int               `json:"nextEpochId"`
	EpochData     NextEpochData     `json:"epochData"`
	EpochDataHash string            `json:"epochDataHash"`
	Proofs        map[string]string `json:"proofs"`
}

func (eda *AggregatedEpochAnnouncementProof) UnmarshalJSON(data []byte) error {
	type alias AggregatedEpochAnnouncementProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	*eda = AggregatedEpochAnnouncementProof(aux)

	return nil
}

func (eda AggregatedEpochAnnouncementProof) MarshalJSON() ([]byte, error) {
	type alias AggregatedEpochAnnouncementProof

	aux := alias(eda)

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	return json.Marshal(aux)
}

type AggregatedAnchorEpochAckProof struct {
	EpochId       int               `json:"epochId"`
	NextEpochId   int               `json:"nextEpochId"`
	EpochDataHash string            `json:"epochDataHash"`
	Proofs        map[string]string `json:"proofs"`
}

func (a *AggregatedAnchorEpochAckProof) UnmarshalJSON(data []byte) error {
	type alias AggregatedAnchorEpochAckProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	*a = AggregatedAnchorEpochAckProof(aux)

	return nil
}

func (a AggregatedAnchorEpochAckProof) MarshalJSON() ([]byte, error) {
	type alias AggregatedAnchorEpochAckProof

	aux := alias(a)

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	return json.Marshal(aux)
}

type AggregatedEpochRotationProof struct {
	EpochId           int               `json:"epochId"`
	NextEpochId       int               `json:"nextEpochId"`
	EpochData         NextEpochData     `json:"epochData"`
	EpochDataHash     string            `json:"epochDataHash"`
	FinishedOnHeight  int64             `json:"finishedOnHeight"`
	FinishedOnBlockId string            `json:"finishedOnBlockId"`
	FinishedOnHash    string            `json:"finishedOnHash"`
	Proofs            map[string]string `json:"proofs"`
}

func (eda *AggregatedEpochRotationProof) UnmarshalJSON(data []byte) error {
	type alias AggregatedEpochRotationProof

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	*eda = AggregatedEpochRotationProof(aux)

	return nil
}

func (eda AggregatedEpochRotationProof) MarshalJSON() ([]byte, error) {
	type alias AggregatedEpochRotationProof

	aux := alias(eda)

	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}

	return json.Marshal(aux)
}
