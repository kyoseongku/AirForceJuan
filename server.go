package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
)



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
    if err := json.Unmarshal(reqBody, &piData); err != nil {
        log.Fatalln(where, err)
    }

    log.Printf("%s Received data %+v\n", where, piData)

    Compute()

    data, err := json.Marshal(piControl)
    if err != nil {
        log.Fatalln(err)
    }

    res.Header().Set("Content-Type", "application/json")
    res.Write(data)
}



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
