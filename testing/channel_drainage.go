package main

import (
 "log"
 "os"
 "os/signal"
 "reflect"
 "syscall"
 "time"
)

var (
 BufferSize = 512
 MaxIter    = 10
 monitored  []interface{}
 c          = make(chan int, BufferSize)
 stopping   bool
)

func RegisterChannel(i interface{}) {
 monitored = append(monitored, i)
}

func MonitorSigTerm() chan bool {
 s := make(chan os.Signal, 1)
 b := make(chan bool)
 signal.Notify(s, syscall.SIGTERM)

  go func(c chan os.Signal, b chan bool) {
  _ = <-c
  log.Println("Cleaning up")
  // tell the caller
  b <- true
  for _, i := range monitored {
   ch := reflect.ValueOf(i)
   if ch.Kind() != reflect.Chan {
    continue
   }
   prev := 0
   iteration := 0
   for {
    if ch.Len() == 0 {
     break
    }

     if prev == ch.Len() {
     iteration++
     // enough?
     if iteration >= MaxIter {
      log.Println("Dropping")
      break
     }
    } else {
     iteration = 0
    }

     prev = ch.Len()
    log.Printf("Draining:%v\n", prev)
    // other goroutines are working, let them
    time.Sleep(1e9)
   }
  }
  os.Exit(1)
 }(s, b)
 return b
}

func main() {
 RegisterChannel(c)
 stop := MonitorSigTerm()

  go func() {
  i := 0
  for {
   if stopping {
    break
   }
   i++
   c <- i
   time.Sleep(1e9)
  }
 }()

  go func() {
  for {
   i := <-c
   log.Printf("rx:%v\n", i)
   // slower read
   time.Sleep(2e9)
  }
 }()

  stopping = <-stop

  // wait for cleanup to finish
 select {}

}