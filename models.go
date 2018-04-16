package main

type PiData struct {
    PropellerArray[] Propeller
    Altitude  float64
    Latitude  float64
    Longitude float64
}

type PiControl struct {
    PropellerArray[] Propeller
}

type Propeller struct {
    Frequency float64
}
