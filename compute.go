package main

import (
)



func NewCompute() {
  piControl = PiControl{
    Props: make([]Propeller, N_Props),
  }

  piControl.Props[0] = Propeller{ Freq: 0.0 }
  piControl.Props[1] = Propeller{ Freq: 0.0 }
  piControl.Props[2] = Propeller{ Freq: 0.0 }
  piControl.Props[3] = Propeller{ Freq: 0.0 }
}



func Compute() {
  // TODO actually compute based on current propeller states, altitude, and location
  piControl.Props[0].Freq = -1.0*piData.Props[0].Freq
  piControl.Props[1].Freq = -1.0*piData.Props[1].Freq
  piControl.Props[2].Freq = -1.0*piData.Props[2].Freq
  piControl.Props[3].Freq = -1.0*piData.Props[3].Freq
}
