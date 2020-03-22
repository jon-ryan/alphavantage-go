package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// readIntradayCSV expects a csv file in the response body. It returns it as a 2D array
func readIntradayCSV(response *http.Response) ([][]string, error) {
	reader := csv.NewReader(response.Body)
	var results [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	return results, nil
}

// readIntradayData expects a JSON object in the response body. It returns it as a 2D array
func readIntradayDataJSON(response *http.Response, interval *int64) ([][]string, error) {
	// Read response as JSON
	// create map interface
	var i map[string]map[string]map[string]interface{}
	json.NewDecoder(response.Body).Decode(&i)

	// create table to store the values in
	recordList := make([]string, 0, 105)

	// Filter out the values contained in the 'Time Series' section
	timeSeriesValue := "Time Series (" + strconv.FormatInt(*interval, 10) + "min)"
	timeSeriesMap := i[timeSeriesValue]
	// iterate over all map elements
	for k, v := range timeSeriesMap {
		var buffer bytes.Buffer
		// get trading date and time
		buffer.WriteString(k)
		// make type assertion and get values
		// get open
		if str, ok := v["1. open"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			return nil, errors.New("Error: 'open' not a string")
		}
		// get high
		if str, ok := v["2. high"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			return nil, errors.New("Error: 'high' not a string")
		}
		// geht low
		if str, ok := v["3. low"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			return nil, errors.New("Error: 'low' not a string")
		}
		// get volume
		if str, ok := v["5. volume"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			return nil, errors.New("Error: 'volume' not a string")
		}

		recordList = append(recordList, buffer.String())

	}

	// sort table
	fmt.Println("Sorting entries")
	sort.Strings(recordList)

	// SPLIT ARRAY INTO TWO DIM ARRAY
	// create table and apend header
	table := make([][]string, 0, 105)
	row := make([]string, 5)
	row[0] = "Time"
	row[1] = "Open"
	row[2] = "High"
	row[3] = "Low"
	row[4] = "Volume"
	table = append(table, row)

	// insert slices
	for _, item := range recordList {
		row := strings.Split(item, ",")
		table = append(table, row)
	}

	return table, nil
}

// readGlobalQuoteCSV expects a csv file in the response body. It returns it as a 2D array
func readGlobalQuoteCSV(response *http.Response) ([][]string, error) {
	reader := csv.NewReader(response.Body)
	var results [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		results = append(results, record)
	}
	return results, nil
}

// readGlobalQuoteData expects a JSON object in the response body. It returns it as a 2D array
func readGlobalQuoteDataJSON(response *http.Response) ([][]string, error) {
	// Read response as JSON
	// create map interface
	var i map[string]map[string]interface{}
	json.NewDecoder(response.Body).Decode(&i)

	table := make([][]string, 0, 2)

	// create row containing the headers
	row := make([]string, 10)
	row[0] = "Symbol"
	row[1] = "Open"
	row[2] = "High"
	row[3] = "Low"
	row[4] = "Price"
	row[5] = "Volume"
	row[6] = "Latest trading day"
	row[7] = "Previous close"
	row[8] = "Change"
	row[9] = "Change percent"

	// append to table
	table = append(table, row)

	// fill header row
	dataRow := make([]string, 10)
	// fill data row
	// get symbol
	if str, ok := i["Global Quote"]["01. symbol"].(string); ok {
		dataRow[0] = str
	} else {
		return nil, errors.New("Error: 'symbol' not a string")
	}
	// get open
	if str, ok := i["Global Quote"]["02. open"].(string); ok {
		dataRow[1] = str
	} else {
		return nil, errors.New("Error: 'open' not a string")
	}
	// get high
	if str, ok := i["Global Quote"]["03. high"].(string); ok {
		dataRow[2] = str
	} else {
		return nil, errors.New("Error: 'high' not a string")
	}
	// get low
	if str, ok := i["Global Quote"]["04. low"].(string); ok {
		dataRow[3] = str
	} else {
		return nil, errors.New("Error: 'low' not a string")
	}
	// get price
	if str, ok := i["Global Quote"]["05. price"].(string); ok {
		dataRow[4] = str
	} else {
		return nil, errors.New("Error: 'price' not a string")
	}
	// get volume
	if str, ok := i["Global Quote"]["06. volume"].(string); ok {
		dataRow[5] = str
	} else {
		return nil, errors.New("Error: 'volume' not a string")
	}
	// get latest trading day
	if str, ok := i["Global Quote"]["07. latest trading day"].(string); ok {
		dataRow[6] = str
	} else {
		return nil, errors.New("Error: 'latest trading day' not a string")
	}
	// get previous close
	if str, ok := i["Global Quote"]["08. previous close"].(string); ok {
		dataRow[7] = str
	} else {
		return nil, errors.New("Error: 'previous close' not a string")
	}
	// get change
	if str, ok := i["Global Quote"]["09. change"].(string); ok {
		dataRow[8] = str
	} else {
		return nil, errors.New("Error: 'change' not a string")
	}
	// get change percent
	if str, ok := i["Global Quote"]["10. change percent"].(string); ok {
		dataRow[9] = str
	} else {
		return nil, errors.New("Error: 'change percent' not a string")
	}

	// append to table
	table = append(table, dataRow)

	return table, nil
}

// WriteToCSV takes a 2D Array and a pointer to the function and writes the data to a csv file
func WriteToCSV(table *[][]string, function *string) error {
	// write to csv file
	// create new file
	fmt.Println("Writing CSV File")
	// truncate existing files, create non-existing files, open in write only mode
	file, err := os.OpenFile(*function+".csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return error
	}
	defer file.Close()

	// create writer
	csvWriter := csv.NewWriter(file)

	// write file
	writerFailiure := csvWriter.WriteAll(*table)
	if writerFailiure != nil {
		return writerFailiure
	}

	return nil
}

func main() {
	// parse cmd line flags
	var apiKey = flag.String("apikey", "demo", "Set the API Key to connet to Alpha Vantage")
	var symbol = flag.String("symbol", "MSFT", "Define the stock you want to get the data from")
	var interval = flag.Int64("interval", 5, "Set the interval between data points. Valid are 1, 5, 15, 30 and 60")
	var function = flag.String("function", "intraday", "Specify what data you want to retrieve. Valid values are 'intraday' and 'quote'")
	var datatype = flag.String("datatype", "json", "Specify if the data should be retrieved as JSON or CSV")

	flag.Parse()

	// check for valid parameters
	if *interval != 1 && *interval != 5 && *interval != 15 && *interval != 30 && *interval != 60 {
		fmt.Println("Invalid interval provided. Provided value:", *interval, "\nValid Values: 1, 5, 15, 30, 60")
		return
	}

	if *function == "intraday" {
		*function = "TIME_SERIES_INTRADAY"
	} else if *function == "quote" {
		*function = "GLOBAL_QUOTE"
	} else {
		fmt.Println("Invalid function call. Provided value:", *function, "\nValid Values: 'intraday' or 'quote'")
		return
	}

	if *datatype != "json" && *datatype != "csv" {
		fmt.Println("Invalid datatype. Provided value:", *datatype, "\nValid Values: 'json' or 'csv'")
	}

	// print flags
	fmt.Println("--- Flags ---")
	fmt.Println("Symbol:", *symbol)
	fmt.Println("API Key:", *apiKey)
	fmt.Println("Interval:", *interval)
	fmt.Println("Function:", *function)
	fmt.Println("Datatype:", *datatype)
	fmt.Println("-------------\n")

	// create new client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// full url: https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=MSFT&interval=5min&apikey=demo

	urlBase := "https://www.alphavantage.co/query?"

	// append function
	urlBase += "function=" + *function
	// append symbol
	urlBase += "&symbol=" + *symbol
	// append interval
	urlBase += "&interval=" + strconv.FormatInt(*interval, 10) + "min"
	// append api key
	urlBase += "&apikey=" + *apiKey

	// if datatype is csv append to url
	if *datatype == "csv" {
		urlBase += "&datatype=csv"
	}

	fmt.Println("Sending Request")

	// send HTTP GET request
	httpResponse, err := client.Get(urlBase)
	if err != nil {
		fmt.Println(err)
	}
	defer httpResponse.Body.Close()

	fmt.Println("Decoding Response")

	// depending on the requested function a different handler has to be called
	var table [][]string
	switch *function {
	case "TIME_SERIES_INTRADAY":
		if *datatype == "csv" {
			table, parseError := readIntradayCSV(httpResponse)
			if table == nil && parseError != nil {
				log.Fatal(parseError)
			}
		} else {
			table, parseError := readIntradayDataJSON(httpResponse, interval)
			if table == nil && parseError != nil {
				log.Fatal(parseError)
			}
		}

		break

	case "GLOBAL_QUOTE":
		if *datatype == "csv" {
			table, parseError := readGlobalQuoteCSV(httpResponse)
			if table == nil && parseError != nil {
				log.Fatal(parseError)
			}
		} else {
			table, parseError := readGlobalQuoteDataJSON(httpResponse)
			if table == nil && parseError != nil {
				log.Fatal(parseError)
			}
		}

		break

	default:
		fmt.Println("Unknown function call")
		return
	}

	writerError := WriteToCSV(&table, function)
	if writerError != nil {
		log.Fatal(writerError)
	}

	fmt.Println("Done")

}
