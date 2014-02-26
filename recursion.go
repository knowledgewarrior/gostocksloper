package main

import "fmt"

func recv(value int, ch chan bool) {
    if value < 0 {
        ch <- true
        return
    }

    fmt.Println(value)
    go recv(value - 1, ch)
}

func main() {
    ch := make(chan bool)
    recv(10, ch)

    <-ch
}