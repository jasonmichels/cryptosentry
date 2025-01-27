package alert

import (
	"log"
	"time"
)

// timePrice captures the price at a specific moment in time.
type timePrice struct {
	timestamp time.Time
	price     float64
}

// historicalData: coin -> slice of (time, price) data
var historicalData = make(map[string][]timePrice)

// CheckThreshold logs if current price crosses a threshold.
func CheckThreshold(coinID string, currentPrice, threshold float64) {
	if currentPrice >= threshold {
		log.Printf("ALERT: %s price ($%.4f) hit threshold ($%.4f)\n",
			coinID, currentPrice, threshold)
	}
}

// CheckPriceMove logs if price jumped by a specified percentage within `window`.
func CheckPriceMove(coinID string, currentPrice float64, window time.Duration, percentage float64) {
	if len(historicalData[coinID]) == 0 {
		return
	}
	oldest := historicalData[coinID][0] // earliest data point in window
	if currentPrice >= oldest.price*(1+(percentage/100)) {
		log.Printf("ALERT: %s price jumped %.1f%%+ in < %s. Old: $%.4f, New: $%.4f\n",
			coinID, percentage, window, oldest.price, currentPrice)
	}
}

// UpdateHistory adds current price data and prunes data older than `window`.
func UpdateHistory(coinID string, currentPrice float64, window time.Duration) {
	now := time.Now()
	// Add new data
	historicalData[coinID] = append(historicalData[coinID], timePrice{
		timestamp: now,
		price:     currentPrice,
	})

	// Prune old data
	cutoff := now.Add(-window)
	var pruned []timePrice
	for _, tp := range historicalData[coinID] {
		if tp.timestamp.After(cutoff) {
			pruned = append(pruned, tp)
		}
	}
	historicalData[coinID] = pruned
}
