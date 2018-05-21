package autodrone

var (
    // NumPropellers Number of Propellers
    NumPropellers = 4
)

// GPS Reader Constants
var (
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
var (
    // PSSpeedOfSound this variable is in centimeters per second
    PSSpeedOfSound float64 = 34300.0

    // PSPollDelay this variable determines the amount of milliseconds
    // before the next poll can be done per sensor.
    PSPollDelay int = 200

    PSensorIDs []PSensor = []PSensor {
        PSensor {
            SensorName: "PS_Front",
            EchoPin:    21,
            TriggerPin: 20,
        },
    }
)
