package autodrone

// Define registers per MPU6050, Register Map and Descriptions, Rev 4.2, 08/19/2013 6 DOF Motion sensor fusion device
// Invensense Inc., www.invensense.com
// See also MPU-6050 Register Map and Descriptions, Revision 4.0, RM-MPU-6050A-00, 9/12/2012 for registers not listed in 
// above document; the MPU6050 and MPU 9150 are virtually identical but the latter has an on-board magnetic sensor

const (
    GA_XGOFFS_TC          byte = 0x00 // Bit 7 PWR_MODE, bits 6:1 XG_OFFS_TC, bit 0 OTP_BNK_VLD                 
    GA_YGOFFS_TC          byte = 0x01
    GA_ZGOFFS_TC          byte = 0x02
    GA_X_FINE_GAIN        byte = 0x03 // [7:0] fine gain
    GA_Y_FINE_GAIN        byte = 0x04
    GA_Z_FINE_GAIN        byte = 0x05
    GA_XA_OFFSET_H        byte = 0x06 // User-defined trim values for accelerometer
    GA_XA_OFFSET_L_TC     byte = 0x07
    GA_YA_OFFSET_H        byte = 0x08
    GA_YA_OFFSET_L_TC     byte = 0x09
    GA_ZA_OFFSET_H        byte = 0x0A
    GA_ZA_OFFSET_L_TC     byte = 0x0B
    GA_SELF_TEST_X        byte = 0x0D
    GA_SELF_TEST_Y        byte = 0x0E
    GA_SELF_TEST_Z        byte = 0x0F
    GA_SELF_TEST_A        byte = 0x10
    GA_XG_OFFS_USRH       byte = 0x13  // User-defined trim values for gyroscope; supported in MPU-6050?
    GA_XG_OFFS_USRL       byte = 0x14
    GA_YG_OFFS_USRH       byte = 0x15
    GA_YG_OFFS_USRL       byte = 0x16
    GA_ZG_OFFS_USRH       byte = 0x17
    GA_ZG_OFFS_USRL       byte = 0x18
    GA_SMPLRT_DIV         byte = 0x19
    GA_CONFIG             byte = 0x1A
    GA_GYRO_CONFIG        byte = 0x1B
    GA_ACCEL_CONFIG       byte = 0x1C
    GA_FF_THR             byte = 0x1D  // Free-fall
    GA_FF_DUR             byte = 0x1E  // Free-fall
    GA_MOT_THR            byte = 0x1F  // Motion detection threshold bits [7:0]
    GA_MOT_DUR            byte = 0x20  // Duration counter threshold for motion interrupt generation, 1 kHz rate, LSB byte = 1 ms
    GA_ZMOT_THR           byte = 0x21  // Zero-motion detection threshold bits [7:0]
    GA_ZRMOT_DUR          byte = 0x22  // Duration counter threshold for zero motion interrupt generation, 16 Hz rate, LSB byte = 64 ms
    GA_FIFO_EN            byte = 0x23
    GA_I2C_MST_CTRL       byte = 0x24
    GA_I2C_SLV0_ADDR      byte = 0x25
    GA_I2C_SLV0_REG       byte = 0x26
    GA_I2C_SLV0_CTRL      byte = 0x27
    GA_I2C_SLV1_ADDR      byte = 0x28
    GA_I2C_SLV1_REG       byte = 0x29
    GA_I2C_SLV1_CTRL      byte = 0x2A
    GA_I2C_SLV2_ADDR      byte = 0x2B
    GA_I2C_SLV2_REG       byte = 0x2C
    GA_I2C_SLV2_CTRL      byte = 0x2D
    GA_I2C_SLV3_ADDR      byte = 0x2E
    GA_I2C_SLV3_REG       byte = 0x2F
    GA_I2C_SLV3_CTRL      byte = 0x30
    GA_I2C_SLV4_ADDR      byte = 0x31
    GA_I2C_SLV4_REG       byte = 0x32
    GA_I2C_SLV4_DO        byte = 0x33
    GA_I2C_SLV4_CTRL      byte = 0x34
    GA_I2C_SLV4_DI        byte = 0x35
    GA_I2C_MST_STATUS     byte = 0x36
    GA_INT_PIN_CFG        byte = 0x37
    GA_INT_ENABLE         byte = 0x38
    GA_DMP_INT_STATUS     byte = 0x39  // Check DMP interrupt
    GA_INT_STATUS         byte = 0x3A
    GA_ACCEL_XOUT_H       byte = 0x3B
    GA_ACCEL_XOUT_L       byte = 0x3C
    GA_ACCEL_YOUT_H       byte = 0x3D
    GA_ACCEL_YOUT_L       byte = 0x3E
    GA_ACCEL_ZOUT_H       byte = 0x3F
    GA_ACCEL_ZOUT_L       byte = 0x40
    GA_TEMP_OUT_H         byte = 0x41
    GA_TEMP_OUT_L         byte = 0x42
    GA_GYRO_XOUT_H        byte = 0x43
    GA_GYRO_XOUT_L        byte = 0x44
    GA_GYRO_YOUT_H        byte = 0x45
    GA_GYRO_YOUT_L        byte = 0x46
    GA_GYRO_ZOUT_H        byte = 0x47
    GA_GYRO_ZOUT_L        byte = 0x48
    GA_EXT_SENS_DATA_00   byte = 0x49
    GA_EXT_SENS_DATA_01   byte = 0x4A
    GA_EXT_SENS_DATA_02   byte = 0x4B
    GA_EXT_SENS_DATA_03   byte = 0x4C
    GA_EXT_SENS_DATA_04   byte = 0x4D
    GA_EXT_SENS_DATA_05   byte = 0x4E
    GA_EXT_SENS_DATA_06   byte = 0x4F
    GA_EXT_SENS_DATA_07   byte = 0x50
    GA_EXT_SENS_DATA_08   byte = 0x51
    GA_EXT_SENS_DATA_09   byte = 0x52
    GA_EXT_SENS_DATA_10   byte = 0x53
    GA_EXT_SENS_DATA_11   byte = 0x54
    GA_EXT_SENS_DATA_12   byte = 0x55
    GA_EXT_SENS_DATA_13   byte = 0x56
    GA_EXT_SENS_DATA_14   byte = 0x57
    GA_EXT_SENS_DATA_15   byte = 0x58
    GA_EXT_SENS_DATA_16   byte = 0x59
    GA_EXT_SENS_DATA_17   byte = 0x5A
    GA_EXT_SENS_DATA_18   byte = 0x5B
    GA_EXT_SENS_DATA_19   byte = 0x5C
    GA_EXT_SENS_DATA_20   byte = 0x5D
    GA_EXT_SENS_DATA_21   byte = 0x5E
    GA_EXT_SENS_DATA_22   byte = 0x5F
    GA_EXT_SENS_DATA_23   byte = 0x60
    GA_MOT_DETECT_STATUS  byte = 0x61
    GA_I2C_SLV0_DO        byte = 0x63
    GA_I2C_SLV1_DO        byte = 0x64
    GA_I2C_SLV2_DO        byte = 0x65
    GA_I2C_SLV3_DO        byte = 0x66
    GA_I2C_MST_DELAY_CTRL byte = 0x67
    GA_SIGNAL_PATH_RESET  byte = 0x68
    GA_MOT_DETECT_CTRL    byte = 0x69
    GA_USER_CTRL          byte = 0x6A  // Bit 7 enable DMP, bit 3 reset DMP
    GA_PWR_MGMT_1         byte = 0x6B // Device defaults to the SLEEP mode
    GA_PWR_MGMT_2         byte = 0x6C
    GA_DMP_BANK           byte = 0x6D  // Activates a specific bank in the DMP
    GA_DMP_RW_PNT         byte = 0x6E  // Set read/write pointer to a specific start address in specified DMP bank
    GA_DMP_REG            byte = 0x6F  // Register in DMP from which to read or to which to write
    GA_DMP_REG_1          byte = 0x70
    GA_DMP_REG_2          byte = 0x71
    GA_FIFO_COUNTH        byte = 0x72
    GA_FIFO_COUNTL        byte = 0x73
    GA_FIFO_R_W           byte = 0x74
    GA_WHO_AM_I_MPU6050   byte = 0x75 // Should return byte = 0x68
)

// Initial input paameters
const(
    AFS_2G  = iota
    AFS_4G  = iota
    AFS_8G  = iota
    AFS_16G = iota
)

const(
    GFS_250DPS  = iota
    GFS_500DPS  = iota
    GFS_1000DPS = iota
    GFS_2000DPS = iota
)
