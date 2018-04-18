package piserver

import (
	"AutoDrone/AD_AutoDrone"
	SC "AutoDrone/AD_ServerConstants"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	// AutoDrone AutoDrone's current state
	AutoDrone autodrone.DataType

	// AutoDroneControl Suggested new state for Autodrone
	AutoDroneControl autodrone.ControlType
)

// StartPiServer ...
func StartPiServer() {
	AutoDrone = autodrone.NewData()
	// polling variables
	var timer = time.NewTimer(time.Duration(SC.PiPollPeriod) * time.Millisecond)
	var channel = make(chan bool)

	log.Printf("Launching pi @ %s%s\n", SC.PiIPAddress, SC.PiPort)

	// not sure what this does yet
	// Pi while looper.
	for {
		select {
		case <-timer.C:
			go DoThePi(channel)
		case <-channel:
			timer = time.NewTimer(time.Duration(SC.PiPollPeriod) * time.Millisecond)
		}
	}
}

// DoThePi ...
func DoThePi(c chan bool) {
	log.Printf("Sending data to %s%s from %s\n", SC.WebIPAddress, SC.WebPort, SC.PiIPAddress)

	// Update readings
	autodrone.UpdateGPS(AutoDrone)
	autodrone.UpdatePropellerArray(AutoDrone)

	body, err := json.Marshal(AutoDrone)
	if err != nil {
		log.Fatalln(err)
	}

	var client = &http.Client{
		Timeout: time.Duration(SC.PiTimeout) * time.Second,
	}

	request, err := http.NewRequest("POST", "http://"+SC.WebIPAddress+SC.WebPort+"/pi", bytes.NewBuffer(body))
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
