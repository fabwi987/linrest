package yaho

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

//GetStocks returns several stocks in one query
func GetStocks(stockSymbols string) (*Response, error) {

	//Build the http request
	query := url.QueryEscape("select * from yahoo.finance.quote where symbol in ('" + stockSymbols + "')")
	destination := "http://query.yahooapis.com/v1/public/yql?q="
	specification := "&format=json&env=http://datatables.org/alltables.env"

	total := destination + query + specification
	//log.Println(total)
	resp, err := http.Get(total)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var stocks *Response

	if err := json.NewDecoder(resp.Body).Decode(&stocks); err != nil {
		return nil, err
	}

	return stocks, nil

}

func GetSingleStocks(stockSymbols string) (*SingleResponse, error) {

	//Build the http request
	query := url.QueryEscape("select * from yahoo.finance.quote where symbol in ('" + stockSymbols + "')")
	destination := "http://query.yahooapis.com/v1/public/yql?q="
	specification := "&format=json&env=http://datatables.org/alltables.env"

	total := destination + query + specification
	//log.Println(total)
	resp, err := http.Get(total)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var stocks *SingleResponse

	if err := json.NewDecoder(resp.Body).Decode(&stocks); err != nil {
		return nil, err
	}

	return stocks, nil

}

type Response struct {
	Query struct {
		Count   int       `json:"count"`
		Created time.Time `json:"created"`
		Lang    string    `json:"lang"`
		Results struct {
			Quote []Quote
		} `json:"results"`
	} `json:"query"`
}

type SingleResponse struct {
	Query struct {
		Count   int       `json:"count"`
		Created time.Time `json:"created"`
		Lang    string    `json:"lang"`
		Results struct {
			Quote Quote
		} `json:"results"`
	} `json:"query"`
}

type Quote struct {
	Symbol               string      `json:"symbol"`
	AverageDailyVolume   string      `json:"AverageDailyVolume"`
	Change               string      `json:"Change"`
	DaysLow              interface{} `json:"DaysLow"`
	DaysHigh             interface{} `json:"DaysHigh"`
	YearLow              string      `json:"YearLow"`
	YearHigh             string      `json:"YearHigh"`
	MarketCapitalization string      `json:"MarketCapitalization"`
	LastTradePriceOnly   string      `json:"LastTradePriceOnly"`
	DaysRange            interface{} `json:"DaysRange"`
	Name                 string      `json:"Name"`
	Symbol1              string      `json:"Symbol1"`
	Volume               string      `json:"Volume"`
	StockExchange        string      `json:"StockExchange"`
}
