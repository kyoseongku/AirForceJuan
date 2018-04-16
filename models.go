package main

type PiData struct {
  Props []Propeller
  Alt   float64
  Lat   float64
  Lng   float64
}

type PiControl struct {
  Props []Propeller
}

type Propeller struct {
  Which string
  Freq  float64
}
