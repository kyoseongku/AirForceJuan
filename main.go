package main

import (
  "log"
  "net/http"
  "os"
  "time"
)

var (
  Port = ":3210"
  Server = "localhost"
  Pi = "192.168.158.221"
  piData PiData       // Pi's current state
  piControl PiControl // Suggested new state for Pi
  T_Timeout = 3       // Seconds
  T_Pi = 1000         // Milliseconds
  N_Props = 4         // Number of propellers
)



func middlewares(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    log.Printf("%s %s %+v\n", req.Method, req.URL, req.Header["User-Agent"])

    fn(res, req)
  }
}



func main() {
  if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "pi") {
    log.Fatalln("Run: \"./AutoDrone server|pi\"\n")
  }
  mode := os.Args[1]

  if mode == "pi" {
    NewPi()

    t := time.NewTimer(time.Duration(T_Pi)*time.Millisecond)
    c := make(chan bool)

    log.Printf("Launching %s @ %s%s\n", mode, Pi, Port)
    for {
      select {
      case <-t.C:
        go DoThePi(c)
      case <-c:
        t = time.NewTimer(time.Duration(T_Pi)*time.Millisecond)
      }
    }
  } else {
    NewCompute()

    http.HandleFunc("/pi", middlewares(HandlePi))
    http.HandleFunc("/ui", middlewares(HandleUI))

    log.Printf("Launching %s @ %s%s\n", mode, Server, Port)
    err := http.ListenAndServe(Port, nil)
    log.Fatalln(err)
  }
}
