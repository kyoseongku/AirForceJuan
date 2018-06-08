package autodrone
// GA stands for gyroscope and accelerometer
import(
    "github.com/d2r2/go-i2c"
    "log"
    "time"
	"sync"
)

var(
    // gaConnection a pointer to the interface that allows reading/writing
    // to and from the module
    gaConnection *(i2c.I2C)
	gaReading GAReading
	ga_rwMutex sync.RWMutex
)

func GA_StartModule() {
    // connect to a I2C bus using the address and bus number
	tmp, err := i2c.NewI2C( GAAddress, GABusNum, false )
    if err != nil { log.Fatal( err ) }
	gaConnection = tmp
    defer gaConnection.Close()

    // Initialize the i2c slave module
	resetMPU6050()
    initializeMPU6050()

	//log.Printf( "Gyro Bias: %f, %f, %f", gyroBias[0], gyroBias[1], gyroBias[2] )
	// Main module loop.
    for {
		// Check if there's new data to read.
		if checkReadyBit() == false { continue }

		// Read in the raw accelerometer, temperature, and gyroscope data
		var accelData = readAccelData()
		var gyroData  = readGyroData()
		var tempData  = readTempData()
		// Parse the raw accelerometer and gyroscope data 
		var pAccelData = parseAccelData( accelData )
		var pGyroData  = parseGyroData ( gyroData  )
		var pTempData  = parseTempData ( tempData  )

		// Update the current value of the quaternion, and also update the temperature.
		updateGAReading( &gaReading, pAccelData, pGyroData, pTempData )
		log.Printf(
			"\nGyro: %d, %d, %d\nAccel: %d, %d, %d\nTemp: %d" +
			"\nGyro: %f, %f, %f\nAccel: %f, %f, %f\nTemp: %f",
			gyroData[0] , gyroData[1] , gyroData[2] ,
			accelData[1], accelData[1], accelData[2],
			tempData,
			pGyroData[0] , pGyroData[1] , pGyroData[2] ,
			pAccelData[1], pAccelData[1], pAccelData[2],
			pTempData )
	}
}

func resetMPU6050() {
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x80 ) // Write a one to bit 7 reset bit; toggle reset device
    time.Sleep( 100 * time.Millisecond )
}

func initializeMPU6050() {
    // Wake up the device.
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x00 ) // Clear the sleep mode bit (6) and enable all sensors.
    time.Sleep( 100 * time.Millisecond )           // Delay 100 ms for PPL to get established on x-axis gyro; should check for PPL ready interrupt.
	gaConnection.WriteRegU8( GA_GYRO_CONFIG, 0x00  )
	gaConnection.WriteRegU8( GA_ACCEL_CONFIG, 0x00 )
}

func checkReadyBit() bool {
	value, err := gaConnection.ReadRegU8( GA_INT_STATUS )
	if err != nil { log.Fatal( err ) }
	if (value & 0x01) == 0 {
		return false
	} else {
		return true
	}
}

func readAccelData() [3]int16 {
	var retBuffer = [3]int16{}
	tmp, err := gaConnection.ReadRegS16BE( GA_ACCEL_XOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[0] = tmp

	tmp, err  = gaConnection.ReadRegS16BE( GA_ACCEL_YOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[1] = tmp

	tmp, err  = gaConnection.ReadRegS16BE( GA_ACCEL_ZOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[2] = tmp
	return retBuffer
}

func readGyroData() [3]int16 {
	var retBuffer = [3]int16{}
	tmp, err := gaConnection.ReadRegS16BE( GA_GYRO_XOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[0] = tmp

	tmp, err  = gaConnection.ReadRegS16BE( GA_GYRO_YOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[1] = tmp

	tmp, err  = gaConnection.ReadRegS16BE( GA_GYRO_ZOUT_H )
	if err != nil { log.Fatal( err ) }
	retBuffer[2] = tmp
	return retBuffer
}

func readTempData() int16 {
	retValue, err := gaConnection.ReadRegS16BE( GA_TEMP_OUT_H )
	if err != nil { log.Fatal( err ) }
	return retValue
}

func parseAccelData( accelData [3]int16 ) [3]float64 {
	var aScale = getAccelScale()
	return [3]float64 {
		float64(accelData[0])*aScale,
		float64(accelData[1])*aScale,
		float64(accelData[2])*aScale,
	}
}

func parseGyroData( gyroData [3]int16 ) [3]float64 {
	var gScale = getGyroScale()
	return [3]float64 {
		float64(gyroData[0])*gScale,
		float64(gyroData[1])*gScale,
		float64(gyroData[2])*gScale,
	}
}

func parseTempData( tempData int16 ) float64 {
	var returnVar = float64( tempData ) / 340.0 + 36.53
	return returnVar
}

func getGyroScale() float64 {
	switch GScale {
		case GFS_250DPS:  { return 250.0  / 32768.0 }
		case GFS_500DPS:  { return 500.0  / 32768.0 }
		case GFS_1000DPS: { return 1000.0 / 32768.0 }
		case GFS_2000DPS: { return 2000.0 / 32768.0 }
	}
	return 0
}

func getAccelScale() float64 {
	switch AScale {
		case AFS_2G:  { return 2.0  / 32768.0 }
		case AFS_4G:  { return 4.0  / 32768.0 }
		case AFS_8G:  { return 8.0  / 32768.0 }
		case AFS_16G: { return 16.0 / 32768.0 }
	}
	return 0
}

func updateGAReading( gaReading *GAReading, g [3]float64, a [3]float64, t float64 ) {
	ga_rwMutex.Lock()

	*gaReading = GAReading {
		GyroX: g[0],
		GyroY: g[1],
		GyroZ: g[2],
		AccelX: a[0],
		AccelY: a[1],
		AccelZ: a[2],
		Temperature: t,
	}

	ga_rwMutex.Unlock()
}

func GetGAReading( get *GAReading ) {
	ga_rwMutex.RLock()

	*get = gaReading

	ga_rwMutex.RUnlock()
}
