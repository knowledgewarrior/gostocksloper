package main

import "fmt"

type HelloFunc func(string)

func SayHello(to string) {
    fmt.Printf("Hello, %s!\n", to)
}

func main() {
    var hf HelloFunc

    hf = SayHello

    hf("world")
}