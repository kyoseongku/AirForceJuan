package main

import (
	PI "AutoDrone/AD_PiServer"
	WEB "AutoDrone/AD_WebServer"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 || (os.Args[1] != "server" && os.Args[1] != "pi") {
		log.Fatalln("Run: \"./AutoDrone server|pi\"")
	}
	var mode = os.Args[1]

	// if the server is the raspberry pi, then perform these actions.
	// else perform these actions on the EC2.
	if mode == "pi" {
		PI.StartPiServer()
	} else {
		WEB.StartWebServer()
	}
}
