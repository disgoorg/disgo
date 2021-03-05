package models

import "encoding/json"

type DataModel struct{}

func (p DataModel) MarshalJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *DataModel) UnmarshalJson(payload []byte) error {
	return json.Unmarshal(payload, p)
}
