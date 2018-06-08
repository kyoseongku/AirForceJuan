package autodrone
// GA stands for gyroscope and accelerometer
import(
    "github.com/d2r2/go-i2c"
    "log"
    "time"
	"math"
	"sync"
)

var(
    // gaConnection a pointer to the interface that allows reading/writing
    // to and from the module
    gaConnection *(i2c.I2C)

	gyroBias  = [3]float64 {0,0,0}
	accelBias = [3]float64 {0,0,0}
	currQuaternion = [4]float64 { 1,0,0,0 }
	currTemperature float64 = 0
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
	calibrateMPU6050()
    initializeMPU6050()

	// Test the MPU 6050
	var testResult = selfTestMPU6050()
	var testPassed = true
	for i := 0; i < len(testResult); i++ {
		if testResult[i] > 1.0 {
			testPassed = false
			break
		}
	}
	if !testPassed {
		log.Fatal( "Device did not pass self test. " )
		return
	}

	log.Printf( "Gyro Bias: %f, %f, %f", gyroBias[0], gyroBias[1], gyroBias[2] )
	// Main module loop.
	var prevTime  = time.Now()
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
		// Calculate the time passed since last filter check.
		var deltaTime = float64( time.Since( prevTime ) ) / float64( time.Second )

		// Save the current time.
		prevTime = time.Now()

		// Update the current value of the quaternion, and also update the temperature.
		ga_rwMutex.Lock()

		updateQuaternion( &currQuaternion, pAccelData, pGyroData, deltaTime )
		currTemperature = pTempData
		//log.Printf( "%f, %f, %f, %f", currQuaternion[0], currQuaternion[1], currQuaternion[2], currQuaternion[3] )

		ga_rwMutex.Unlock()

    }
}

func initializeMPU6050() {
    // Wake up the device.
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x00 ) // Clear the sleep mode bit (6) and enable all sensors.
    time.Sleep( 100 * time.Millisecond )           // Delay 100 ms for PPL to get established on x-axis gyro; should check for PPL ready interrupt.

    // Get a stable time source.
    (*gaConnection).WriteRegU8( GA_PWR_MGMT_1, 0x01 ) // Set the clock source to be PPL with x-axis gyro reference, bits 2:0 = 001

    // Configure Gyro and Accelerometer
    // Disable FSYNC and set accelerometer and gyro bandwidth to 44 and 42 Hz, respectively; 
    // DLPF_CFG = bits 2:0 = 010; this sets the sample rate at 1 kHz for both
    // Maximum delay is 4.9 ms which is just over a 200 Hz maximum rate
    gaConnection.WriteRegU8( GA_CONFIG, 0x03 )

    // Set sample rate = gyroscope output rate/(1 + SMPLRT_DIV)
    gaConnection.WriteRegU8( GA_SMPLRT_DIV, 0x04 ) // Use a 200 Hz rate; the same rate set in CONFIG above

    // Range selects FS_SEL and AFS_SEL are 0 - 3, so 2-bit values are left-shifted into positions 4:3
    // Set gyroscope full scale range
    c, err := gaConnection.ReadRegU8( GA_GYRO_CONFIG )
    if err != nil { log.Fatal( err ) }
    gaConnection.WriteRegU8( GA_GYRO_CONFIG, c & ^byte(0xE0)   ) // Clear self-test bits [7:5] 
    gaConnection.WriteRegU8( GA_GYRO_CONFIG, c & ^byte(0x18)   ) // Clear AFS bits [4:3]
    gaConnection.WriteRegU8( GA_GYRO_CONFIG, c |  byte(GScale) << 3 ) // Set full scale range for the gyro

    // Set accelerometer configuration
    c, err  = gaConnection.ReadRegU8( GA_ACCEL_CONFIG )
    if err != nil { log.Fatal( err ) }
    gaConnection.WriteRegU8( GA_ACCEL_CONFIG, c & ^byte(0xE0)   ) // Clear self-test bits [7:5] 
    gaConnection.WriteRegU8( GA_ACCEL_CONFIG, c & ^byte(0x18)   ) // Clear AFS bits [4:3]
    gaConnection.WriteRegU8( GA_ACCEL_CONFIG, c | byte(AScale) << 3 ) // Set full scale range for the accelerometer 

    // Configure Interrupts and Bypass Enable
    // Set interrupt pin active high, push-pull, and clear on read of INT_STATUS, enable I2C_BYPASS_EN so additional chips 
    // can join the I2C bus and all can be controlled by the Rpi as master
    gaConnection.WriteRegU8( GA_INT_PIN_CFG, 0x22 )
    gaConnection.WriteRegU8( GA_INT_ENABLE , 0x01 )  // Enable data ready (bit 0) interrupt
}

// Function which accumulates gyro and accelerometer data after device initialization. It calculates the average
// of the at-rest readings and then loads the resulting offsets into accelerometer and gyro bias registers.
func calibrateMPU6050() {
    // reset device, reset all registers, clear gyro and accelerometer bias registers
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x80 ) // Write a one to bit 7 reset bit; toggle reset device
    time.Sleep( 100 * time.Millisecond )

    // get stable time source
    // Set clock source to be PLL with x-axis gyroscope reference, bits 2:0 = 001
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x01 )
    gaConnection.WriteRegU8( GA_PWR_MGMT_2, 0x00 )
    time.Sleep( 200 * time.Millisecond )

    // Configure device for bias calculation
    gaConnection.WriteRegU8( GA_INT_ENABLE  , 0x00 ) // Disable all interrupts
    gaConnection.WriteRegU8( GA_FIFO_EN     , 0x00 ) // Disable FIFO
    gaConnection.WriteRegU8( GA_PWR_MGMT_1  , 0x00 ) // Turn on internal clock source
    gaConnection.WriteRegU8( GA_I2C_MST_CTRL, 0x00 ) // Disable I2C master
    gaConnection.WriteRegU8( GA_USER_CTRL   , 0x00 ) // Disable FIFO and I2C master modes
    gaConnection.WriteRegU8( GA_USER_CTRL   , 0x0C ) // Reset FIFO and DMP
    time.Sleep( 150 * time.Millisecond );

    // Configure MPU6050 gyro and accelerometer for bias calculation
    gaConnection.WriteRegU8( GA_CONFIG      , 0x01 ) // Set low-pass filter to 188 Hz
    gaConnection.WriteRegU8( GA_SMPLRT_DIV  , 0x00 ) // Set sample rate to 1 kHz
    gaConnection.WriteRegU8( GA_GYRO_CONFIG , 0x00 ) // Set gyro full-scale to 250 degrees per second, maximum sensitivity
    gaConnection.WriteRegU8( GA_ACCEL_CONFIG, 0x00 ) // Set accelerometer full-scale to 2 g, maximum sensitivity

    // Configure FIFO to capture accelerometer and gyro data for bias calculation
    gaConnection.WriteRegU8( GA_USER_CTRL, 0x40 ) // Enable FIFO  
    gaConnection.WriteRegU8( GA_FIFO_EN  , 0x78 ) // Enable gyro and accelerometer sensors for FIFO  (max size 1024 bytes in MPU-6050)
    time.Sleep( 80 * time.Millisecond )           // accumulate 80 samples in 80 milliseconds = 960 bytes

    // At end of sample accumulation, turn off FIFO sensor read
    gaConnection.WriteRegU8( GA_FIFO_EN, 0x00 )                    // Disable gyro and accelerometer sensors for FIFO
    fifo_count, err := gaConnection.ReadRegU16BE( GA_FIFO_COUNTH ) // read FIFO sample count
    if err != nil { log.Fatal( err ) }
    var packet_count = fifo_count/12; // How many sets of full gyro and accelerometer data for averaging
    var tmpGyroBias  = [3]int32 {0,0,0}
    var tmpAccelBias = [3]int32 {0,0,0}

    for ii := uint16(0); ii < packet_count; ii++ {
        var accelBuffer = [3]int16 {0, 0, 0}
        var gyroBuffer  = [3]int16 {0, 0, 0}
        data, b, err := gaConnection.ReadRegBytes( GA_FIFO_R_W, 12 ) // read data for averaging
        if err != nil { log.Printf( "%i", b ); log.Fatal( err ) }
        // Form signed 16-bit integer for each sample in FIFO
        accelBuffer[0] = int16(data[ 0] << 8) | int16(data[ 1])
        accelBuffer[1] = int16(data[ 2] << 8) | int16(data[ 3])
        accelBuffer[2] = int16(data[ 4] << 8) | int16(data[ 5])
        gyroBuffer[0]  = int16(data[ 6] << 8) | int16(data[ 7])
        gyroBuffer[1]  = int16(data[ 8] << 8) | int16(data[ 9])
        gyroBuffer[2]  = int16(data[10] << 8) | int16(data[11])

        // Sum individual signed 16-bit biases to get accumulated signed 32-bit biases
        tmpAccelBias[0] += int32( accelBuffer[0] )
        tmpAccelBias[1] += int32( accelBuffer[1] )
        tmpAccelBias[2] += int32( accelBuffer[2] )
        tmpGyroBias[0]  += int32( gyroBuffer[0]  )
        tmpGyroBias[1]  += int32( gyroBuffer[1]  )
        tmpGyroBias[2]  += int32( gyroBuffer[2]  )
    }

    // Normalize sums to get average count biases
    tmpAccelBias[0] /= int32( packet_count )
    tmpAccelBias[1] /= int32( packet_count )
    tmpAccelBias[2] /= int32( packet_count )
    tmpGyroBias[0]  /= int32( packet_count )
    tmpGyroBias[1]  /= int32( packet_count )
    tmpGyroBias[2]  /= int32( packet_count )

    // Remove gravity from the z-axis accelerometer bias calculation
    if tmpAccelBias[2] > 0 { tmpAccelBias[2] -= int32( AccelSensitivity )
    } else                 { tmpAccelBias[2] += int32( AccelSensitivity )
    }

    // Construct the gyro biases for push to the hardware gyro bias registers, which are reset to zero upon device startup
    var bias = [6]uint8 {0,0,0,0,0,0}
    bias[0] = uint8( ( -tmpGyroBias[0]/4 >> 8 ) & 0xFF ) // Divide by 4 to get 32.9 LSB per deg/s to conform to expected bias input format
    bias[1] = uint8( ( -tmpGyroBias[0]/4      ) & 0xFF ) // Biases are additive, so change sign on calculated average gyro biases
    bias[2] = uint8( ( -tmpGyroBias[1]/4 >> 8 ) & 0xFF )
    bias[3] = uint8( ( -tmpGyroBias[1]/4      ) & 0xFF )
    bias[4] = uint8( ( -tmpGyroBias[2]/4 >> 8 ) & 0xFF )
    bias[5] = uint8( ( -tmpGyroBias[2]/4      ) & 0xFF )

    // Push gyro biases to hardware registers; works well for gyro but not for accelerometer
    // gaConnection.WriteRegU8( GA_XG_OFFS_USRH, bias[0] )
    // gaConnection.WriteRegU8( GA_XG_OFFS_USRL, bias[1] )
    // gaConnection.WriteRegU8( GA_YG_OFFS_USRH, bias[2] )
    // gaConnection.WriteRegU8( GA_YG_OFFS_USRL, bias[3] )
    // gaConnection.WriteRegU8( GA_ZG_OFFS_USRH, bias[4] )
    // gaConnection.WriteRegU8( GA_ZG_OFFS_USRL, bias[5] )

    // construct gyro bias in deg/s for later manual subtraction
    gyroBias[0] = float64( tmpGyroBias[0] ) / float64( GyroSensitivity )
    gyroBias[1] = float64( tmpGyroBias[1] ) / float64( GyroSensitivity )
    gyroBias[2] = float64( tmpGyroBias[2] ) / float64( GyroSensitivity )

    // Construct the accelerometer biases for push to the hardware accelerometer bias registers. These registers contain
    // factory trim values which must be added to the calculated accelerometer biases; on boot up these registers will hold
    // non-zero values. In addition, bit 0 of the lower byte must be preserved since it is used for temperature
    // compensation calculations. Accelerometer bias registers expect bias input as 2048 LSB per g, so that
    // the accelerometer biases calculated above must be divided by 8.

    var factoryAccelBias = [3]int32{0, 0, 0}; // A place to hold the factory accelerometer trim biases

    offset, err := gaConnection.ReadRegS16BE( GA_XA_OFFSET_H ) // Read factory accelerometer trim values
    if err != nil { log.Fatal( err ) }
    factoryAccelBias[0] = int32( offset )

    offset, err  = gaConnection.ReadRegS16BE( GA_YA_OFFSET_H ) // Read factory accelerometer trim values
    if err != nil { log.Fatal( err ) }
    factoryAccelBias[1] = int32( offset )

    offset, err  = gaConnection.ReadRegS16BE( GA_ZA_OFFSET_H ) // Read factory accelerometer trim values
    if err != nil { log.Fatal( err ) }
    factoryAccelBias[2] = int32( offset )

	var mask int32 = 1;              // Define mask for temperature compensation bit 0 of lower byte of accelerometer bias registers
    var mask_bit = [3]uint8{0, 0, 0}; // Define array to hold mask bit for each accelerometer bias axis

    for ii := 0; ii < 3; ii++ {
        if factoryAccelBias[ii] & mask != 0 { mask_bit[ii] = 0x01 } // If temperature compensation bit is set, record that fact in mask_bit
    }

    // Construct total accelerometer bias, including calculated average accelerometer bias from above
    factoryAccelBias[0] -= ( tmpAccelBias[0]/8 ) // Subtract calculated averaged accelerometer bias scaled to 2048 LSB/g (16 g full scale)
    factoryAccelBias[1] -= ( tmpAccelBias[1]/8 )
    factoryAccelBias[2] -= ( tmpAccelBias[2]/8 )


	var data = [6]byte {0,0,0,0,0,0}
    data[0] = uint8( ( factoryAccelBias[0] >> 8) & 0xFF )
    data[1] = uint8( ( factoryAccelBias[0]     ) & 0xFF )
    data[1] = uint8( data[1] | mask_bit[0] ); // preserve temperature compensation bit when writing back to accelerometer bias registers

    data[2] = uint8( ( factoryAccelBias[1] >> 8) & 0xFF )
    data[3] = uint8( ( factoryAccelBias[1]     ) & 0xFF )
    data[3] = uint8( data[3] | mask_bit[1] ); // preserve temperature compensation bit when writing back to accelerometer bias registers

    data[4] = uint8( ( factoryAccelBias[2] >> 8) & 0xFF )
    data[5] = uint8( ( factoryAccelBias[2]     ) & 0xFF )
    data[5] = uint8( data[5] | mask_bit[2] ); // preserve temperature compensation bit when writing back to accelerometer bias registers

    // Push accelerometer biases to hardware registers; doesnt work well for accelerometer
    // gaConnection.WriteRegU8( GA_XA_OFFSET_H   , data[0] )
    // gaConnection.WriteRegU8( GA_XA_OFFSET_L_TC, data[1] )
    // gaConnection.WriteRegU8( GA_YA_OFFSET_H   , data[2] )
    // gaConnection.WriteRegU8( GA_YA_OFFSET_L_TC, data[3] )
    // gaConnection.WriteRegU8( GA_ZA_OFFSET_H   , data[4] )
    // gaConnection.WriteRegU8( GA_ZA_OFFSET_L_TC, data[5] )


    // Output scaled accelerometer biases for manual subtraction in the main program
    accelBias[0] = float64(tmpAccelBias[0]) / float64( AccelSensitivity )
    accelBias[1] = float64(tmpAccelBias[1]) / float64( AccelSensitivity )
    accelBias[2] = float64(tmpAccelBias[2]) / float64( AccelSensitivity )
}


func resetMPU6050() {
    gaConnection.WriteRegU8( GA_PWR_MGMT_1, 0x80 ) // Write a one to bit 7 reset bit; toggle reset device
    time.Sleep( 100 * time.Millisecond )
}


func selfTestMPU6050() [6]float64{
	var returnValue = [6]float64{ 0,0,0,0,0,0 }
	var rawData     = [4]uint8{ 0,0,0,0 }
	var selfTest    = [6]uint8{ 0,0,0,0 }
	var factoryTrim = [6]float64{ 0,0,0,0,0,0 }

	gaConnection.WriteRegU8( GA_ACCEL_CONFIG, 0xF0 ) // Enable self test on all three axes and set accelerometer range to +/- 8 g
	gaConnection.WriteRegU8( GA_GYRO_CONFIG, 0xE0  ) // Enable self test on all three axes and set gyo range to +/- 250 degree/s
	time.Sleep( 250 * time.Millisecond ) // Delay a while to let the device execute a self test.

	// Get the X/Y/Z/Mixed axes Test results
	data, err := gaConnection.ReadRegU8( GA_SELF_TEST_X );
	if err != nil { log.Fatal( err ) }
	rawData[0] = data

	data, err  = gaConnection.ReadRegU8( GA_SELF_TEST_Y );
	if err != nil { log.Fatal( err ) }
	rawData[1] = data

	data, err  = gaConnection.ReadRegU8( GA_SELF_TEST_Z );
	if err != nil { log.Fatal( err ) }
	rawData[2] = data

	data, err  = gaConnection.ReadRegU8( GA_SELF_TEST_A );
	if err != nil { log.Fatal( err ) }
	rawData[3] = data

	// Extract the acceleration test results first 
	selfTest[0] = ( (rawData[0] >> 3) | (rawData[3] & 0x30) >> 4 ) // XA_TEST result is a five bit unsigned integer 
	selfTest[1] = ( (rawData[1] >> 3) | (rawData[3] & 0x0C) >> 2 ) // YA_TEST result is a five bit unsigned integer 
	selfTest[2] = ( (rawData[2] >> 3) | (rawData[3] & 0x03) >> 0 ) // ZA_TEST result is a five bit unsigned integer 

	// Extract the gyration test results
	selfTest[3] = rawData[0] & 0x1F // XG_TEST result is a five bit unsigned integer
	selfTest[4] = rawData[1] & 0x1F // YG_TEST result is a five bit unsigned integer
	selfTest[5] = rawData[2] & 0x1F // ZG_TEST result is a five bit unsigned integer

	// Process results to allow final comparison with factory set values
	factoryTrim[0] = (4096.0*0.34)*(math.Pow((0.92/0.34),((float64(selfTest[0]) - 1.0)/30.0))); // FT[Xa] factory trim calculation
    factoryTrim[1] = (4096.0*0.34)*(math.Pow((0.92/0.34),((float64(selfTest[1]) - 1.0)/30.0))); // FT[Ya] factory trim calculation
	factoryTrim[2] = (4096.0*0.34)*(math.Pow((0.92/0.34),((float64(selfTest[2]) - 1.0)/30.0))); // FT[Za] factory trim calculation
	factoryTrim[3] = ( 25.0*131.0)*(math.Pow(1.046,(float64(selfTest[3]) - 1.0)));              // FT[Xg] factory trim calculation
	factoryTrim[4] = (-25.0*131.0)*(math.Pow(1.046,(float64(selfTest[4]) - 1.0)));              // FT[Yg] factory trim calculation
	factoryTrim[5] = ( 25.0*131.0)*(math.Pow(1.046,(float64(selfTest[5]) - 1.0)));              // FT[Zg] factory trim calculation

	// Report results as a ratio of (STR - FT)/FT; the change from Factory Trim of the Self-Test Response
	// To get to percent, must multiply by 100 and subtract result from 100
    for i := 0; i < len(returnValue); i++ {
	      returnValue[i] = 100.0 + 100.0*(float64(selfTest[i]) - factoryTrim[i])/factoryTrim[i]; // Report percent differences
	}
	return returnValue
}

func checkReadyBit() bool {
	value, err := gaConnection.ReadRegU8( GA_INT_STATUS )
	if err != nil { log.Fatal( err ) }
	if (value & 0x01) == 0 { return false
	} else { return true }
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

func parseAccelData( accelData [3]int16 ) [3]float64 {
	var aScale = getAccelScale()
	return [3]float64 {
		float64(accelData[0])*aScale - accelBias[0],
		float64(accelData[1])*aScale - accelBias[1],
		float64(accelData[2])*aScale - accelBias[2],
	}
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
	//log.Printf( "GYRO : %i, %i, %i ", retBuffer[0], retBuffer[1], retBuffer[2] )
	return retBuffer
}

func parseGyroData( gyroData [3]int16 ) [3]float64 {
	var gScale = getGyroScale()
	var retVal = [3]float64 {
		float64(gyroData[0])*gScale - gyroBias[0],
		float64(gyroData[1])*gScale - gyroBias[1],
		float64(gyroData[2])*gScale - gyroBias[2],
	}
	retVal[0] *= math.Pi/180.0
	retVal[1] *= math.Pi/180.0
	retVal[2] *= math.Pi/180.0
	return retVal
}

func readTempData() int16 {
	retValue, err := gaConnection.ReadRegS16BE( GA_TEMP_OUT_H )
	if err != nil { log.Fatal( err ) }
	return retValue
}

// Temperature of the GA is in degrees centigrade.
func parseTempData( tempData int16 ) float64 {
	return float64( tempData ) / 340.0  + 36.53
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

// Implementation of Sebastian Madgwick's "...efficient orientation filter for... inertial/magnetic sensor arrays"
// (see http://www.x-io.co.uk/category/open-source/ for examples and more details)
// which fuses acceleration and rotation rate to produce a quaternion-based estimate of relative
// device orientation -- which can be converted to yaw, pitch, and roll. Useful for stabilizing quadcopters, etc.
// The performance of the orientation filter is at least as good as conventional Kalman-based filtering algorithms
// but is much less computationally intensive---it can be performed on a 3.3 V Pro Mini operating at 8 MHz!
func updateQuaternion( q *[4]float64, g [3]float64, a [3]float64, deltaTime float64 ) {
	// Quick variable rename
	var q1, q2, q3, q4 = (*q)[0], (*q)[1], (*q)[2], (*q)[3]
	var ax, ay, az     = a[0], a[1], a[2]
	var gx, gy, gz     = g[0], g[1], g[2]
	log.Printf( "\nGyro: %f, %f, %f\nAccel: %f, %f, %f\nDelta: %f", gx, gy, gz, ax, ay, az, deltaTime )
	// vector Normal
	var norm float64 = 0
	// Objective function elements
	var f1, f2, f3 float64 = 0, 0, 0
	// Objective function jacobian elements
	var J_11_24, J_12_23, J_13_22, J_14_21, J_32, J_33 float64 = 0,0,0,0,0,0
	// quaternion dot products
	var qd1, qd2, qd3, qd4 float64 = 0,0,0,0
	// hat dot products
	var hd1, hd2, hd3, hd4 float64 = 0,0,0,0
	// Gyroscope bias and error
	var gyroBiasX , gyroBiasY , gyroBiasZ float64 = 0,0,0
	var gyroErrorX, gyroErrorY, gyroErrorZ float64 = 0,0,0
	// Auxiliary variables to avoid repeated arithmetic
	var _hq1, _hq2, _hq3, _hq4 float64 = 0.5*q1, 0.5*q2, 0.5*q3, 0.5*q4
	var _2q1, _2q2, _2q3, _2q4 float64 = 2.0*q1, 2.0*q2, 2.0*q3, 2.0*q4
	// Normalize the accelerometer measurement
	norm = math.Sqrt( ax*ax + ay*ay + az*az )
	if norm == 0.0 { return } // handle NaN
	norm = 1.0 / norm
	ax, ay, az = ax*norm, ay*norm, az*norm
	// Compute the objective functions
	f1 =       _2q2*q4 - _2q1*q3 - ax
	f2 =       _2q1*q2 + _2q3*q4 - ay
	f3 = 1.0 - _2q2*q2 - _2q3*q3 - az
	// Compute the jacobian
	J_11_24 = _2q3
	J_12_23 = _2q4
	J_13_22 = _2q1
	J_14_21 = _2q2
	J_32    = 2.0 * J_14_21
	J_33    = 2.0 * J_11_24
	// Compute the gradient ( Matrix Multiplication )
	hd1 = J_14_21*f2 - J_11_24*f1
	hd4 = J_14_21*f1 + J_11_24*f2
	hd2 = J_12_23*f1 + J_13_22*f2 - J_32*f3
	hd3 = J_12_23*f2 - J_13_22*f1 - J_33*f3
	// Normalize the gradient
	norm = math.Sqrt( hd1*hd1 + hd2*hd2 + hd3*hd3 + hd4*hd4 )
	norm = 1.0/norm
	hd1 *= norm
	hd2 *= norm
	hd3 *= norm
	hd4 *= norm
	// Compute estimated gyroscope bias
	gyroErrorX = _2q1*hd2 - _2q2*hd1 - _2q3*hd4 + _2q4*hd3
	gyroErrorY = _2q1*hd3 + _2q2*hd4 - _2q3*hd1 - _2q4*hd2
	gyroErrorZ = _2q1*hd4 - _2q2*hd3 + _2q3*hd2 - _2q4*hd1
	// Compute and remove gyroscope bias
	gyroBiasX += gyroErrorX * deltaTime * GAzeta
	gyroBiasY += gyroErrorY * deltaTime * GAzeta
	gyroBiasZ += gyroErrorZ * deltaTime * GAzeta
	gx -= gyroBiasX
	gy -= gyroBiasY
	gz -= gyroBiasZ
	// Compute the quaternion derivative
	qd1 = _hq2*gx - _hq3*gy - _hq4*gz
	qd2 = _hq1*gx + _hq3*gz - _hq4*gy
	qd3 = _hq1*gy - _hq2*gz + _hq4*gx
	qd4 = _hq1*gz + _hq2*gy - _hq3*gx
	// Compute, and then integrate estimated quaternion derivative
	q1 += (qd1 - (GAbeta*hd1))*deltaTime
	q2 += (qd2 - (GAbeta*hd2))*deltaTime
	q3 += (qd3 - (GAbeta*hd3))*deltaTime
	q4 += (qd4 - (GAbeta*hd4))*deltaTime
	// Normalize the Quaternion
	norm = math.Sqrt( q1*q1 + q2*q2 + q3*q3 + q4*q4 )
	if norm == 0.0 { return } // handle NaN
	norm = 1.0/norm
	(*q)[0] = q1*norm
	(*q)[1] = q2*norm
	(*q)[2] = q3*norm
	(*q)[3] = q4*norm
}

func GetQuaternion() [4]float64 {
	ga_rwMutex.RLock()

	var retVal = [4]float64{
		currQuaternion[0],
		currQuaternion[1],
		currQuaternion[2],
		currQuaternion[3],
	}

	ga_rwMutex.RUnlock()
	return retVal
}

/*
func GetEuler() [4]float64 {
	ga_rwMutex.RLock()
	
	
	ga_rwMutex.RUnlock()
}
*/

func GetTemperature() float64 {
	return currTemperature
}
