package stockhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"service/stock_ticker/configuration"
)

type avMetadata struct {
	Info          string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type dayItem struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

// GetStockData returns processed data
func GetStockData(conf configuration.StockDataConfig) error {
	link := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", conf.APIKey, conf.Symbol)
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Printf("Http code: %d", resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read responce: %s", err)
		return err
	}
	var respDataMap map[string]json.RawMessage
	err = json.Unmarshal(data, &respDataMap)
	if err != nil {
		log.Printf("Parsing resp body failed: %s", err)
		return err
	}
	if val, ok := respDataMap["Error Message"]; ok {
		var errMsg string
		err = json.Unmarshal(val, &errMsg)
		if err != nil {
			log.Printf("Parsing error message filed: %s", err)
			return err
		}
		log.Printf("Alphavantage api error %s", errMsg)
		err = errors.New(errMsg)
		return err
	}
	var metadata avMetadata
	err = json.Unmarshal(respDataMap["Meta Data"], &metadata)
	if err != nil {
		log.Println("Parsing metadata failed")
		log.Fatal(err)
	}
	var dayItemsArr map[string]dayItem
	err = json.Unmarshal(respDataMap["Time Series (Daily)"], &dayItemsArr)
	if err != nil {
		log.Println("Parsing dayItemsArr failed")
		log.Fatal(err)
	}
	log.Println(dayItemsArr)
	return nil
}
