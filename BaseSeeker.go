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
  "github.com/mreiferson/go-httpclient"
  "os"
  "io"
  "math/rand"
  "encoding/csv"
  "net/http"
  "time"
  "strings"
  "log"
  "strconv"
)

var ERROR *log.Logger

func getSymbols(file string) {
  symbols, err := readLines(file)
  if err != nil { log.Println(err) }

  for _, symbol := range symbols {
   getStocks(symbol)
  }
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil { log.Println(err) }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
} //readlines

func Init(errorHandle io.Writer) {
  ERROR = log.New(errorHandle,
  "ERROR: ",
  log.Ldate|log.Ltime|log.Lshortfile)
}

func getStocks(symbol string) {

  t := time.Now().Format("2006-01-02")
  tArray := strings.Split(t, "-")

  nowyear := tArray[0]
  nmonth := tArray[1]
  nmonth2, err := strconv.Atoi(nmonth)
  if err != nil { log.Println(err) }
  nowmonth := (nmonth2 - 1)
  nowday := tArray[2]

  p := time.Now().AddDate(0, 0, -729).Format("2006-01-02")
  pArray := strings.Split(p, "-")
  thenyear := pArray[0]
  tmonth := pArray[1]
  tmonth2, err := strconv.Atoi(tmonth)
  if err != nil { log.Println(err) }
  thenmonth := (tmonth2 - 1)
  thenday := pArray[2]

  url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%s&c=%s&d=%d&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)

  transport := &httpclient.Transport{
    ConnectTimeout:        1*time.Second,
    RequestTimeout:        5*time.Second,
    ResponseHeaderTimeout: 2*time.Second,
  }
  defer transport.Close()

  client := &http.Client{Transport: transport}
  req, _ := http.NewRequest("GET", url, nil)
  resp, err := client.Do(req)
      if err != nil {
         log.Println(err)
        return
    }else {
      defer resp.Body.Close()

      if resp.StatusCode > 399 {
        resp.Body.Close()
        return
      }

      csvReader := csv.NewReader(resp.Body)
      records, err := csvReader.ReadAll()
      if err != nil {
        log.Println(err)
        return
      }
       // if not >120 lines, skip
      lineCount := 0
      for _ = range records {
        lineCount += 1
      }
      if lineCount < 121 {
        return
      }

      // if no data, skip - check closing price
      for _, record := range records {
        c := record[4]
        if c == "0.00" {
          return
        }
      }

    // // okay, good to go
      records = append(records[:0], records[0+1:]...)
      db, err := sql.Open("sqlite3", "db/"+symbol+".db")
      if err != nil { log.Println(err) }
      defer db.Close()
      _, err = db.Exec("CREATE TABLE stockhistory (id INTEGER NOT NULL PRIMARY KEY, ydate TEXT, closeprice INTEGER);")
      if err != nil { log.Println(err) }

      for _, record := range records {
        d := record[0]
        c := record[4]
        tx, err := db.Begin()
        if err != nil { log.Println(err) }
        insert_stmt, err := tx.Prepare("insert into stockhistory(ydate,closeprice) values(?,?)")
        if err != nil { log.Println(err) }
        defer insert_stmt.Close()
        _, err = insert_stmt.Exec(d,c)
        if err != nil { log.Println(err) }
        tx.Commit()
        insert_stmt.Close()
      }
  amt := time.Duration(rand.Intn(250))
  time.Sleep(time.Millisecond * amt)
  } //else
} // getStocks

func getSlope(symbol string, ntd float64, slope float64, ch chan bool) {
  if (slope < 0.01) && (slope > -0.01) {
    fname := "Slopes.csv"
    f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil { log.Println(err) }
    defer f.Close()

    b := bufio.NewWriter(f)
    defer func() {
    if err = b.Flush(); err != nil {
          //fmt.Println(err)
    }
    }()
    slope := slope * -1
    fmt.Fprint(b, symbol+",")
    fmt.Fprint(b, slope)
    fmt.Fprint(b, "\n")
    ch <- true
    return
}

  db, err := sql.Open("sqlite3", "db/"+symbol+"db")
  if err != nil { log.Println(err) }
  defer db.Close()
  rows, err := db.Query("select sum(id) as sumx, sum(closeprice) as sumy, sum(id * closeprice) as sumxy, sum(id * id) as sumxx from(select id, closeprice from stockhistory order by ydate desc limit ?);", ntd)
  if err != nil { log.Println(err) }
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
    go getSlope(symbol, ntd + 1.00, slope, ch)
  } // for rows
  rows.Close()
} //getSlope

func main() {
  Init(os.Stderr)
  os.Remove("Slopes.csv")
  os.Remove("baseseeker_log.txt")
  os.RemoveAll("./db")
  os.Mkdir("./db", 0700)
  logf, err := os.OpenFile("baseseeker_log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil { fmt.Println(err) }
  defer logf.Close()
  log.SetOutput(logf)

getSymbols("symbols1.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols2.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols3.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols4.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols5.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols6.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols7.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols8.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols9.txt")
time.Sleep(1800 * time.Second)
getSymbols("symbols9.txt")
time.Sleep(1800 * time.Second)

  for _, symbol := range symbols {
    _, err := os.Stat("db/"+symbol+".db")
    if err != nil {
    break
    }
    ch1 := make(chan bool)
    getSlope(symbol, 120.00, 2.00, ch1)
    <-ch1
  }

} // func main