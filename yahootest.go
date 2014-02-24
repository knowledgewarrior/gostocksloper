

package main

import (
    "fmt"
	"encoding/csv"
	"net/http"
	"log"
)

func main() {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=GOOG&a=01&b=18&c=2014&d=01&e=21&f=2014&g=d")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error retrieving url: %s", err)
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error reading csv: %s", err)
	}
	//fmt.Println(records)
	records = append(records[:0], records[0+1:]...) 
	//fmt.Println(records)

	for _, record := range records {
   	fmt.Println(record)
   }
}
