package qe_connector

import (
	"encoding/json"
	"strconv"
)

// FlexInt64 handles JSON values that may be either a number or a string-encoded number.
// Protobuf JSON serialization encodes int64 as strings, while standard encoding/json
// uses numbers. This type transparently handles both representations.
type FlexInt64 int64

func (f *FlexInt64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	// Try as number first (standard JSON)
	var n int64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexInt64(n)
		return nil
	}

	// Try as quoted string (protobuf JSON format for int64)
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return &json.UnmarshalTypeError{
			Value: string(data),
			Type:  nil,
		}
	}

	if s == "" {
		*f = 0
		return nil
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*f = FlexInt64(n)
	return nil
}

func (f FlexInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(f))
}

func (f FlexInt64) Int64() int64 {
	return int64(f)
}
