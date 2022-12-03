package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
)

type Pool struct {
	Payouts           *big.Int `json:"payouts" bson:"payouts"`
	PayoutsUSD        float64  `json:"payouts_usd" bson:"payouts_usd"`
	MinutesTranscoded float64  `json:"minutes_transcoded" bson:"minutes_transcoded"`
}

func (p Pool) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Pool) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte array failed")
	}

	return json.Unmarshal(b, &p)
}
