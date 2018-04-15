package main

import (
  "log"
  "net/http"
  "os"
  "time"
)

var (
  Port = ":3210"
  Host = "192.168.158.221"
  T_Pi = 1000
)



func middlewares(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    log.Printf("%s %s %+v\n", req.Method, req.URL, req.Header["User-Agent"])

    fn(res, req)
  }
}



func main() {
  if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "pi") {
    log.Fatalln("Run: \"./AirForceJuan server|pi\"\n")
  }
  mode := os.Args[1]

  if mode == "pi" {
    t := time.NewTimer(time.Duration(T_Pi)*time.Millisecond)
    c := make(chan bool)
    for {
      select {
      case <-t.C:
        go DoThePi(c)
      case <-c:
        t = time.NewTimer(time.Duration(T_Pi)*time.Millisecond)
      }
    }
  } else {
    http.HandleFunc("/phone", middlewares(HandlePhone))
    http.HandleFunc("/pi", middlewares(HandlePi))
    http.HandleFunc("/web", middlewares(HandleWeb))

    log.Printf("Launching %s%s\n", mode, Port)
    err := http.ListenAndServe(Port, nil)
    log.Fatalln(err)
  }
}

