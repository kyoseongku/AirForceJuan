package main

import (
  "log"
  "net/http"
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



func phoneHandler(res http.ResponseWriter, req *http.Request) {
  where := "phone"
  log.Println(where, "hello")

  data := []byte("phone handler")
  res.Write(data)
}



func rpiHandler(res http.ResponseWriter, req *http.Request) {
  where := "rpi"
  log.Println(where, "hello")

  data := []byte("rpi handler")
  res.Write(data)
}



func webHandler(res http.ResponseWriter, req *http.Request) {
  where := "web"
  log.Println(where, "hello")

  data := []byte("web handler")
  res.Write(data)
}



func main() {
  http.HandleFunc("/phone", middlewares(phoneHandler))
  http.HandleFunc("/rpi", middlewares(rpiHandler))
  http.HandleFunc("/web", middlewares(webHandler))

  log.Println("Launching server", port)
  err := http.ListenAndServe(port, nil)
  log.Fatalln(err)
}
