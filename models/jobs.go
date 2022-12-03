package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
	"time"
)

type Job struct {
	Timestamp  int64         `json:"timestamp" bson:"timestamp"`
	Transcoder string        `json:"transcoder" bson:"transcoder"` // eth address
	Node       string        `json:"node" bson:"node"`             // NEVER EXPOSE EXTERNALLY ON GET ROUTES (USED TO TRACK JOBS PER NODES FOR AUTHENTICATED ROUTES)
	Payout     *big.Int      `json:"payout" bson:"payout"`         // in wei
	PayoutUSD  float64       `json:"payout_usd" bson:"payout_usd"` // in usd (just used as historical data, not used for accounting)
	Pixels     int64         `json:"pixels" bson:"pixels"`         // amount of pixels
	Duration   time.Duration `json:"duration" bson:"duration"`     // segment length
	Profiles   []string      `json:"profiles" bson:"profiles"`     // video profiles
}

func (j Job) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *Job) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte array failed")
	}

	return json.Unmarshal(b, &j)
}
