package types

func (d DeviceSetting) JSON() ([]byte, error) { return AsJSON(d) }
func (d DeviceSetting) String() string        { panic("TODO") }

func (d DeviceSettings) JSON() ([]byte, error) { return AsJSON(d) }
func (d DeviceSettings) String() string        { panic("TODO") }

func (d DeviceNic) JSON() ([]byte, error) { return AsJSON(d) }
func (d DeviceNic) String() string        { panic("TODO") }

func (d DeviceLocation) JSON() ([]byte, error) { return AsJSON(d) }
func (d DeviceLocation) String() string        { panic("TODO") }
