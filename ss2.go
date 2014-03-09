/*
Jason Fowler February 2014
Linear Regression Analysis of Yahoo Stock Data using formula:
Slope(b) = (N*ΣXY - (ΣX)(ΣY)) / (N*ΣX2 - (ΣX)2)
*/

package main

import (
    "fmt"
    "bufio"
    "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"encoding/csv"
	"net/http"
	"time"
	"strings"
)

type CreateDBFunc func(string)
type GetSlopeFunc func(string)

func main() {
    symbols, err := readLines("testsymbols.txt") // very small
    //symbols, err := readLines("stocks-testing.txt") // 1649 symbols
  	if err != nil {
    	fmt.Println("readLines error reading: %s", err)
  	}
  	for _, symbol := range symbols {
         var cdb CreateDBFunc
	   cdb = createDB
	   cdb(symbol)
   }
         for _, symbol := range symbols {
	  	records, err := getYahooInfo(symbol)
	  	if err != nil {
	    	fmt.Println(symbol+": cannot get yahoo info")
	  	}
	  	  for _, record := range records {
			d := record[0]
		   	c := record[4]
               // fmt.Println(d,c)
		   	db, err := sql.Open("sqlite3", symbol+".db")
				if err != nil {
					fmt.Println(err)
				}
				defer db.Close()

				tx, err := db.Begin()
				if err != nil {
					fmt.Println(err)
				}
				insert_stmt, err := tx.Prepare("insert into stockhistory(ydate,closeprice) values(?,?)")
				if err != nil {
					fmt.Println(err)
				}
				defer insert_stmt.Close()
					_, err = insert_stmt.Exec(d,c)
					if err != nil {
						fmt.Println(err)
					}
				tx.Commit()
		  }
        }


 //  	for _, symbol := range symbols {
 //  		var gsf GetSlopeFunc
	// 	gsf = getSlope
	// 	gsf(symbol)
	// }

}

func createDB(s string) {
    //fmt.Printf("%s\n", s)
    os.Remove(s+".db")
    db, err := sql.Open("sqlite3", s+".db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE stockhistory (id INTEGER NOT NULL PRIMARY KEY, ydate TEXT, closeprice INTEGER);")
    if err != nil {
        fmt.Println("could not create table:", err)
    }
    _, err = db.Exec("CREATE TABLE slopedata (id INTEGER NOT NULL PRIMARY KEY, slope FLOAT64, tradingdays INTEGER);")
    if err != nil {
        fmt.Println("could not create table:", err)
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

func getYahooInfo(symbol string) ([][]string, error){

	t := time.Now().Format("2006-01-02")
	tArray := strings.Split(t, "-")

 	nowyear := tArray[0]
 	nowmonth := tArray[1]
 	nowday := tArray[2]

  	thenyear := 2014
  	thenmonth := 00
  	thenday := 01

	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%d&c=%d&d=%s&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
	resp, err := http.Get(url)
	if resp.StatusCode != 200 {
		break
	}
	if err != nil {
		fmt.Println(symbol+"error retrieving url: %s", err)
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
          fmt.Println(symbol+" error reading csv: %s", err)
	}
	records = append(records[:0], records[0+1:]...)
	return records, nil
}


func getSlope(s string) {
      ntd := 35.00
	db, err := sql.Open("sqlite3", s+".db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	rows, err := db.Query("select sum(id) as sumx, sum(closeprice) as sumy, sum(id * closeprice) as sumxy, sum(id * id) as sumxx from(select id, closeprice from stockhistory order by ydate desc limit ?);", ntd)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		for rows.Next() {
			var sumx float64
			var sumy float64
			var sumxy float64
			var sumxx float64
			rows.Scan(&sumx, &sumy, &sumxy, &sumxx)

			ntdsumxy := ntd * sumxy
			sumxsumy := sumx * sumy
			ntdsumxx := ntd * sumxx
			sumxsumx := sumx * sumx

			slope := (ntdsumxy - sumxsumy) / (ntdsumxx - sumxsumx)
                  fmt.Println(s,slope)
		}
		rows.Close()
}
