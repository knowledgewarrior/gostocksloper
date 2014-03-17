

package main

import (
    "fmt"
	"encoding/csv"
	"net/http"
	"log"
	"io"
)

func main() {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=GOOG&a=01&b=18&c=2014&d=01&e=21&f=2014&g=d")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error retrieving url: %s", err)
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	lineCount := 0
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// record is an array of string so is directly printable
		fmt.Println("Record", lineCount, "is", record, "and has", len(record), "fields")
		// and we can iterate on top of that
		for i := 0; i < len(record); i++ {
			fmt.Println(" ", record[i])
		}
		fmt.Println()
		lineCount += 1
	}
}
