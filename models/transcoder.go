package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/big"
)

type Transcoder struct {
	Timestamp         int64    `json:"timestamp" bson:"timestamp"`                   // creation time
	Address           string   `json:"address" bson:"address"`                       // transcoder eth address
	PendingBalance    *big.Int `json:"pending_balance" bson:"pending_balance"`       // pending balance
	Payout            *big.Int `json:"payout" bson:"payout"`                         // total paid out
	PayoutUSD         float64  `json:"payout_usd" bson:"payout_usd"`                 // total paid out in USD (USD value of pending balance is calculated upon payout and added here)
	SecondsTranscoded float64  `json:"seconds_transcoded" bson:"seconds_transcoded"` // seconds of video transcoded (time.Duration.Seconds)
}

type Node struct {
	ID                string  `json:"ID" bson:"ID"`                                 // id (name if provided or IP)
	Transcoder        string  `json:"transcoder" bson:"transcoder"`                 // eth address of transcoder
	Region            string  `json:"region" bson:"region"`                         // should use IATA codes of server locations
	Active            bool    `json:"active" bson:"active"`                         // has live connection
	SecondsTranscoded float64 `json:"seconds_transcoded" bson:"seconds_transcoded"` // seconds of video transcoded (time.Duration.Seconds)
}

func (t Transcoder) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Transcoder) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte array failed")
	}

	return json.Unmarshal(b, &t)
}

func (n Node) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *Node) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to byte array failed")
	}

	return json.Unmarshal(b, &n)
}
