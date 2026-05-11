package qe_connector

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// signLikeBackend reproduces the backend's signing pipeline
// (`apiAuth.CollectParamsAndBodyForSign`) for tests:
//   - URL query keys (excluding `signature`) → url.Values p
//   - JSON body top-level keys → url.Values b (`scalarToString`-style for
//     scalars, `json.Marshal` for arrays/objects)
//   - merged := p ∪ b ; signature := HMAC_SHA256(secret, merged.Encode())
//
// Keep this in sync with backend `internal/server/gin/middleware.go`.
func signLikeBackend(t *testing.T, secret, rawQuery string, bodyBytes []byte) string {
	t.Helper()
	merged := url.Values{}
	if rawQuery != "" {
		qv, err := url.ParseQuery(rawQuery)
		if err != nil {
			t.Fatalf("parse query: %v", err)
		}
		for k, vs := range qv {
			if strings.EqualFold(k, "signature") {
				continue
			}
			for _, v := range vs {
				merged.Add(k, v)
			}
		}
	}
	if len(bytes.TrimSpace(bodyBytes)) > 0 {
		var obj map[string]interface{}
		dec := json.NewDecoder(bytes.NewReader(bodyBytes))
		dec.UseNumber()
		if err := dec.Decode(&obj); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		for k, v := range obj {
			if strings.EqualFold(k, "signature") {
				continue
			}
			if s, ok := backendScalarToString(v); ok {
				merged.Add(k, s)
				continue
			}
			if b, err := json.Marshal(v); err == nil {
				merged.Add(k, string(b))
			}
		}
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(merged.Encode()))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func backendScalarToString(v interface{}) (string, bool) {
	switch tv := v.(type) {
	case string:
		return tv, true
	case json.Number:
		return tv.String(), true
	case float64:
		if math.IsNaN(tv) || math.IsInf(tv, 0) {
			return "", false
		}
		return strconv.FormatFloat(tv, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(tv), true
	case nil:
		return fmt.Sprint(tv), true
	default:
		return "", false
	}
}

// TestV2SignAlignment_CreateMasterOrder mounts a tiny test server that mimics
// the backend's signature verification, then drives the real V2 SDK call path.
// If the SDK and the backend ever drift on signing rules, this test fails.
func TestV2SignAlignment_CreateMasterOrder(t *testing.T) {
	const apiKey = "test-api-key"
	const secret = "test-secret"

	var sawCorrectSignature bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-MBX-APIKEY"); got != apiKey {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			http.Error(w, "bad content-type", http.StatusBadRequest)
			return
		}
		body, _ := io.ReadAll(r.Body)
		expected := signLikeBackend(t, secret, r.URL.RawQuery, body)
		got := r.URL.Query().Get("signature")
		if got != expected {
			t.Errorf("signature mismatch:\n  got     = %s\n  want    = %s\n  query   = %s\n  body    = %s",
				got, expected, r.URL.RawQuery, string(body))
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		sawCorrectSignature = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":{"masterOrderId":"mo_xxx","status":"NEW","clientOrderId":"cli_1"}}`))
	}))
	defer srv.Close()

	client := NewClient(apiKey, secret, srv.URL)
	reply, err := client.NewCreateMasterOrderV2Service().
		ApiKeyId("binding-uuid").
		Exchange("Binance").
		MarketType("PERP").
		Symbol("BTCUSDT").
		Side("buy").
		Algorithm("TWAP").
		ExecutionDurationSeconds(3600).
		StartTimeMs(1760000000000).
		TotalQuantity("0.5").
		MarginType("U").
		WorstPrice("70000").
		MustComplete(true).
		MakerRateLimit("0.1").
		PovLimit("0.8").
		EnableMake(true).
		ClientOrderId("cli_1").
		Notes("v2 sign alignment test").
		Do(context.Background())

	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if !sawCorrectSignature {
		t.Fatal("server never observed a request with a matching signature")
	}
	if reply.MasterOrderId != "mo_xxx" || reply.Status != "NEW" || reply.ClientOrderId != "cli_1" {
		t.Fatalf("unexpected reply: %#v", reply)
	}
}

// TestV2SignAlignment_BatchCancel exercises the array-valued body path
// (`masterOrderIds: [...]`), which goes through the json.Marshal fallback
// in both the SDK and backend signing logic.
func TestV2SignAlignment_BatchCancel(t *testing.T) {
	const apiKey = "test-api-key"
	const secret = "test-secret"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		expected := signLikeBackend(t, secret, r.URL.RawQuery, body)
		got := r.URL.Query().Get("signature")
		if got != expected {
			t.Errorf("batch-cancel signature mismatch:\n  got  = %s\n  want = %s\n  body = %s",
				got, expected, string(body))
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"message":{"successCount":2,"failedOrders":[{"masterOrderId":"mo_c","reason":"already done"}]}}`))
	}))
	defer srv.Close()

	client := NewClient(apiKey, secret, srv.URL)
	reply, err := client.NewBatchCancelMasterOrdersV2Service().
		MasterOrderIds([]string{"mo_a", "mo_b", "mo_c"}).
		Reason("portfolio rebalance").
		Do(context.Background())
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if reply.SuccessCount != 2 || len(reply.FailedOrders) != 1 || reply.FailedOrders[0].MasterOrderId != "mo_c" {
		t.Fatalf("unexpected reply: %#v", reply)
	}
}

// TestCreateMasterOrderV2Validation covers SDK-side validation messages so
// callers fail fast before hitting the network.
func TestCreateMasterOrderV2Validation(t *testing.T) {
	c := NewClient("k", "s", "http://localhost:0")

	cases := []struct {
		name string
		mut  func(s *CreateMasterOrderV2Service)
		want string
	}{
		{
			name: "missing apiKeyId",
			mut: func(s *CreateMasterOrderV2Service) {
				s.Exchange("Binance").MarketType("SPOT").Symbol("BTCUSDT").
					Side("buy").Algorithm("TWAP").ExecutionDurationSeconds(60).TotalQuantity("0.1")
			},
			want: "apiKeyId is required",
		},
		{
			name: "executionDuration too short",
			mut: func(s *CreateMasterOrderV2Service) {
				s.ApiKeyId("k").Exchange("Binance").MarketType("SPOT").Symbol("BTCUSDT").
					Side("buy").Algorithm("TWAP").ExecutionDurationSeconds(10).TotalQuantity("0.1")
			},
			want: "executionDurationSeconds must be greater than 10",
		},
		{
			name: "both quantity and notional",
			mut: func(s *CreateMasterOrderV2Service) {
				s.ApiKeyId("k").Exchange("Binance").MarketType("SPOT").Symbol("BTCUSDT").
					Side("buy").Algorithm("TWAP").ExecutionDurationSeconds(60).
					TotalQuantity("0.1").OrderNotional("100")
			},
			want: "exactly one of totalQuantity / orderNotional",
		},
		{
			name: "target position needs quantity",
			mut: func(s *CreateMasterOrderV2Service) {
				s.ApiKeyId("k").Exchange("Binance").MarketType("PERP").Symbol("BTCUSDT").
					Side("buy").Algorithm("TWAP").ExecutionDurationSeconds(60).
					OrderNotional("100").IsTargetPosition(true)
			},
			want: "totalQuantity is required when isTargetPosition is true",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := c.NewCreateMasterOrderV2Service()
			tc.mut(s)
			_, err := s.Do(context.Background())
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want substring %q", err, tc.want)
			}
		})
	}
}
