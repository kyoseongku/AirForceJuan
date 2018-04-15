package main

import (
  "log"
)



func DoThePi(c chan bool) {
  log.Printf("Did the Pi to %s%s\n", Host, Port)

  c <-true
}

