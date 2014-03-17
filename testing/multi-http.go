package main
 
import (
	"fmt"
	"net/http"
	"bufio"
	"os"
	"strings"
  	"strconv"
  	"time"
)
 
var urls = []string{
	"http://ichart.finance.yahoo.com/table.csv?s=GOOG&a=00&b=18&c=2014&d=01&e=21&f=2014&g=d",
	"http://ichart.finance.yahoo.com/table.csv?s=AAPL&a=00&b=18&c=2014&d=01&e=21&f=2014&g=d",
	"http://ichart.finance.yahoo.com/table.csv?s=F&a=00&b=18&c=2014&d=01&e=21&f=2014&g=d",
}
 
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

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

 
func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse, len(urls)) // buffered
	responses := []*HttpResponse{}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := http.Get(url)
			ch <- &HttpResponse{url, resp, err}
		}(url)
	}
 
	for {
		select {
		case r := <-ch:
			//fmt.Printf("%s was fetched\n", r.url)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		default:
			fmt.Printf(".")
			time.Sleep(5e7)
		}
	}
	return responses
 
}
 
func main() {
	results := asyncHttpGets(urls)
	for _, result := range results {
		fmt.Printf("%s status: %s\n", result.url, result.response.Status)
	}
}