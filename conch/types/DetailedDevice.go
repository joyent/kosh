package types

func (d DetailedDevice) JSON() ([]byte, error) {
	return AsJSON(d)
}

func (d DetailedDevice) String() string {
	return "[ comming soon ]"
}

func (d DeviceReport) JSON() ([]byte, error) {
	return AsJSON(d)
}

func (d DeviceReport) String() string {
	return "[ comming soon ]"
}
