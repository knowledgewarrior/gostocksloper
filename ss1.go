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
	"log"
	"os"
	"encoding/csv"
	"net/http"
	"time"
	"strings"
)

type CreateDBFunc func(string)
type GetSlopeFunc func(string)

func main() {
    symbols, err := readLines("stocks-testing.txt")
  	if err != nil {
    	log.Fatalf("readLines error reading: %s", err)
  	}
  	for _, symbol := range symbols {

	   var cdb CreateDBFunc
	   cdb = createDB
	   cdb(symbol)

	  	records, err := getYahooInfo(symbol)
	  	if err != nil {
	    	log.Fatalf("cannot get yahoo info for: %s", err)
	  	}
	  	for _, record := range records {
			d := record[0]
		   	c := record[4]
		   	v := record[5]
		   	//fmt.Println(d,c,v)

		   	db, err := sql.Open("sqlite3", symbol+".db")
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				tx, err := db.Begin()
				if err != nil {
					log.Fatal(err)
				}
				insert_stmt, err := tx.Prepare("insert into stockhistory(ydate,closeprice,volume) values(?,?,?)")
				if err != nil {
					log.Fatal(err)
				}
				defer insert_stmt.Close()
					_, err = insert_stmt.Exec(d,c,v)
					if err != nil {
						log.Fatal(err)
					}
				tx.Commit()
		}
  	}

  	for _, symbol := range symbols {
  		var gsf GetSlopeFunc
		gsf = getSlope
		gsf(symbol)
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
	_, err = db.Exec("CREATE TABLE stockhistory (id INTEGER NOT NULL PRIMARY KEY, ydate TEXT, volume INTEGER, closeprice INTEGER);")
    if err != nil {
        log.Fatalln("could not create table:", err)
    }
    _, err = db.Exec("CREATE TABLE slopedata (id INTEGER NOT NULL PRIMARY KEY, slope FLOAT64, tradingdays INTEGER);")
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
	if err != nil {
		log.Fatalf("error retrieving url: %s", err)
	}
	defer resp.Body.Close()
	csvReader := csv.NewReader(resp.Body)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("error reading csv: %s", err)
	}
	records = append(records[:0], records[0+1:]...)
	return records, nil
}


func getSlope(s string) {

      ntd := 5.00

	db, err := sql.Open("sqlite3", s+".db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("select sum(id) as sumx, sum(closeprice) as sumy, sum(id * closeprice) as sumxy, sum(id * id) as sumxx from(select id, closeprice from stockhistory order by ydate desc limit ?);", ntd)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var sumx float64
			var sumy float64
			var sumxy float64
			var sumxx float64
			rows.Scan(&sumx, &sumy, &sumxy, &sumxx)

			//fmt.Println(sumx, sumy, sumxy, sumxx)


			ntdsumxy := ntd * sumxy
			sumxsumy := sumx * sumy
			ntdsumxx := ntd * sumxx
			sumxsumx := sumx * sumx

			slope := (ntdsumxy - sumxsumy) / (ntdsumxx - sumxsumx)
		//if ((slope >= -0.001) && (slope <= 0.001)) {
			if ((slope >= -0.001) && (slope <= 0.1)) {
				fmt.Println(s,slope)
			}
			// else if (ntd <=35) {
			// ntd++
			// getSlope(s, ntd)
			// }
		}
		rows.Close()
}



