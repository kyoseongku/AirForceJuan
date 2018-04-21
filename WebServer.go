package main

import (
	"AutoDrone/model"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// StartWebServer ...
func StartWebServer() {
	AutoDroneControl = autodrone.NewCompute()
	AutoDrone = autodrone.NewData()
	http.HandleFunc("/pi", middlewares(HandlePi))
	http.HandleFunc("/ui", middlewares(HandleUI))

	log.Printf("Launching web server @ %s%s\n", WebIPAddress, WebPort)
	err := http.ListenAndServe(WebPort, nil)
	log.Fatalln(err)
}

func middlewares(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %+v\n", req.Method, req.URL, req.Header["User-Agent"])
		fn(res, req)
	}
}

// HandleUI ...
func HandleUI(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.NotFound(res, req)
		return
	}

	where := "ui"
	log.Println(where, "hello")

	data := []byte("web handler\n")
	res.Write(data)
}

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
