package autodrone
// GPS stands for global positionin system 
import(
    "github.com/jacobsa/go-serial/serial"
    "strings"
    "strconv"
    "sync"
    "log"
    "os"
)

// file-local variables
var (
    // gpsLogger logger for debugging purposes
    gpsLogger = log.New( os.Stderr, "", 0 )

    // rw_mutex ( lookup readers-writer lock ) 
    rw_mutex sync.RWMutex

    // gpsReading the current reading of the GPS
    gpsReading GPSReading
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

    // initialize gpsReading
    updateGPSReading( GPSReading {
        Altitude:     0,
        Latitude:     0,
        Longitude:    0,
        Timestamp:    "0",
        LatDirection: 'N',
        LngDirection: 'W' })

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
            parseStr( cumulativeString )
            cumulativeString = strBuf[ tokenIndex : len(strBuf) ]
        }
    }
}

func parseStr( str string ) {
    if len(str) == 0 { return }
    if len(str) <= 6 { return }
    if str[0] != '$' { return }

    //gpsLogger.Printf( "%s", str )
    var strTokens = strings.Split( str, ",")

    // update the gps state
    switch strTokens[0] {
        case "$GPGGA": updateGPGGA( strTokens )
        case "$GPGSA": updateGPGSA( strTokens )
        case "$GPGSV": updateGPGSV( strTokens )
        case "$GPGLL": updateGPGLL( strTokens )
        case "$GPRMC": updateGPRMC( strTokens )
        case "$GPVTG": updateGPVTG( strTokens )
    }
}

func updateGPGGA( strTokens []string ) {
    var currReading GPSReading
    GetGPSReading( &currReading )

    // obtain the current and token time stamp
    currTimestamp, err  := strconv.ParseFloat( currReading.Timestamp, 64 )
    if err != nil { log.Fatal( err ) }
    tokenTimestamp, err := strconv.ParseFloat( strTokens[1], 64 )
    if err != nil { log.Fatal( err ) }

    // check the timestamp in case of duplicate reading
    if( currTimestamp >= tokenTimestamp ) { return }

    // obtain latitude, longitude, altitude
    var latStr = strTokens[2]
    var lngStr = strTokens[4]
    altitudeReading, err := strconv.ParseFloat( strTokens[9], 64 )
    if err != nil { log.Fatal( err ) }
    latDegReading, err := strconv.ParseFloat( latStr[0:2],64 )
    if err != nil { log.Fatal( err ) }
    latMinReading, err := strconv.ParseFloat( latStr[2:len(latStr)], 64 )
    if err != nil { log.Fatal( err ) }
    lngDegReading, err := strconv.ParseFloat( lngStr[0:3],64 )
    if err != nil { log.Fatal( err ) }
    lngMinReading, err := strconv.ParseFloat( lngStr[3:len(lngStr)], 64 )
    if err != nil { log.Fatal( err ) }

    // place into temporary storage
    currReading = GPSReading {
        Altitude:     altitudeReading,
        Latitude:     latDegReading + latMinReading / 60.0,
        Longitude:    lngDegReading + lngMinReading / 60.0,
        Timestamp:    strTokens[1],
        LatDirection: strTokens[3][0],
        LngDirection: strTokens[5][0] }

    // update the gps reading
    updateGPSReading( currReading )
}

func updateGPGSA( strTokens []string ) {
    // POSSIBLY TO ANALYZE SATELITES USED
}

func updateGPGSV( strTokens []string ) {
    // POSSIBLY TO ANALYZE SATELITES USED
}

func updateGPGLL( strTokens []string ) {
    // ONLY OBTAIN LATITUDE AND LONGITUDE, BUT NO ALTITUDE
}

func updateGPRMC( strTokens []string ) {
    // POSSIBLY TO ANALYZE SPEED/VELOCITY DATA
}

func updateGPVTG( strTokens []string ) {
    // POSSIBLY TO ANALYZE SPEED/VELOCITY DATA
}

// UpdateGPSReading updates the local file variable gpsReading while keeping sync
func updateGPSReading( update GPSReading ) {
    rw_mutex.Lock()
    gpsReading = update
    rw_mutex.Unlock()
}

// GetGPSReading copies the GPS reading from the gpsReading variable
func GetGPSReading( get *GPSReading ) {
    rw_mutex.RLock()
    *get = gpsReading
    rw_mutex.RUnlock()
}
