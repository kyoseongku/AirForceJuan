package main

import (
  "bytes"
  "encoding/json"
  "io/ioutil"
  "log"
  "math/rand"
  "net/http"
  "time"
)

var piData PiData



func updateGPS() {
  // TODO get actual GPS reading
  piData.Alt = rand.Float64()*10.0
  piData.Lat = rand.Float64()*10.0
  piData.Lng = rand.Float64()*10.0
}



func updateProps() {
  // TODO get actual propeller reading
  piData.Props[0].Freq = rand.Float64()*10.0
  piData.Props[1].Freq = rand.Float64()*10.0
  piData.Props[2].Freq = rand.Float64()*10.0
  piData.Props[3].Freq = rand.Float64()*10.0
}



func NewPi() {
  piData = PiData{
    Props: make([]Propeller, N_Props),
  }

  piData.Props[0] = Propeller{ Freq: 0.0, Which: "Front-Left" }
  piData.Props[1] = Propeller{ Freq: 0.0, Which: "Front-Right" }
  piData.Props[2] = Propeller{ Freq: 0.0, Which: "Rear-Right" }
  piData.Props[3] = Propeller{ Freq: 0.0, Which: "Rear-Left" }
}



func DoThePi(c chan bool) {
  log.Printf("Sending data to %s%s from %s\n", Server, Port, Pi)

  // Update readings
  updateGPS()
  updateProps()

  body, err := json.Marshal(piData)
  if err != nil {
    log.Fatalln(err)
  }

  client := &http.Client{
    Timeout: time.Duration(T_Timeout)*time.Second,
  }

  request, err := http.NewRequest("POST", "http://"+Server+Port+"/pi", bytes.NewBuffer(body))
  if err != nil {
    log.Fatalln(err)
  }

  request.Header.Set("Content-Type", "application/json")
  response, err := client.Do(request)
  if err != nil {
    log.Fatalln(err)
  }
  defer response.Body.Close()

  respBody, err := ioutil.ReadAll(response.Body)

  log.Println(string(respBody))

  c <-true
}
