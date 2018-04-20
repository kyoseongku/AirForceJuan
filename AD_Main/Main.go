package main

import (
	piserver "AutoDrone/AD_PiServer"
	webserver "AutoDrone/AD_WebServer"
	"log"
	"os"
)

// Main i'd rather have the main function in a folder
func main() {
	if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "pi") {
		log.Fatalln("Run: \"./AutoDrone server|pi\"")
	}
	var mode = os.Args[1]

	// if the server is the raspberry pi, then perform these actions.
	// else perform these actions on the EC2.
	if mode == "pi" {
		piserver.StartPiServer()
	} else if mode == "server" {
		webserver.StartWebServer()
	}
}
