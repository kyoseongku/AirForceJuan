package piserver

import (
	SC "AutoDrone/AD_ServerConstants"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	// PiData Pi's current state
	PiData PiDataType

	// PiControl Suggested new state for Pi
	PiControl PiControlType
)

// StartPiServer ...
func StartPiServer() {
	PiData = NewPi()
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
	updateGPS()
	updatePropellerArray()

	body, err := json.Marshal(PiData)
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

func updateGPS() {
	// TODO: get actual GPS reading
	PiData.Altitude = rand.Float64() * 10.0
	PiData.Latitude = rand.Float64() * 10.0
	PiData.Longitude = rand.Float64() * 10.0
}

func updatePropellerArray() {
	// TODO: get actual propeller reading
	PiData.PropellerArray[0].Frequency = rand.Float64() * 10.0
	PiData.PropellerArray[1].Frequency = rand.Float64() * 10.0
	PiData.PropellerArray[2].Frequency = rand.Float64() * 10.0
	PiData.PropellerArray[3].Frequency = rand.Float64() * 10.0
}

// NewPi ...
func NewPi() PiDataType {
	var piData = PiDataType{
		PropellerArray: make([]PropellerType, SC.PiNumPropellers),
	}

	piData.PropellerArray[0] = PropellerType{Frequency: 0.0}
	piData.PropellerArray[1] = PropellerType{Frequency: 0.0}
	piData.PropellerArray[2] = PropellerType{Frequency: 0.0}
	piData.PropellerArray[3] = PropellerType{Frequency: 0.0}

	return piData
}
