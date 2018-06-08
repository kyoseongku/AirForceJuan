package autodrone
// PS Stands for Proximity Sensors
import(
    "github.com/stianeikeland/go-rpio"
    "github.com/eapache/queue"
    "sort"
    //"log"
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
        pSensors[ index ].CurrReadings = queue.New()
        var echoPin    = rpio.Pin( pSensors[ index ].EchoPin    )
        var triggerPin = rpio.Pin( pSensors[ index ].TriggerPin )
        echoPin.Input()
        echoPin.PullDown()
        triggerPin.Output()
    }

    // main module loop that updates the sensors as fast as possible
    for {
        for index := 0; index < len( pSensors ); index++ {
            measureDistance( &pSensors[ index ] )
        }
    }
}

func measureDistance( sensor *PSensor ) {
    var echoPin    = rpio.Pin( (*sensor).EchoPin    )
    var triggerPin = rpio.Pin( (*sensor).TriggerPin )
    var input      = rpio.Low

    var startTime       = time.Now()
    var functionRuntime = time.Now()

    // Trigger a proximity measurement
    triggerPin.High()
    time.Sleep( time.Duration( PSTriggerPulse ) * time.Microsecond )
    triggerPin.Low()

    // Start the stopwatch, timeout in PSTimeout milliseconds
    for input == rpio.Low {
        var elapsedTime = float64( time.Since( functionRuntime ) ) / float64( time.Millisecond )
        if elapsedTime > PSTimeout { return }
        input = echoPin.Read()
    }

    startTime = time.Now()

    // Stop the stopwatch, timeout in PSTimeout milliseconds
    for input == rpio.High {
        var elapsedTime = float64( time.Since( functionRuntime ) ) / float64( time.Millisecond )
        if elapsedTime > PSTimeout { return }
        input = echoPin.Read()
    }

    // calculate and write to temporary storage
    var deltaTime = float64( time.Since( startTime ) ) / float64( time.Second )
    var psReading = PSReading {
        Distance:  deltaTime * PSSpeedOfSound / 2,
        Timestamp: startTime.String(),
    }

    // clip the value to min and max range
    if psReading.Distance < PSMinRange { psReading.Distance = PSMinRange }
    if psReading.Distance > PSMaxRange { psReading.Distance = PSMaxRange }

    // Update the sensor reading
    updateSensor( sensor, psReading )
    //log.Printf( "%s: %f", (*sensor).SensorName, GetMedianProximity( (*sensor).SensorName ) )
}

func updateSensor( sensor *PSensor, update PSReading ) {
    // Get the current readings for the sensor
    var currReadings = (*sensor).CurrReadings

    (*sensor).rw_mutex.Lock()

    // Add the next reading to the queue
    (*currReadings).Add( update )

    // if the queue size is larger the maximum read buffer, then remove the oldest
    if (*currReadings).Length() > PSMaxReadBuffer { (*currReadings).Remove() }

    (*sensor).rw_mutex.Unlock()
}

// GetMedianProximity Apply a median filter to the latest samples before obtaining a distance reading ( measured in millimeters )
func GetMedianProximity( sensorName string ) float64 {
    var proximity = -1.0
    var floatBuffer [ PSMaxReadBuffer ]float64
    // Scan through the array of proximity sensors and find the correct sensor.
    for index := 0; index < len( pSensors ); index++ {
        if sensorName == pSensors[ index ].SensorName {
            var currSensor = &pSensors[ index ]
            var currReadings = (*currSensor).CurrReadings

            (*currSensor).rw_mutex.RLock()

            var currLength = (*currReadings).Length()

            // break out if no data
            if currLength == 0 {
                (*currSensor).rw_mutex.RUnlock(); break
            }

            // fill the float buffer with the latest distance readings
            for index = 0; index < currLength; index++ {
                floatBuffer[ index ] = (*currReadings).Get( index ).(PSReading).Distance
            }

            (*currSensor).rw_mutex.RUnlock()

            // sort the buffer from least to greatest (probably nlogn)
            sort.Float64s( floatBuffer[:currLength] )

            // write the median reading into the return variable
            proximity = floatBuffer[ currLength / 2 ]
            break;
        }
    }
    return proximity
}

// GetLatestProximity get the latest sensor reading ( measured in millimeters )
func GetLatestProximity( sensorName string ) float64 {
    var proximity = -1.0
    // Scan through the array of proximity sensors and find the correct sensor.
    for index := 0; index < len( pSensors ); index++ {
        if sensorName == pSensors[ index ].SensorName {
            var currSensor = &pSensors[ index ]
            var currReadings = (*currSensor).CurrReadings

            (*currSensor).rw_mutex.RLock()

            var currLength = (*currReadings).Length()

            // break out if no data
            if currLength == 0 {
                (*currSensor).rw_mutex.RUnlock(); break
            }

            // write the latest reading into the return variable
            proximity = (*currReadings).Peek().(PSReading).Distance

            (*currSensor).rw_mutex.RUnlock()

            break;
        }
    }
    return proximity
}
