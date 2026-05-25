package qe_connector

import (
	"testing"

	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

func TestExchangeBybitEnumValue(t *testing.T) {
	if got := string(trading_enums.ExchangeBybit); got != "Bybit" {
		t.Fatalf("ExchangeBybit = %q, want %q", got, "Bybit")
	}
}
