package autodrone

import(
    "github.com/eapache/queue"
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

// ---------- GPS Module Structs --------- //
type GPSReading struct {
    Altitude  float64
    Latitude  float64
    Longitude float64
    Timestamp string
    LatDirection byte
    LngDirection byte
}

// --------- PS Module Structs ---------- // 
// PSReading the distance is measured in millimeters.
type PSReading struct {
    Distance  float64
    Timestamp string
}

type PSensor struct {
    CurrReadings *queue.Queue
    rw_mutex    sync.RWMutex
    SensorName  string
    EchoPin     uint8
    TriggerPin  uint8
}

// --------- GA Module Structs --------- //
// Note: MPU6050 is big endian while RPi is little endian
type GAReading struct {
    GyroX  float64
    GyroY  float64
    GyroZ  float64
    AccelX float64
    AccelY float64
    AccelZ float64
}
