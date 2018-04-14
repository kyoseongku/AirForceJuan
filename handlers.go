package handlers

import (
  "log"
  "net/http"
)



func Phone(http.ResponseWriter, *http.Request) {
  where := "web"
  log.Prinln(where, "hello")
}



func Rpi(http.ResponseWriter, *http.Request) {
  where := "web"
  log.Prinln(where, "hello")
}



func Web(http.ResponseWriter, *http.Request) {
  where := "web"
  log.Prinln(where, "hello")
}
