package autodrone

// DataType ...
type DataType struct {
    PropellerArray []PropellerType
    Altitude       float64
    Latitude       float64
    Longitude      float64
}

// ControlType ...
type ControlType struct {
    PropellerArray []PropellerType
}

// PropellerType ...
type PropellerType struct {
    Frequency float64
}

//
type GPSValueType struct {
    Altitude  float64
    Latitude  float64
    Longitude float64
}

