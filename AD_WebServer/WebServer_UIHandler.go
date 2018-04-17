package webserver

import (
	"log"
	"net/http"
)

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
