package configuration

import (
	"fmt"
	"os"
	"strconv"
)

// StockDataConfig config for data
type StockDataConfig struct {
	NDays  uint64
	Symbol string
	APIKey string
}



// NewConfigurationFromEnv creates new config from env variables
func NewConfigurationFromEnv() (StockDataConfig, error) {
	ndays, err := strconv.ParseUint(os.Getenv("NDAYS"), 10, 64)
	config := StockDataConfig{}
	if err != nil {
		return config, fmt.Errorf("Parsing NDAYS env var error error: %v", err)
	}
	config.NDays = ndays
	config.Symbol = os.Getenv("SYMBOL")
	config.APIKey = os.Getenv("APIKEY")
	return config, nil
}
