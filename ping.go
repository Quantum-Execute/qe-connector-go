package qe_connector

import (
	"context"
	"net/http"
)

// PingService ping to server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/ping",
		secType:  secTypeNone,
	}
	m := params{}
	r.setParams(m)
	_, err = s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}
	return nil
}
