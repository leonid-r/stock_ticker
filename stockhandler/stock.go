package stockhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"service/stock_ticker/configuration"
	"sort"
	"strconv"
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

// StockData information about specific stock
type StockData struct {
	Symbol        string             `json:"Symbol"`
	LastRefreshed string             `json:"LastRefreshed"`
	TimeZone      string             `json:"TimeZone"`
	CloseAvarage  string             `json:"CloseAvarage"`
	TimeSeries    map[string]dayItem `json:"TimeSeriesDaily"`
}

// GetStockData returns processed data
func GetStockData(conf configuration.StockDataConfig) (StockData, error) {
	link := fmt.Sprintf("https://www.alphavantage.co/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", conf.APIKey, conf.Symbol)
	resp, err := http.Get(link)
	stockData := StockData{}
	if err != nil {
		log.Printf("Failed to make request: %s", err)
		return stockData, err
	}
	defer resp.Body.Close()
	log.Printf("Http code: %d", resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read responce: %s", err)
		return stockData, err
	}
	var respDataMap map[string]json.RawMessage
	err = json.Unmarshal(data, &respDataMap)
	if err != nil {
		log.Printf("Parsing resp body failed: %s", err)
		return stockData, err
	}
	if val, ok := respDataMap["Error Message"]; ok {
		var errMsg string
		err = json.Unmarshal(val, &errMsg)
		if err != nil {
			log.Printf("Parsing error message filed: %s", err)
			return stockData, err
		}
		log.Printf("Alphavantage api error %s", errMsg)
		err = errors.New(errMsg)
		return stockData, err
	}
	var metadata avMetadata
	err = json.Unmarshal(respDataMap["Meta Data"], &metadata)
	if err != nil {
		log.Println("Parsing metadata failed")
		return stockData, err
	}
	var dayItemsArr map[string]dayItem
	err = json.Unmarshal(respDataMap["Time Series (Daily)"], &dayItemsArr)
	if err != nil {
		log.Println("Parsing dayItemsArr failed")
		return stockData, err
	}
	var timeSeriesDataNDays map[string]dayItem
	timeSeriesDataNDays, err = getLastNDaysStockData(int(conf.NDays), dayItemsArr)
	if err != nil {
		log.Println("Getting lst nDays for source data failed")
		return stockData, err
	}
	closeAvarage, err := stockCloseAverage(timeSeriesDataNDays)
	if err != nil {
		return stockData, err
	}
	var closeAvarageStr string = strconv.FormatFloat(closeAvarage, 'f', 4, 64)
	stockData = StockData{
		Symbol:        metadata.Symbol,
		LastRefreshed: metadata.LastRefreshed,
		TimeZone:      metadata.TimeZone,
		CloseAvarage:  closeAvarageStr,
		TimeSeries:    timeSeriesDataNDays,
	}
	return stockData, nil
}

func getLastNDaysStockData(NDays int, timeSeriesData map[string]dayItem) (map[string]dayItem, error) {
	var err error
	timeSeriesDataNDays := make(map[string]dayItem)
	if len(timeSeriesData) < NDays {
		err = errors.New("Requested data size is bigger than source data")
	}
	keys := make([]string, 0, len(timeSeriesData))
	for k := range timeSeriesData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys[len(keys)-NDays:] {
		timeSeriesDataNDays[k] = timeSeriesData[k]
	}
	return timeSeriesDataNDays, err
}

func stockCloseAverage(timeSeriesData map[string]dayItem) (float64, error) {
	var closeAvarage float64 = 0
	var err error = nil
	for _, v := range timeSeriesData {
		close, err := strconv.ParseFloat(v.Close, 64)
		if err != nil {
			log.Printf("Calculating close avarage failed: %s", err)
			return closeAvarage, err
		}
		closeAvarage += close
	}
	closeAvarage = closeAvarage / float64(len(timeSeriesData))
	return closeAvarage, err
}
