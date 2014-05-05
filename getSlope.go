package main

import (
  "fmt"
  "bufio"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "path/filepath"
  "os"
  "io"
  "log"
  "strings"
  "strconv"
)

var ERROR *log.Logger

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

func getSlope(symbol string, ntd float64, slope float64) {
  if (slope < 0.001) && (slope > -0.001) {

  	fname := "Slopes.csv"
    f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil { log.Println(err) }
    defer f.Close()

    b := bufio.NewWriter(f)
    defer func() {
      if err = b.Flush(); err != nil { log.Println(err) }
    }()
    slope := slope * -1.00
    fmt.Fprint(b, symbol)
    fmt.Fprint(b, ",")
    strconv.FormatFloat(slope, 'f', 32, 64)
    fmt.Fprint(b, slope)
    fmt.Fprint(b, "\n")

    return
  }

  db, err := sql.Open("sqlite3", "db/"+symbol)
  if err != nil { log.Println(err) }
  defer db.Close()

  var sumx float64
  var sumy float64
  var sumxy float64
  var sumxx float64

  stmt, err := db.Prepare("select sum(id) as sumx, sum(closeprice) as sumy, " +
  "sum(id * closeprice) as sumxy, sum(id * id) as sumxx from(select id, " +
  "closeprice from stockhistory order by ydate desc limit ?);")
  if err != nil { log.Println(err) }
  defer stmt.Close()
  rows, err := stmt.Query(ntd)
  if err != nil { log.Println(err) }
  defer rows.Close()
  for rows.Next() {

    ntdsumxy := ntd * sumxy
    sumxsumy := sumx * sumy
    ntdsumxx := ntd * sumxx
    sumxsumx := sumx * sumx
    slope := (ntdsumxy - sumxsumy) / (ntdsumxx - sumxsumx)

    go getSlope(symbol, ntd + 1.00, slope)
  } // for rows next
} //getSlope

func walkFiles(location string) (chan string) {
    chann := make(chan string)
    go func(){
        filepath.Walk(location, func(path string, _ os.FileInfo, _ error)(err error){
            chann <- path
            return
        })
        defer close(chann)
    }()
        return chann
}

func main() {
  Init(os.Stderr)
  os.Remove("Slopes.csv")
  os.Remove("baseseeker_log.txt")
  logf, err := os.OpenFile("baseseeker_log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  if err != nil { fmt.Println(err) }
  defer logf.Close()
  log.SetOutput(logf)

	dbdir := "db"
  chann := walkFiles(dbdir)
  for symbols := range chann {
    if symbols == "db" { continue }
    symbol := strings.TrimLeft(symbols, "db/")

    getSlope(symbol, 120.00, 1.00)

  }

} // func main
