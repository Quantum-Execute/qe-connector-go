// Preview a published option order; pass -submit to place a real order.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	qe "github.com/Quantum-Execute/qe-connector-go"
	"github.com/Quantum-Execute/qe-connector-go/constant/enums/trading_enums"
)

func main() {
	symbol := flag.String("symbol", "", "full name from the published OPTION list")
	quantity := flag.String("quantity", "", "base-asset quantity, not contracts")
	premium := flag.String("premium-limit", "", "premium per base asset, in USDC/BTC/ETH")
	side := flag.String("side", "buy", "buy or sell")
	submit := flag.Bool("submit", false, "place a real order using QE_API_KEY/QE_API_SECRET/QE_EXCHANGE_API_KEY_ID")
	flag.Parse()
	if *symbol == "" || *quantity == "" || *premium == "" || (*side != "buy" && *side != "sell") {
		log.Fatal("provide -symbol, -quantity, -premium-limit and a buy/sell -side")
	}
	client := qe.NewClient(os.Getenv("QE_API_KEY"), os.Getenv("QE_API_SECRET"))
	if base := os.Getenv("QE_BASE_URL"); base != "" {
		client.BaseURL = base
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rules, err := client.NewTradingPairsV2Service().Exchange(trading_enums.ExchangeDeribit).
		MarketType(trading_enums.MarketTypeOption).Symbol(*symbol).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(rules.Items) != 1 || rules.Items[0].Symbol != *symbol || !rules.Items[0].RulesReady || rules.Items[0].Status == "EXPIRED" {
		log.Fatal("published active option with executable price rules required")
	}
	printJSON("Contract rules", rules.Items[0])
	printJSON("Order preview (server revalidates quantity, ticks and expiry)", map[string]interface{}{
		"apiKeyId": os.Getenv("QE_EXCHANGE_API_KEY_ID"), "exchange": "Deribit", "marketType": "OPTION",
		"symbol": *symbol, "side": *side, "algorithm": "TWAP", "executionDurationSeconds": 1800,
		"totalQuantity": *quantity, "worstPrice": *premium, "mustComplete": false, "tailOrderProtection": false,
	})
	if !*submit {
		fmt.Println("Preview only. Add -submit to place this order.")
		return
	}
	for _, key := range []string{"QE_API_KEY", "QE_API_SECRET", "QE_EXCHANGE_API_KEY_ID"} {
		if os.Getenv(key) == "" {
			log.Fatalf("missing %s", key)
		}
	}
	reply, err := client.NewCreateMasterOrderV2Service().ApiKeyId(os.Getenv("QE_EXCHANGE_API_KEY_ID")).
		Exchange(trading_enums.ExchangeDeribit).MarketType(trading_enums.MarketTypeOption).
		Symbol(*symbol).Side(trading_enums.OrderSide(*side)).Algorithm(trading_enums.AlgorithmTWAP).
		ExecutionDurationSeconds(1800).TotalQuantity(*quantity).WorstPrice(*premium).
		MustComplete(false).TailOrderProtection(false).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	printJSON("Created order", reply)
}

func printJSON(label string, value interface{}) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:\n%s\n", label, data)
}
