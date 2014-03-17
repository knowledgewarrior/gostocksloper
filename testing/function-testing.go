package main

import "fmt"

type inc func(digit int) int

func myInc(value int) int {
   return value +1
}

func getIncbynFunction(n int) inc {
   return func(value int) int {
             return value+n
   }
}

func main() {
        f := myInc
        g := getIncbynFunction
        h := g(4)
        i := g(6)
        fmt.Println(f(2))
        fmt.Println(h(5))
        fmt.Println(i(1))
}