package main

import (
    "fmt"
)

func main() {
    arr := [5]float64{1,2,3,4,5}
	x := arr[0:5]
    fmt.Println(x)
    x = append(x[:0], x[0+1:]...) 
    fmt.Println(x)

}