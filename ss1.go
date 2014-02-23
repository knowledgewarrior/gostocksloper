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
	"time"
	"strings"
)

type CreateDBFunc func(string)

type StockHist struct {
	symbol string
	//volume int
	price  float64
}

func main() {
    symbols, err := readLines("stocks-testing.txt")
  	if err != nil {
    	log.Fatalf("readLines error reading: %s", err)
  }
  	for _, symbol := range symbols {
   	//fmt.Println(symbol)

    var cdb CreateDBFunc
    cdb = createDB
    cdb(symbol)
  
  	t := time.Now().Format("2006-01-02")
	tArray := strings.Split(t, "-")
 	
 	nowyear := tArray[0]
 	nowmonth := tArray[1]
 	nowday := tArray[2]
 
  	thenyear := 2014
  	thenmonth := 00
  	thenday := 01

  	results, err := getYahooInfo(symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
		if err == nil {
			//results = append(results, closingPrice)
			fmt.Println(results)
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

func getYahooInfo(symbol string, thenmonth int, thenday int, thenyear int, nowmonth string, nowday string, nowyear string) (StockHist, error) {
	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%s&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
	resp, err := http.Get(url)
	if err != nil {
		return StockHist{}, errors.New(fmt.Sprintf("Error making an HTTP request for stock %s.", symbol))
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err2 := csvReader.ReadAll()
	if err2 != nil {
		return StockHist{}, errors.New(fmt.Sprintf("Error parsing CSV values for stock %s.", symbol))
	}
	closingPrice, err3 := strconv.ParseFloat(records[1][4], 64)
	if err3 != nil {
		return StockHist{}, err3
	}
	return StockHist{symbol, closingPrice}, nil
}


// func getYahooInfo(symbol string, thenmonth int, thenday int, thenyear int, nowmonth string, nowday string, nowyear string) (StockPrice, error) {
// 	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%s&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return StockPrice{}, errors.New(fmt.Sprintf("Error making an HTTP request for stock %s.", symbol))
// 	}
// 	defer resp.Body.Close()
// 	csvReader := csv.NewReader(resp.Body)
// 	records, err2 := csvReader.ReadAll()
// 	//fmt.Println(records)
// 	if err2 != nil {
// 		return StockPrice{}, errors.New(fmt.Sprintf("Error parsing CSV values for stock %s.", symbol))
// 	}

// 	// for key, record := range records {
//  //   	fmt.Println(key, record)
//  //   	}
//    	//closingPrice, err3 := strconv.ParseFloat(records[1][4], 64)
//    	for goodrecords, err3 := range records {
// 	if err3 != nil {
// 		return StockPrice{}, err3
// 	}
// 	//return StockPrice{symbol, closingPrice}, nil
// 	return StockPrice{goodrecords}, nil
// 	}
// }






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
	


