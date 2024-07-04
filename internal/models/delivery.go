package models

import "encoding/json"

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func (d *Delivery) FromJSON(data []byte) error {
	return json.Unmarshal(data, d)
}
