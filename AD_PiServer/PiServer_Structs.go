package piserver

// PiDataType ...
type PiDataType struct {
	PropellerArray []PropellerType
	Altitude       float64
	Latitude       float64
	Longitude      float64
}

// PiControlType ...
type PiControlType struct {
	PropellerArray []PropellerType
}

// PropellerType ...
type PropellerType struct {
	Frequency float64
}
