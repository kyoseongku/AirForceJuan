package main

import (
    "github.com/stianeikeland/go-rpio"
	"AutoDrone/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// StartPiServer ...
func StartPiServer() {
	AutoDrone = autodrone.NewData()

    // initialize the GPIO drivers in rpio
    err := rpio.Open()
    if err != nil { log.Fatalln( err ) }
    defer rpio.Close()

    // Start modules
	//go autodrone.GPS_StartModule()
    //go autodrone.PS_StartModule()
    go autodrone.GA_StartModule()

    // polling variables
	var timer = time.NewTimer(time.Duration(PiPollPeriod) * time.Millisecond)
	var channel = make(chan bool)

	log.Printf("Launching pi @ %s%s\n", PiIPAddress, PiPort)

    for {
        select{ }
    }

    for {
		select {
		case <-timer.C:
			go DoThePi(channel)
		case <-channel:
			timer = time.NewTimer(time.Duration(PiPollPeriod) * time.Millisecond)
		}
	}
}


// DoThePi ...
func DoThePi(c chan bool) {
	log.Printf("Sending data to %s%s from %s\n", WebIPAddress, WebPort, PiIPAddress)

	// Update readings
	autodrone.UpdateGPS( &AutoDrone )
	autodrone.UpdatePropellerArray( &AutoDrone )
    log.Printf( "\n A: %f\n Lat: %f\n Lng: %f\n\n", AutoDrone.Altitude, AutoDrone.Latitude, AutoDrone.Longitude )

	body, err := json.Marshal(AutoDrone)
	if err != nil {
		log.Fatalln(err)
	}

	var client = &http.Client{
		Timeout: time.Duration(PiTimeout) * time.Second,
	}

	request, err := http.NewRequest("POST", "http://"+WebIPAddress+WebPort+"/pi", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		errSplit := strings.Split(err.Error(), " ")
		if len(errSplit) == 8 && errSplit[6]+errSplit[7] == "connectionrefused" {
			log.Println("Can't reach server: connection refused")
			c <- true
			return
		} else if len(errSplit) == 10 && errSplit[3]+errSplit[4] == "requestcanceled" {
			log.Println("Can't reach server: connection timed out")
			c <- true
			return
		}
		log.Fatalln(err)
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)

	log.Println(string(respBody))

	c <- true
}
