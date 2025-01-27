package config_test

import (
	"os"
	"testing"

	"github.com/jasonmichels/cryptosentry/internal/config"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "config_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some JSON
	sample := `{
      "coins": [
        { "id": "bitcoin", "threshold": 50000 }
      ],
      "intervalSeconds": 15,
      "priceMoveWindowMinutes": 2,
      "priceMovePercentage": 15
    }`
	if _, err := tempFile.WriteString(sample); err != nil {
		t.Fatalf("Failed to write temp config: %v", err)
	}
	tempFile.Close()

	cfg, err := config.LoadConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	// Basic assertions
	if len(cfg.Coins) != 1 {
		t.Errorf("Expected 1 coin, got %d", len(cfg.Coins))
	}
	if cfg.Coins[0].ID != "bitcoin" {
		t.Errorf("Expected coin ID 'bitcoin', got '%s'", cfg.Coins[0].ID)
	}
	if cfg.IntervalSeconds != 15 {
		t.Errorf("Expected IntervalSeconds=15, got %d", cfg.IntervalSeconds)
	}
	if cfg.PriceMoveWindowMinutes != 2 {
		t.Errorf("Expected PriceMoveWindowMinutes=2, got %d", cfg.PriceMoveWindowMinutes)
	}
	if cfg.PriceMovePercentage != 15 {
		t.Errorf("Expected PriceMovePercentage=15, got %.2f", cfg.PriceMovePercentage)
	}
}

// BenchmarkLoadConfig measures the performance of loading config.
func BenchmarkLoadConfig(b *testing.B) {
	// Create a temporary config file or reference a local fixture file
	tempFile, err := os.CreateTemp("", "config_bench_*.json")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	sample := `{
      "coins": [
        { "id": "bitcoin", "threshold": 50000 },
        { "id": "cardano", "threshold": 2.0 }
      ],
      "intervalSeconds": 15,
      "priceMoveWindowMinutes": 2,
      "priceMovePercentage": 15
    }`
	if _, err := tempFile.WriteString(sample); err != nil {
		b.Fatalf("Failed to write temp config: %v", err)
	}
	tempFile.Close()

	// Reset timer so that file creation time isn't counted in the benchmark loop
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := config.LoadConfig(tempFile.Name())
		if err != nil {
			b.Fatalf("LoadConfig returned error: %v", err)
		}
	}
}
