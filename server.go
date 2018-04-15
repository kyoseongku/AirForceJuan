package main

import (
  "log"
  "net/http"
)



func HandlePhone(res http.ResponseWriter, req *http.Request) {
  where := "phone"
  log.Println(where, "hello")

  data := []byte("phone handler")
  res.Write(data)
}



func HandlePi(res http.ResponseWriter, req *http.Request) {
  where := "pi"
  log.Println(where, "hello")

  data := []byte("pi handler")
  res.Write(data)
}



func HandleWeb(res http.ResponseWriter, req *http.Request) {
  where := "web"
  log.Println(where, "hello")

  data := []byte("web handler")
  res.Write(data)
}

