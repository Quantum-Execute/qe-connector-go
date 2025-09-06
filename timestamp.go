package qe_connector

import (
	"context"
	"encoding/json"
	"net/http"
)

// TimestampService get service timestamp milli
type TimestampService struct {
	c *Client
}

// Do send request
func (s *TimestampService) Do(ctx context.Context, opts ...RequestOption) (res int64, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/timestamp",
		secType:  secTypeNone,
	}
	m := params{}
	r.setParams(m)
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	resp := new(TimestampMessage)
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return 0, err
	}
	return resp.ServerTimeMilli, nil
}

type TimestampMessage struct {
	ServerTimeMilli int64 `json:"serverTimeMilli"`
}
