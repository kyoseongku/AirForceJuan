package main

import (
	"AutoDrone/model"
	"log"
	"os"
)

// AutoDrone and AutoDroneControl global variables are used differently between PiServer and WebServer.
// Webserver uses these global variables as a temporary state storage for calculation purposes.
// PiServer uses these global variables as state machine.
var (
	// AutoDrone AutoDrone's current state
	AutoDrone autodrone.DataType

	// AutoDroneControl Suggested new state for Autodrone
	AutoDroneControl autodrone.ControlType
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "pi") {
		log.Fatalln("Run: \"./AutoDrone server|pi\"")
	}
	var mode = os.Args[1]

	// directed to the piserver package and the webserver package
	if mode == "pi" {
		StartPiServer()
	} else if mode == "server" {
		StartWebServer()
	}
}
