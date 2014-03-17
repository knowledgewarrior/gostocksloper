package main 

import ("fmt";"io";"io/ioutil";"os") 

const file = "temp.txt" 

func write(flag int, text string) { 
        f, err:=os.Open(file) 
        if err != nil { fmt.Println(err); return } 
        n, err := io.WriteString(f, text) 
        if err != nil { fmt.Println(n, err); return } 
        f.Close() 
        data, err := ioutil.ReadFile(file) 
        if err != nil { fmt.Println(err); return } 
        fmt.Println(string(data)) 
} 

func main() { 
        write(os.O_CREATE|os.O_TRUNC|os.O_RDWR, "new") 
        for i := 0; i < 2; i++ { 
                write(os.O_APPEND|os.O_RDWR, "|append") 
        } 
} 
