package alphavantage

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// readIntradayData reads converts the JSON object to a map and reads the necessary data
func readIntradayData(httpResponseBody io.Reader, interval *int64) [][]string {
	// Read response as JSON
	// create map interface
	var i map[string]map[string]map[string]interface{}
	json.NewDecoder(httpResponseBody).Decode(&i)

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
			fmt.Println("Error: 'open' not a string")
			return nil
		}
		// get high
		if str, ok := v["2. high"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			fmt.Println("Error: 'high' not a string")
			return nil
		}
		// geht low
		if str, ok := v["3. low"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			fmt.Println("Error: 'low' not a string")
			return nil
		}
		// get volume
		if str, ok := v["5. volume"].(string); ok {
			buffer.WriteString(",")
			buffer.WriteString(str)
		} else {
			fmt.Println("Error: 'volume' not a string")
			return nil
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
	return table
}

// readGlobalQuoteData converts the JSON object to a map and reads the data
func readGlobalQuoteData(httpResponseBody io.Reader) [][]string {
	// Read response as JSON
	// create map interface
	var i map[string]map[string]interface{}
	json.NewDecoder(httpResponseBody).Decode(&i)

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
		fmt.Println("Error: 'symbol' not a string")
		return nil
	}
	// get open
	if str, ok := i["Global Quote"]["02. open"].(string); ok {
		dataRow[1] = str
	} else {
		fmt.Println("Error: 'open' not a string")
		return nil
	}
	// get high
	if str, ok := i["Global Quote"]["03. high"].(string); ok {
		dataRow[2] = str
	} else {
		fmt.Println("Error: 'high' not a string")
		return nil
	}
	// get low
	if str, ok := i["Global Quote"]["04. low"].(string); ok {
		dataRow[3] = str
	} else {
		fmt.Println("Error: 'low' not a string")
		return nil
	}
	// get price
	if str, ok := i["Global Quote"]["05. price"].(string); ok {
		dataRow[4] = str
	} else {
		fmt.Println("Error: 'price' not a string")
		return nil
	}
	// get volume
	if str, ok := i["Global Quote"]["06. volume"].(string); ok {
		dataRow[5] = str
	} else {
		fmt.Println("Error: 'volume' not a string")
		return nil
	}
	// get latest trading day
	if str, ok := i["Global Quote"]["07. latest trading day"].(string); ok {
		dataRow[6] = str
	} else {
		fmt.Println("Error: 'latest trading day' not a string")
		return nil
	}
	// get previous close
	if str, ok := i["Global Quote"]["08. previous close"].(string); ok {
		dataRow[7] = str
	} else {
		fmt.Println("Error: 'previous close' not a string")
		return nil
	}
	// get change
	if str, ok := i["Global Quote"]["09. change"].(string); ok {
		dataRow[8] = str
	} else {
		fmt.Println("Error: 'chagne' not a string")
		return nil
	}
	// get change percent
	if str, ok := i["Global Quote"]["10. change percent"].(string); ok {
		dataRow[9] = str
	} else {
		fmt.Println("Error: 'change percent' not a string")
		return nil
	}

	// append to table
	table = append(table, dataRow)

	return table
}

// WriteToCSV takes a 2D Array and a pointer to the function and writes the data to a csv file
func WriteToCSV(table *[][]string, function *string) {
	// write to csv file
	// create new file
	fmt.Println("Writing CSV File")
	// truncate existing files, create non-existing files, open in write only mode
	file, err := os.OpenFile(*function+".csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create writer
	csvWriter := csv.NewWriter(file)

	// write file
	writerFailiure := csvWriter.WriteAll(*table)
	if writerFailiure != nil {
		fmt.Println(writerFailiure)
	}
}
