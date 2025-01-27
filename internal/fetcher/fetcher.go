package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PriceMap map[string]map[string]float64

func FetchPrices(coinIDs []string, vsCurrency string) (map[string]float64, error) {
	// Join coinIDs into a comma-separated string
	// e.g., "aurora-near,vechain,cardano,sophiaverse,bitcoin"
	var coinList string
	for i, c := range coinIDs {
		if i == 0 {
			coinList = c
		} else {
			coinList += "," + c
		}
	}

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		coinList, vsCurrency,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var pm PriceMap
	if err := json.NewDecoder(resp.Body).Decode(&pm); err != nil {
		return nil, err
	}

	// Flatten structure into coin -> price
	results := make(map[string]float64)
	for coinID, data := range pm {
		results[coinID] = data[vsCurrency]
	}

	return results, nil
}
