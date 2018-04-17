package webserver

import (
	PI "AutoDrone/AD_PiServer"
	SC "AutoDrone/AD_ServerConstants"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	// PiData Web server's interpretation of pi data
	PiData PI.PiDataType

	// PiControl Web server's suggestion of pi's next state
	PiControl PI.PiControlType
)

// HandlePi ...
func HandlePi(res http.ResponseWriter, req *http.Request) {
	where := "pi"

	if req.Method != "POST" {
		http.NotFound(res, req)
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(reqBody, &PiData); err != nil {
		log.Fatalln(where, err)
	}

	log.Printf("%s Received data %+v\n", where, PI.PiData)

	Compute()

	data, err := json.Marshal(PiControl)
	if err != nil {
		log.Fatalln(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

// NewCompute ...
func NewCompute() PI.PiControlType {
	// allocate the required number of propellers to the pi control
	var piControl = PI.PiControlType{
		PropellerArray: make([]PI.PropellerType, SC.PiNumPropellers),
	}

	piControl.PropellerArray[0] = PI.PropellerType{Frequency: 0.0}
	piControl.PropellerArray[1] = PI.PropellerType{Frequency: 0.0}
	piControl.PropellerArray[2] = PI.PropellerType{Frequency: 0.0}
	piControl.PropellerArray[3] = PI.PropellerType{Frequency: 0.0}

	return piControl
}

// Compute ...
func Compute() {
	// TODO actually compute based on current propeller states, altitude, and location
	PiControl.PropellerArray[0].Frequency = -1.0 * PiData.PropellerArray[0].Frequency
	PiControl.PropellerArray[1].Frequency = -1.0 * PiData.PropellerArray[1].Frequency
	PiControl.PropellerArray[2].Frequency = -1.0 * PiData.PropellerArray[2].Frequency
	PiControl.PropellerArray[3].Frequency = -1.0 * PiData.PropellerArray[3].Frequency
}
