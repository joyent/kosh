package types

func (v ValidationStateWithResults) JSON() ([]byte, error) {
	return AsJSON(v)
}

func (v ValidationStateWithResults) String() string {
	return ""
}
