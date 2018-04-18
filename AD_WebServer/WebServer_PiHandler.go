package webserver

import (
	"AutoDrone/AD_AutoDrone"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	// AutoDrone Web server's interpretation of the autodrone
	AutoDrone autodrone.DataType

	// AutoDroneControl Web server's suggestion of the autodrone's next state
	AutoDroneControl autodrone.ControlType
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
	if err := json.Unmarshal(reqBody, &AutoDrone); err != nil {
		log.Fatalln(where, err)
	}

	log.Printf("%s Received data %+v\n", where, AutoDrone)

	autodrone.ComputeControl(AutoDrone, &AutoDroneControl)

	data, err := json.Marshal(AutoDroneControl)
	if err != nil {
		log.Fatalln(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}
