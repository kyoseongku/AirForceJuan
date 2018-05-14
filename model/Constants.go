package autodrone

var (
    // NumPropellers Number of Propellers
    NumPropellers = 4
)

// GPS Reader Constants
var (
    // BaudRate the stream flow of the serial port
    GPSBaudRate uint = 9600

    // SerialPort the port value ( directory ) that the serial stream can be read in
    GPSSerialPort = "/dev/ttyAMA0"

    // BufferSize when reading from the serial port, how much data should be pulled from the stream?
    GPSBufferSize uint = 256

    GPSDataBits uint = 8

    GPSStopBits uint = 1

    GPSMinBufferRead uint = 64
)

