package main

import (
	"fmt"
	"net/http"
	"service/stock_ticker/configuration"
	"service/stock_ticker/stockhandler"
)

func stockData(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func main() {

	// http.HandleFunc("/data", stockData)

	// http.ListenAndServe(":8090", nil)
	conf, _ := configuration.NewConfigurationFromEnv()
	stockhandler.GetStockData(conf)
}
