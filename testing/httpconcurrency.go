package main
 
import (
	"fmt"
	"net/http"
	"time"
)
 
var urls = []string{
	"http://apple.ca",
	"http://golang.org",
	"http://google.ca",
}
 
type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}
 
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
			fmt.Printf("%s was fetched\n", r.url)
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