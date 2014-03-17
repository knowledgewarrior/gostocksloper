package main

import (
    "fmt"
    "time"
    "strings"
)

func main () {
	//t := time.Now()
	//fmt.Println(t.Format("2006-01-02"))
	t := time.Now().Format("2006-01-02")
	tArray := strings.Split(t, "-")
 	//fmt.Println(tArray)
 	y := tArray[0]
 	m := tArray[1]
 	d := tArray[2]
 	fmt.Println(y,m,d)
}