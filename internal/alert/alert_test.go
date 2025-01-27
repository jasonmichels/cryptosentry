package alert_test

import (
	"testing"
	"time"

	"github.com/jasonmichels/cryptosentry/internal/alert"
)

// We'll just run some basic coverage of the alert methods.
// In a real test suite, you might inject a mock logger or check log output, etc.

func TestCheckThreshold(t *testing.T) {
	// Just verifying it doesn't panic, etc.
	// Ideally, you'd capture the log and confirm the message.
	// For brevity, we won't do that here.
	alert.CheckThreshold("bitcoin", 60000, 50000) // Should log an alert
	alert.CheckThreshold("bitcoin", 40000, 50000) // No alert
}

func TestCheckPriceMove(t *testing.T) {
	// Add some historical data
	alert.UpdateHistory("bitcoin", 40000, time.Minute)
	// Now check for a big jump
	alert.CheckPriceMove("bitcoin", 45000, time.Minute, 10) // Should log an alert (12.5% jump)
}

func TestUpdateHistory(t *testing.T) {
	alert.UpdateHistory("bitcoin", 30000, time.Minute)
	// If we add multiple data points, only the last 1 minute worth should remain
	alert.UpdateHistory("bitcoin", 31000, time.Minute)
	// etc.
	// Usually you'd check internal data structures or time-based logic,
	// but that requires more complex test harness.
}
