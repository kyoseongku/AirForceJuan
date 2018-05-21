package autodrone
// PS Stands for Proximity Sensors
import(
    "github.com/stianeikeland/go-rpio"
    // "log"
    "time"
)

var(
    pSensors []PSensor
)

// PS_StartModule ...
func PS_StartModule() {
    // initialize the GPIO pins - set the IO mode for each pin and
    // set the echo pin as a pulldown circuit ( to prefer to hold the echo signal at 0V )
    // ( look up pulldown resistor circuit )
    pSensors = make( []PSensor, len( PSensorIDs ) )
    for index := 0; index < len( PSensorIDs ); index++ {
        pSensors[ index ] = PSensorIDs[ index ]
        var currSensor = pSensors[ index ]
        var echoPin    = rpio.Pin( currSensor.EchoPin    )
        var triggerPin = rpio.Pin( currSensor.TriggerPin )
        echoPin.Input()
        triggerPin.Output()
        triggerPin.PullDown()
    }

    // main module loop that updates the sensors as fast as possible
    for {
        for index := 0; index < len( pSensors ); index++ {
            measureDistance( &pSensors[ index ] )
        }
    }
}

// TODO: Optimize the distance measurement by playing with poll delay 
// and the better management of which sensor is off its cooldown period
// possibly add a safeguard for the locking case where the rpi is not able
// to read the echo at all.
func measureDistance( sensor *PSensor ) {
    var echoPin    = rpio.Pin( (*sensor).EchoPin    )
    var triggerPin = rpio.Pin( (*sensor).TriggerPin )
    var input      = rpio.Low
    var startTime  = time.Now()

    // Trigger a proximity measurement
    triggerPin.High()
    time.Sleep( 10 * time.Microsecond )
    triggerPin.Low()

    // Start the stopwatch // CAUTION: CAN LOCK THE MODULE
    for input == rpio.Low {
        input = echoPin.Read()
    }

    startTime = time.Now()

    // Stop the stopwatch // CAUTION: CAN LOCK THE MODULE
    for input == rpio.High {
        input = echoPin.Read()
    }

    // calculate and write to temporary storage
    var deltaTime = float64( time.Since( startTime ) ) / float64( time.Second )
    var psReading = PSReading {
        Distance:  deltaTime * PSSpeedOfSound / 2,
        Timestamp: time.Since( startTime ).String(),
    }

    // Each sensor's measuring rate should be at minimum per 60 ms or greater
    // I set this sleep to 200 for extreme caution due to locking of the
    // module from a bad echo read.
    time.Sleep( time.Duration( PSPollDelay ) * time.Millisecond )

    // Update the sensor reading
    // log.Printf( "Reading: %f",  )
    updateSensor( sensor, psReading )
}

func updateSensor( sensor *PSensor, update PSReading ) {
    (*sensor).rw_mutex.Lock()
    (*sensor).CurrReading = update
    (*sensor).rw_mutex.Unlock()
}

func GetProximityReading( sensorName string, get *PSReading ) {
    // Scan through the array of proximity sensors and find the correct one.
    for index := 0; index < len( pSensors ); index++ {
        var currSensor = &pSensors[ index ]
        // get the sensor reading.
        if sensorName == currSensor.SensorName {
            (*currSensor).rw_mutex.RLock()
            *get = currSensor.CurrReading
            (*currSensor).rw_mutex.RUnlock()
            return
        }
    }
}

