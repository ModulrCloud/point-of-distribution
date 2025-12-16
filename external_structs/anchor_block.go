package external_structs

import (
	"encoding/json"
	"fmt"
)

type AnchorBlock struct {
	Creator   string                 `json:"creator"`
	Time      int64                  `json:"time"`
	Epoch     string                 `json:"epoch"`
	ExtraData ExtraDataToAnchorBlock `json:"extraData"`
	Index     int                    `json:"index"`
	PrevHash  string                 `json:"prevHash"`
	Sig       string                 `json:"sig"`
}

type ExtraDataToAnchorBlock struct {
	AggregatedAnchorRotationProofs     []AggregatedAnchorRotationProof     `json:"aggregatedAnchorRotationProofs,omitempty"`
	AggregatedLeaderFinalizationProofs []AggregatedLeaderFinalizationProof `json:"aggregatedLeaderFinalizationProofs,omitempty"`
	Rest                               map[string]string                   `json:"rest,omitempty"`
}

type blockExtraDataAlias struct {
	AggregatedAnchorRotationProofs     []AggregatedAnchorRotationProof     `json:"aggregatedAnchorRotationProofs,omitempty"`
	AggregatedLeaderFinalizationProofs []AggregatedLeaderFinalizationProof `json:"aggregatedLeaderFinalizationProofs,omitempty"`
	Rest                               map[string]string                   `json:"rest,omitempty"`
}

func (extra ExtraDataToAnchorBlock) MarshalJSON() ([]byte, error) {
	if len(extra.AggregatedAnchorRotationProofs) == 0 && len(extra.AggregatedLeaderFinalizationProofs) == 0 {
		if len(extra.Rest) == 0 {
			return []byte("{}"), nil
		}
		return json.Marshal(extra.Rest)
	}
	alias := blockExtraDataAlias(extra)
	if alias.Rest == nil {
		alias.Rest = map[string]string{}
	}
	return json.Marshal(alias)
}

func (extra *ExtraDataToAnchorBlock) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*extra = ExtraDataToAnchorBlock{}
		return nil
	}
	var alias blockExtraDataAlias
	if err := json.Unmarshal(data, &alias); err == nil && (alias.Rest != nil || alias.AggregatedAnchorRotationProofs != nil || alias.AggregatedLeaderFinalizationProofs != nil) {
		*extra = ExtraDataToAnchorBlock(alias)
		return nil
	}
	var fields map[string]string
	if err := json.Unmarshal(data, &fields); err == nil {
		extra.Rest = fields
		extra.AggregatedAnchorRotationProofs = nil
		extra.AggregatedLeaderFinalizationProofs = nil
		return nil
	}
	return fmt.Errorf("invalid extraData payload")
}
