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

// FlexDecimalString handles JSON values for decimal-shaped fields (averagePrice,
// filledNotional, etc.) that may come back as either a JSON number or a string.
// Backend V2 is contracted to return strings to avoid JS precision loss, but
// older backend builds or alternative deployments may still emit numbers. This
// type unmarshal both forms into a canonical string, leaving downstream callers
// free to use shopspring/decimal or big.Float as they like.
type FlexDecimalString string

func (f *FlexDecimalString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*f = ""
		return nil
	}
	// JSON string: trim the surrounding quotes via json.Unmarshal.
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexDecimalString(s)
		return nil
	}
	// JSON number: preserve the raw textual form to keep precision (json.Number).
	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return err
	}
	*f = FlexDecimalString(n.String())
	return nil
}

func (f FlexDecimalString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexDecimalString) String() string {
	return string(f)
}
