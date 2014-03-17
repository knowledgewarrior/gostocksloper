package main

import (
    "fmt"
  	"io/ioutil"
)


// var urls = []string{
//   ioutil.ReadFile("urls.txt")
// }

func main() {
	urls,_ := ioutil.ReadFile("urls.txt")
    fmt.Println(urls)
}


