package alphavantage

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	urlBase = "https://www.alphavantage.co/query?"
)

func buildQueryURL(function, symbol, apikey *string) string {
	var buffer bytes.Buffer
	// build string buffer
	buffer.WriteString(urlBase)
	// append query function
	buffer.WriteString("function=")
	buffer.WriteString(*function)
	// append symbol
	buffer.WriteString("&symbol=")
	buffer.WriteString(*symbol)
	// append apikey
	buffer.WriteString("&apikey=")
	buffer.WriteString(*apikey)
	// set datatype to csv
	buffer.WriteString("&datatype=csv")

	// return query string
	return buffer.String()
}

func sendRequestToAV(function, symbol, apikey *string) ([][]string, error) {
	// build query url
	queryURL := buildQueryURL(function, symbol, apikey)

	// send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, responseError := client.Get(queryURL)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	// return data
	return readCSVBody(response)
}

func readCSVBody(response *http.Response) ([][]string, error) {
	// read csv body
	reader := csv.NewReader(response.Body)
	var table [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		table = append(table, record)
	}
	return table, nil
}

// GetTimeSeriesIntraday returns the intraday data of the specified stock
// symbol: Symbol of the stock in question
// interval: Set the interval between datapoints. Valid options are 1, 5, 15, 30 and 60
// apikey: Set your API key
func GetTimeSeriesIntraday(symbol string, interval int64, apikey string) ([][]string, error) {
	// check if interval fulfills requirements
	if interval != 1 && interval != 5 && interval != 15 && interval != 30 && interval != 60 {
		return nil, errors.New("Interval does not meet the requirements. Possible values are 1, 5, 15, 30 and 60. Provided:" + strconv.FormatInt(interval, 10))
	}

	// build query url
	var buffer bytes.Buffer
	buffer.WriteString(urlBase)
	// append query function
	buffer.WriteString("function=TIME_SERIES_INTRADAY")
	// append symbol
	buffer.WriteString("&symbol=")
	buffer.WriteString(symbol)
	// append interval
	buffer.WriteString("&interval=")
	buffer.WriteString(strconv.FormatInt(interval, 10))
	buffer.WriteString("min")
	// append api key
	buffer.WriteString("&apikey=")
	buffer.WriteString(apikey)
	// set the datatype to csv
	buffer.WriteString("&datatype=csv")

	// send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	response, responseError := client.Get(buffer.String())
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	// return data
	return readCSVBody(response)
}

// GetTimeSeriesDaily returns the daily data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesDaily(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_DAILY"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)

}

// GetTimeSeriesDailyAdjusted returns the adjusted daily data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesDailyAdjusted(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_DAILY_ADJUSTED"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}

// GetTimeSeriesWeekly returns the weekly data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesWeekly(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_WEEKLY"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}

// GetTimeSeriesWeeklyAdjusted returns the adjusted weekly data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesWeeklyAdjusted(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_WEEKLY_ADJUSTED"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}

// GetTimeSeriesMonthly returns the monthly data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesMonthly(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_MONTHLY"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}

// GetTimeSeriesMonthlyAdjusted returns the adjusted monthly data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetTimeSeriesMonthlyAdjusted(symbol, apikey string) ([][]string, error) {
	// function
	function := "TIME_SERIES_MONTHLY_ADJUSTED"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}

// GetQuoteEndpoint returns the global quote data of the specified stock
// symbol: Symbol of the stock in question
// apikey: Set your API key
func GetQuoteEndpoint(symbol, apikey string) ([][]string, error) {
	// function
	function := "GLOBAL_QUOTE"
	// get data
	return sendRequestToAV(&function, &symbol, &apikey)
}
