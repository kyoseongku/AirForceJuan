package main

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "strings"
    "time"
)



func updateGPS() {
    // TODO get actual GPS reading
    piData.Altitude  = rand.Float64()*10.0
    piData.Latitude  = rand.Float64()*10.0
    piData.Longitude = rand.Float64()*10.0
}



func updatePropellerArray() {
    // TODO get actual propeller reading
    piData.PropellerArray[0].Frequency = rand.Float64()*10.0
    piData.PropellerArray[1].Frequency = rand.Float64()*10.0
    piData.PropellerArray[2].Frequency = rand.Float64()*10.0
    piData.PropellerArray[3].Frequency = rand.Float64()*10.0
}



func NewPi() {
    piData = PiData{
        PropellerArray: make([]Propeller, N_Propellers ),
    }

    piData.PropellerArray[0] = Propeller{ Frequency: 0.0 }
    piData.PropellerArray[1] = Propeller{ Frequency: 0.0 }
    piData.PropellerArray[2] = Propeller{ Frequency: 0.0 }
    piData.PropellerArray[3] = Propeller{ Frequency: 0.0 }
}



func DoThePi(c chan bool) {
    log.Printf("Sending data to %s%s from %s\n", Server, Port, Pi)

    // Update readings
    updateGPS()
    updatePropellerArray()

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
        errSplit := strings.Split(err.Error(), " ")

        if len(errSplit) == 8 && errSplit[6]+errSplit[7] == "connectionrefused" {
            log.Println("Can't reach server: connection refused")
            c <-true
            return
        } else if len(errSplit) == 10 && errSplit[3]+errSplit[4] == "requestcanceled" {
            log.Println("Can't reach server: connection timed out")
            c <-true
            return
        }

        log.Fatalln(err)
    }
    defer response.Body.Close()

    respBody, err := ioutil.ReadAll(response.Body)

    log.Println(string(respBody))

    c <-true
}
