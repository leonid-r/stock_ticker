package main

import (
	"encoding/json"
	"log"
	"net/http"
	"service/stock_ticker/configuration"
	"service/stock_ticker/stockhandler"
)

func stockData(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conf, err := configuration.NewConfigurationFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	stockData, err := stockhandler.GetStockData(conf)
	if err != nil {
		log.Println("Getting stock data failed")
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.NewEncoder(w).Encode(stockData)
	if err != nil {
		log.Println("Creating responce payload failed")
		http.Error(w, err.Error(), 500)
		return
	}
}

func main() {

	http.HandleFunc("/data", stockData)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
