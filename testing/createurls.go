package main

import (
    "fmt"
  	"bufio"
	"os"
	"strings"
  	"strconv"
  	"time"
)

func check(e error) {
    if e != nil {
      fmt.Println(e)
    }
} //check

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


func main() {
	fname := "urls.txt"
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	check(err)
	defer f.Close()
    symbols, err := readLines("symbols-medium.txt")
    check(err)
    t := time.Now().Format("2006-01-02")
	tArray := strings.Split(t, "-")
	nowyear := tArray[0]
	nmonth := tArray[1]
	nmonth2, err := strconv.Atoi(nmonth)
	
	  nowmonth := (nmonth2 - 1)
	  nowday := tArray[2]

	  p := time.Now().AddDate(0, 0, -729).Format("2006-01-02")
	  pArray := strings.Split(p, "-")
	  thenyear := pArray[0]
	  tmonth := pArray[1]
	  tmonth2, err := strconv.Atoi(tmonth)
	 
	  thenmonth := (tmonth2 - 1)
	  thenday := pArray[2]
	  for _, symbol := range symbols {
	  url := fmt.Sprintf("http://ichart.finance.yahoo.com/table.csv?s=%s&a=%d&b=%s&c=%s&d=%d&e=%s&f=%s&g=d", symbol, thenmonth, thenday, thenyear, nowmonth, nowday, nowyear)
	  fmt.Println(url)
	  b := bufio.NewWriter(f)
	  defer func() {
	  if err = b.Flush(); err != nil {
	        //fmt.Println(err)
	  }
	  }()
	  fmt.Fprint(b, url)
	  fmt.Fprint(b, "\n")
	  }
}
