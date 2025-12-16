package external_structs

import "encoding/json"

type Transaction struct {
	V       uint           `json:"v"`
	From    string         `json:"from"`
	To      string         `json:"to"`
	Amount  uint64         `json:"amount"`
	Fee     uint64         `json:"fee"`
	Sig     string         `json:"sig"`
	Nonce   uint64         `json:"nonce"`
	Payload map[string]any `json:"payload"`
}

func (t Transaction) MarshalJSON() ([]byte, error) {
	type alias Transaction
	if t.Payload == nil {
		t.Payload = make(map[string]any)
	}
	return json.Marshal((alias)(t))
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type alias Transaction
	var aux alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Payload == nil {
		aux.Payload = make(map[string]any)
	}
	*t = Transaction(aux)
	return nil
}
