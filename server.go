package main

import (
  "io/ioutil"
  "log"
  "net/http"
)



func HandlePi(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.NotFound(res, req)
    return
  }

  where := "pi"
  log.Println(where, "Received data")

  reqBody, err := ioutil.ReadAll(req.Body)
  if err != nil {
    log.Fatalln(err)
  }

  log.Println(string(reqBody))

  res.Header().Set("Content-Type", "application/json")

  data := []byte("pi handler\n")
  res.Write(data)
}



func HandleUI(res http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" {
    http.NotFound(res, req)
    return
  }

  where := "ui"
  log.Println(where, "hello")

  data := []byte("web handler\n")
  res.Write(data)
}
