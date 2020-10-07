package types

func NewDeviceLinks(uris ...string) DeviceLinks {
	links := []Link{}
	for _, u := range uris {
		links = append(links, Link(u))
	}
	return DeviceLinks{links}
}
