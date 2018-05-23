package autodrone

const (
    // NumPropellers Number of Propellers
    NumPropellers = 4
)

// GPS Reader Constants
const (
    // BaudRate the stream flow of the serial port
    GPSBaudRate uint = 9600

    // SerialPort the port value ( directory ) that the serial stream can be read in ( CANNOT USE BLUETOOTH!! )
    GPSSerialPort = "/dev/ttyAMA0"

    // BufferSize when reading from the serial port, how much data should be pulled from the stream?
    GPSBufferSize uint = 256

    GPSDataBits uint = 8

    GPSStopBits uint = 1

    GPSMinBufferRead uint = 64
)

// Proximity Sensors Constants
const (
    // PSSpeedOfSound this variable is in millimeters per second
    PSSpeedOfSound float64 = 343000.0

    // PSTriggerPulse the pulse width of the trigger signal in microseconds
    PSTriggerPulse uint = 10

    // PSTimeout exit out of a measurement if no response comes from the sensor in milliseconds
    PSTimeout float64 = 50

    // PSMinRange the minimum range that can be calculated for the proximity sensor in milliseconds.
    PSMinRange float64 = 100

    // PSMaxRange the maximum range that can be calculated for the proximity sensor in milliseconds.
    PSMaxRange float64 = 4000

    // PSMaxReadBuffer the maximum amount of latest readings held for each sensor at any given time. 
    PSMaxReadBuffer int = 31
)

var (
    PSensorIDs []PSensor = []PSensor {
        PSensor {
            SensorName: "PS_Front",
            EchoPin:    21,
            TriggerPin: 20,
        },
        // Add more proximity sensors here!
        /*
        PSensor {
            SensorName: "PS_Back",
            EchoPin:    26,
            TriggerPin: 19,
        },
        */
    }
)
