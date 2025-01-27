package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jasonmichels/cryptosentry/internal/alert"
	"github.com/jasonmichels/cryptosentry/internal/config"
	"github.com/jasonmichels/cryptosentry/internal/fetcher"
)

// We'll track the last price of each coin to determine "up" or "down".
var lastPrices = make(map[string]float64)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	// 2. Build coin ID list
	coinIDs := make([]string, len(cfg.Coins))
	for i, coinConf := range cfg.Coins {
		coinIDs[i] = coinConf.ID
	}

	// 3. Validate coin IDs
	if err := validateCoins(coinIDs); err != nil {
		log.Fatalf("Coin ID validation failed: %v", err)
	}

	// 4. Start the ticker
	ticker := time.NewTicker(cfg.IntervalDuration())
	defer ticker.Stop()

	for {
		// 5. Fetch prices
		prices, err := fetcher.FetchPrices(coinIDs, "usd")
		if err != nil {
			log.Println("Error fetching prices:", err)
			<-ticker.C
			continue
		}

		// 6. Log all coin prices & trend
		for coinID, currentPrice := range prices {
			// Determine trend
			lastPrice, seen := lastPrices[coinID]
			var arrow string
			if !seen || lastPrice == 0 {
				arrow = " (no prior data)"
			} else {
				switch {
				case currentPrice > lastPrice:
					arrow = " ↑" // trending up
				case currentPrice < lastPrice:
					arrow = " ↓" // trending down
				default:
					arrow = " →" // unchanged
				}
			}

			log.Printf("Price update: %s => $%.4f%s\n", coinID, currentPrice, arrow)

			// Store the current price for next iteration
			lastPrices[coinID] = currentPrice
		}

		log.Println("----------------------------------------------------")

		// 7. Check thresholds and big jumps
		for _, coinConf := range cfg.Coins {
			coinID := coinConf.ID
			currentPrice := prices[coinID]

			// 1) Threshold check
			alert.CheckThreshold(coinID, currentPrice, coinConf.Threshold)

			// 2) Big move check
			alert.CheckPriceMove(
				coinID,
				currentPrice,
				cfg.PriceMoveWindowDuration(),
				cfg.PriceMovePercentage,
			)

			// 3) Update history
			alert.UpdateHistory(coinID, currentPrice, cfg.PriceMoveWindowDuration())
		}

		<-ticker.C // Wait for next interval
	}
}

// validateCoins checks if CoinGecko recognizes all coinIDs
func validateCoins(coinIDs []string) error {
	log.Println("Validating coin IDs with CoinGecko...")
	prices, err := fetcher.FetchPrices(coinIDs, "usd")
	if err != nil {
		return err
	}

	var invalid []string
	for _, id := range coinIDs {
		if _, found := prices[id]; !found {
			invalid = append(invalid, id)
		} else {
			log.Printf("Coin ID '%s' is valid.", id)
		}
	}
	if len(invalid) > 0 {
		log.Printf("Invalid coin IDs found: %v", invalid)
		return fmt.Errorf("invalid coins: %v", invalid)
	}
	log.Println("All coin IDs are valid!")
	return nil
}
