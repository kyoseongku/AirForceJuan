package autodrone

import(
    "github.com/jacobsa/go-serial/serial"
    "strings"
    "log"
    "os"
)

// file-local variables
var (
    gpsLogger = log.New( os.Stderr, "", 0 )
)

// GPS_StartModule run through a go routine from the main function.
func GPS_StartModule() {
    // set up the configuration structure
    var config = serial.OpenOptions {
        PortName:        GPSSerialPort,
        BaudRate:        GPSBaudRate,
        DataBits:        GPSDataBits,
        StopBits:        GPSStopBits,
        MinimumReadSize: GPSMinBufferRead }

    // open the serial port for reading
    stream, err := serial.Open( config )
    if err != nil {
        log.Fatal( err )
    }
    defer stream.Close()

    // allocate buffer space
    buf := make([]byte, GPSBufferSize )

    // cumulative string storage
    var cumulativeString = ""

    // while the autodrone is active, read from pin
    for {
        numBytes, err := stream.Read( buf )
        if err != nil {
             log.Fatal( err )
        }
        var strBuf = string( buf[:numBytes] )

        // search for the start token in the current string buffer
        var tokenIndex = strings.Index( strBuf, "$" )

        // parse the string once accumulated up to the next token index 
        if tokenIndex == -1 {
            cumulativeString += strBuf
        } else {
            cumulativeString += strBuf[ 0 : tokenIndex ]
            ParseStr( cumulativeString )
            cumulativeString = strBuf[ tokenIndex : len(strBuf) ]
        }
    }
}

func ParseStr( str string ) {
    if str[0] != '$' { return }

    // TODO: Parse the string and store the value into shared memory.
    //gpsLogger.Printf( "%s", str )
}
