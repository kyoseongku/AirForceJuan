package autodrone

import "math/rand"

// NewData ...
func NewData() DataType {
    var dataType = DataType{
        PropellerArray: make([]PropellerType, NumPropellers),
    }

    dataType.PropellerArray[0] = PropellerType{Frequency: 0.0}
    dataType.PropellerArray[1] = PropellerType{Frequency: 0.0}
    dataType.PropellerArray[2] = PropellerType{Frequency: 0.0}
    dataType.PropellerArray[3] = PropellerType{Frequency: 0.0}

    return dataType
}

// NewCompute ...
func NewCompute() ControlType {
    // allocate the required number of propellers to the pi control
    var controlType = ControlType{
        PropellerArray: make([]PropellerType, NumPropellers),
    }

    controlType.PropellerArray[0] = PropellerType{Frequency: 0.0}
    controlType.PropellerArray[1] = PropellerType{Frequency: 0.0}
    controlType.PropellerArray[2] = PropellerType{Frequency: 0.0}
    controlType.PropellerArray[3] = PropellerType{Frequency: 0.0}

    return controlType
}

// Compute ...
func ComputeControl(autoDrone DataType, control *ControlType) {
    // TODO actually compute based on current propeller states, altitude, and location.
    (*control).PropellerArray[0].Frequency = -1.0 * autoDrone.PropellerArray[0].Frequency
    (*control).PropellerArray[1].Frequency = -1.0 * autoDrone.PropellerArray[1].Frequency
    (*control).PropellerArray[2].Frequency = -1.0 * autoDrone.PropellerArray[2].Frequency
    (*control).PropellerArray[3].Frequency = -1.0 * autoDrone.PropellerArray[3].Frequency
}

// UpdatePropellerArray ...
func UpdatePropellerArray(autoDrone DataType) {
    // TODO: get actual propeller reading
    autoDrone.PropellerArray[0].Frequency = rand.Float64() * 10.0
    autoDrone.PropellerArray[1].Frequency = rand.Float64() * 10.0
    autoDrone.PropellerArray[2].Frequency = rand.Float64() * 10.0
    autoDrone.PropellerArray[3].Frequency = rand.Float64() * 10.0
}

// UpdateGPS ...
func UpdateGPS(autoDrone DataType) {
    // TODO: get actual GPS reading
    autoDrone.Altitude = rand.Float64() * 10.0
    autoDrone.Latitude = rand.Float64() * 10.0
    autoDrone.Longitude = rand.Float64() * 10.0
}
