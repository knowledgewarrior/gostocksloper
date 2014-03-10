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
      "strconv"
)


type GetStocksFunc func(string)
type GetSlopesFunc func(string)

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
} //readlines

func writeLines(lines []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, line := range lines {
    fmt.Fprintln(w, line)
  }
  return w.Flush()
}

func check(e error) {
    if e != nil {
        panic(e)
    }
} //check

func getStocks(symbol string) {

      t := time.Now().Format("2006-01-02")
	tArray := strings.Split(t, "-")

 	nowyear := tArray[0]
 	nmonth := tArray[1]
      nmonth2, err := strconv.Atoi(nmonth)
      check(err)
      nowmonth := (nmonth2 - 1)
 	nowday := tArray[2]

  	p := time.Now().AddDate(0, 0, -60).Format("2006-01-02")
      pArray := strings.Split(p, "-")
      thenyear := pArray[0]
      tmonth := pArray[1]
      tmonth2, err := strconv.Atoi(tmonth)
      check(err)
      thenmonth := (tmonth2 - 1)
      thenday := pArray[2]
      //fmt.Println(thenyear, thenmonth, thenday, " - ", nowyear,nowmonth,nowday)

	url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%s&c=%s&d=%d&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
	resp, err := http.Get(url)
	if resp.StatusCode > 399 {
		fmt.Println(symbol+": not on yahoo finance")
		return
	}
      defer resp.Body.Close()

	os.Remove(symbol+".db")
      db, err := sql.Open("sqlite3", symbol+".db")
	if err != nil {
		fmt.Println("error opening db")
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE stockhistory (id INTEGER NOT NULL PRIMARY KEY, ydate TEXT, closeprice INTEGER);")
      if err != nil {
        fmt.Println("could not create table")
      }
     _, err = db.Exec("CREATE TABLE slopedata (id INTEGER NOT NULL PRIMARY KEY, slope FLOAT64, tradingdays INTEGER);")
      if err != nil {
          fmt.Println("could not create table")
      }

      csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	check(err)
	records = append(records[:0], records[0+1:]...)
		for _, record := range records {
			d := record[0]
		   	c := record[4]

			tx, err := db.Begin()
				if err != nil {
					fmt.Println("error with db")
				}
			insert_stmt, err := tx.Prepare("insert into stockhistory(ydate,closeprice) values(?,?)")

			defer insert_stmt.Close()
			_, err = insert_stmt.Exec(d,c)
			tx.Commit()
		}

	var gslf GetSlopesFunc
	gslf = getSlope
	gslf(symbol)

} // getStocks

func getSlope(symbol string) {
	//get slope
	ntd := 35.00

	db, err := sql.Open("sqlite3", symbol+".db")
	check(err)
	defer db.Close()
	rows, err := db.Query("select sum(id) as sumx, sum(closeprice) as sumy, sum(id * closeprice) as sumxy, sum(id * id) as sumxx from(select id, closeprice from stockhistory order by ydate desc limit ?);", ntd)
		if err != nil {
			fmt.Println("error with select")
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
                  //fmt.Println(symbol,slope)

                  os.Remove("Slopes.csv")
                  fname := "Slopes.csv"
                  f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
                  check(err)
                  defer f.Close()

                  b := bufio.NewWriter(f)
                  defer func() {
                  if err = b.Flush(); err != nil {
                        fmt.Println(err)
                  }
                  }()
                  fmt.Fprint(b, symbol+",")
                  fmt.Fprint(b, slope)
                  fmt.Fprint(b, "\n")
            } // for rows
            rows.Close()
} //getSlope

func main(){
	symbols, err := readLines("testsymbols.txt")
    //symbols, err := readLines("stocks-testing.txt")
  	if err != nil {
    	fmt.Println("readLines error reading")
  	}
	for _, symbol := range symbols {
		var gsf GetStocksFunc
	 	gsf = getStocks
	 	gsf(symbol)
	}

} // func main





