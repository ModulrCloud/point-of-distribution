package external_structs

import "encoding/json"

type CoreBlock struct {
	Creator      string               `json:"creator"`
	Time         int64                `json:"time"`
	Epoch        string               `json:"epoch"`
	Transactions []Transaction        `json:"transactions"`
	ExtraData    ExtraDataToCoreBlock `json:"extraData"`
	Index        int                  `json:"index"`
	PrevHash     string               `json:"prevHash"`
	Sig          string               `json:"sig"`
}

type ExtraDataToCoreBlock struct {
	DelayedTransactionsBatch DelayedTransactionsBatch `json:"delayedTxsBatch"`
	Rest                     map[string]string        `json:"rest"`
}

func (ed ExtraDataToCoreBlock) MarshalJSON() ([]byte, error) {
	type alias ExtraDataToCoreBlock

	aux := alias(ed)

	// Normalize empty maps to nil so JSON uses `null` instead of {}
	if aux.Rest != nil && len(aux.Rest) == 0 {
		aux.Rest = nil
	}

	return json.Marshal(aux)
}

func (ed *ExtraDataToCoreBlock) UnmarshalJSON(data []byte) error {

	type alias ExtraDataToCoreBlock

	var aux alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Rest == nil {
		aux.Rest = make(map[string]string)
	}

	*ed = ExtraDataToCoreBlock(aux)

	return nil

}

type DelayedTransactionsBatch struct {
	EpochIndex          int                 `json:"epochIndex"`
	DelayedTransactions []map[string]string `json:"delayedTransactions"`
	Proofs              map[string]string   `json:"proofs"`
}

func (dtb DelayedTransactionsBatch) MarshalJSON() ([]byte, error) {
	type alias DelayedTransactionsBatch

	if dtb.DelayedTransactions == nil {
		dtb.DelayedTransactions = make([]map[string]string, 0)
	}

	if dtb.Proofs == nil {
		dtb.Proofs = make(map[string]string)
	}

	return json.Marshal(alias(dtb))
}

func (dtb *DelayedTransactionsBatch) UnmarshalJSON(data []byte) error {
	type alias DelayedTransactionsBatch
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.DelayedTransactions == nil {
		aux.DelayedTransactions = make([]map[string]string, 0)
	}
	if aux.Proofs == nil {
		aux.Proofs = make(map[string]string)
	}
	*dtb = DelayedTransactionsBatch(aux)
	return nil
}
