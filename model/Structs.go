package autodrone

import(
    "sync"
)
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

// ---------- GPS Module Structs ----------
type GPSReading struct {
    Altitude  float64
    Latitude  float64
    Longitude float64
    Timestamp string
    LatDirection byte
    LngDirection byte
}

// --------- PS Module Structs ----------
// PSReading the distance is measured in centimeters.
type PSReading struct {
    Distance  float64
    Timestamp string
}

type PSensor struct {
    CurrReading PSReading
    SensorName  string
    EchoPin     uint8
    TriggerPin  uint8
    rw_mutex    sync.RWMutex
}
