package main

import (
  "log"
  "net/http"

  "AirForceJuan/handlers"
)

var (
  port = ":3210"
)



func middlewares(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    log.Printf("%s %s %+v\n", req.Method, req.URL, req.Header["User-Agent"])

    fn(res, req)
  }
}



func main() {
  http.HandleFunc("/api/phone", middlewares(handlers.Phone))
  http.HandleFunc("/api/rpi", middlewares(handlers.Rpi))
  http.HandleFunc("/api/web", middlewares(handlers.Web))

  log.Println("Launching server", port)
  err := http.ListenAndServe(port, nil)
  log.Fatalln(err)
}