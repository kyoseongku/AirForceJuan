package webserver

import (
	PI "AutoDrone/AD_PiServer"
	SC "AutoDrone/AD_ServerConstants"
	"log"
	"net/http"
)

// StartWebServer ...
func StartWebServer() {
	PiControl = NewCompute()
	PiData = PI.NewPi()
	http.HandleFunc("/pi", middlewares(HandlePi))
	http.HandleFunc("/ui", middlewares(HandleUI))

	log.Printf("Launching web server @ %s%s\n", SC.WebIPAddress, SC.WebPort)
	err := http.ListenAndServe(SC.WebPort, nil)
	log.Fatalln(err)
}

func middlewares(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s %+v\n", req.Method, req.URL, req.Header["User-Agent"])
		fn(res, req)
	}
}
