package ingest

import (
	"encoding/json"
	"errors"
	"time"
)

type ExecutionStatus uint8

const (
	ExecutionSuccess = iota
	ExecutionUnknownError
	ExecutionDatabaseError
	ExecutionIOError
	ExecutionImageError
)

type duration time.Duration

func (d duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).Seconds())
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = duration(time.Duration(value * float64(time.Second)))
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Execution struct {
	Status       ProcessStatus `json:"status"`
	ErrorType    string        `json:"error_type,omitempty"`
	ErrorMessage string        `json:"error_message,omitempty"`
	ExecutedAt   time.Time     `json:"executed_at"`
	Inspect      duration      `json:"inspect"`
	Collect      duration      `json:"collect"`
	Plan         duration      `json:"plan"`
	Execute      duration      `json:"execute"`
	Store        duration      `json:"store"`
	Overall      duration      `json:"overall"`
}
