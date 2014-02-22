package main

import (
    "fmt"
    "bufio"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"
)

type CreateDBFunc func(string)
type StockPrice struct {
	symbol string
	price  float64
}

func main() {
    symbols, err := readLines("test-symbols.txt")
  	if err != nil {
    log.Fatalf("readLines error reading: %s", err)
  }
  	for _, symbol := range symbols {
    fmt.Println(symbol)

    var cdb CreateDBFunc
    cdb = createDB
    cdb(symbol)
  
  year := 2009
  closingPrice, err := getYahooInfo(symbol, year)
		if err == nil {
			fmt.Println(closingPrice)
		}
  }  
}

func createDB(s string) {
    //fmt.Printf("%s\n", s)
    os.Remove(s+".db")
    db, err := sql.Open("sqlite3", s+".db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE stockhistory (id INTEGER NOT NULL PRIMARY KEY, symbol TEXT, ydate TEXT, volume INTEGER, close, INTEGER);")
    if err != nil {
        log.Fatalln("could not create table:", err)
    }
    _, err = db.Exec("CREATE TABLE slopedata (id INTEGER NOT NULL PRIMARY KEY, symbol TEXT, slope FLOAT64, tradingdays INTEGER);")
    if err != nil {
        log.Fatalln("could not create table:", err)
    }
}

func getYahooInfo(symbol string, year int) (StockPrice, error) {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=11&b=01&c=%d&d=11&e=31&f=%d&g=m",
		symbol, year, year)
	resp, err := http.Get(url)
	if err != nil {
		return StockPrice{}, errors.New(fmt.Sprintf("Error making an HTTP request for stock %s.", symbol))
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err2 := csvReader.ReadAll()
	if err2 != nil {
		return StockPrice{}, errors.New(fmt.Sprintf("Error parsing CSV values for stock %s.", symbol))
	}
	closingPrice, err3 := strconv.ParseFloat(records[1][4], 64)
	if err3 != nil {
		return StockPrice{}, err3
	}
	return StockPrice{symbol, closingPrice}, nil
}


func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}



//func insertDb() {


	// tx, err := db.Begin()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// insert_stmt, err := tx.Prepare("insert into stocks(symbol) values(?)")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer insert_stmt.Close()
	// 	_, err = insert_stmt.Exec(symbol)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// tx.Commit()
	
	// rows, err := db.Query("select symbol from stocks")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var symbol string
	// 	rows.Scan(&symbol)
	// 	fmt.Println(symbol)
	// }
	// rows.Close()
	//}
	


