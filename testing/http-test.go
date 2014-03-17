/*
Jason Fowler February 2014
Linear Regression Analysis of Yahoo Stock Data using formula:
Slope(b) = (N*ΣXY - (ΣX)(ΣY)) / (N*ΣX2 - (ΣX)2)
*/

package main

import (
  "fmt"
  "bufio"
  _ "github.com/mattn/go-sqlite3"
  "os"
  "io"
  "io/ioutil"
  "net/http"
  "time"
  "strings"
  "log"
  "strconv"
)

var ERROR *log.Logger

type HttpResponse struct {
  Url      string
  ByteStr  []byte
  Response *http.Response
  Err      error
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  check(err)
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

func check(e error) {
    if e != nil {
      log.Println(e)
    }
} //check

func Get(url string) (chan *HttpResponse) {

  channel  := make(chan *HttpResponse)
  client   := &http.Client{}
  req, _ := http.NewRequest("GET", url, nil)
  resp, _ := client.Do(req)
  fmt.Println(resp.StatusCode)
   if resp.StatusCode > 399 {
    os.Exit(1)
  }


  go func(){
      resp, _ := client.Do(req)
      
      defer resp.Body.Close()

      bs, err := ioutil.ReadAll(resp.Body)
      check(err)

      channel <- &HttpResponse{url, bs, resp, err}
  }()

  return channel
}

func getStocks(symbol string) {

  t := time.Now().Format("2006-01-02")
  tArray := strings.Split(t, "-")

  nowyear := tArray[0]
  nmonth := tArray[1]
  nmonth2, err := strconv.Atoi(nmonth)
  check(err)
  nowmonth := (nmonth2 - 1)
  nowday := tArray[2]

  p := time.Now().AddDate(0, 0, -29).Format("2006-01-02")
  pArray := strings.Split(p, "-")
  thenyear := pArray[0]
  tmonth := pArray[1]
  tmonth2, err := strconv.Atoi(tmonth)
  check(err)
  thenmonth := (tmonth2 - 1)
  thenday := pArray[2]

  url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%s&c=%s&d=%d&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
  //fmt.Println(url)
  channel := Get(url)
  //resp := <-channel
  _ = <-channel
  //fmt.Println(string(resp.ByteStr))

 
  // csvReader := csv.NewReader(resp.Body)
  // records, err := csvReader.ReadAll()
  // check(err)
  //  // if not >120 lines, skip
  // lineCount := 0
  // for _ = range records {
  //   lineCount += 1
  // }
  // if lineCount < 121 {
  //   return
  // }

  // // if no data, skip - check closing price
  // for _, record := range records {
  //   c := record[4]
  //   if c == "0.00" {
  //     return
  //   }
  // }

} // getStocks


func main() {
  Init(os.Stderr)
  //os.Remove("Slopes.csv")
  os.Remove("http_log.txt")
  //os.RemoveAll("./db")
  //os.Mkdir("./db", 0700)
  logf, err := os.OpenFile("http_log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
  check(err)
  defer logf.Close()
  log.SetOutput(logf)
  symbols, err := readLines("symbols-small.txt")
  check(err)
  for _, symbol := range symbols {
    getStocks(symbol)
  }

} // func main
