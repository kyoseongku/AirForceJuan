package main

import (
)



func NewCompute() {
    piControl = PiControl {
        PropellerArray: make( []Propeller, N_Propellers ),
    }

    piControl.PropellerArray[0] = Propeller{ Frequency: 0.0 }
    piControl.PropellerArray[1] = Propeller{ Frequency: 0.0 }
    piControl.PropellerArray[2] = Propeller{ Frequency: 0.0 }
    piControl.PropellerArray[3] = Propeller{ Frequency: 0.0 }
}



func Compute() {
    // TODO actually compute based on current propeller states, altitude, and location
    piControl.PropellerArray[0].Frequency = -1.0*piData.PropellerArray[0].Frequency
    piControl.PropellerArray[1].Frequency = -1.0*piData.PropellerArray[1].Frequency
    piControl.PropellerArray[2].Frequency = -1.0*piData.PropellerArray[2].Frequency
    piControl.PropellerArray[3].Frequency = -1.0*piData.PropellerArray[3].Frequency
}
