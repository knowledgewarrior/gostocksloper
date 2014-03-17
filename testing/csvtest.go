package main

import (
    "encoding/csv"
    "fmt"
    "os"
)
func checkError(e error){
    if e != nil {
        panic(e)
    }
}
func writeCSV(){
    fmt.Println("Writing csv")
    f, err := os.Create("./test.csv")
    checkError(err)
    defer f.Close()

    w := csv.NewWriter(f)
    s := "Cr@zy text with , and \\ and \" etc"
    record := []string{ 
      "Unquoted string",
      s,
    }
    fmt.Println(record)
    w.Write(record)

    record = []string{ 
      "Quoted string",
      fmt.Sprintf("%q",s),
    }
    fmt.Println(record)
    w.Write(record)
    w.Flush()
}
func readCSV(){
    fmt.Println("Reading csv")
    file, err := os.Open("./test.csv")
    defer file.Close();
    cr := csv.NewReader(file)
    records, err := cr.ReadAll()
    checkError(err)
    for _, record := range records {
        fmt.Println(record)
    }
}
func main() {
   writeCSV()
   readCSV()
}
