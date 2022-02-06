package sonoff

import "fmt"

type SocketDevice struct {
	SonoffDevice
}

func NewSocket(id string, ip string, api int, txt int, seq int) Device {
	s := &SocketDevice{
		SonoffDevice: SonoffDevice{
			deviceType: BULB,
			deviceId:   id,
			deviceIp:   ip,
			seq:        seq,
			apiVers:    api,
			txtVers:    txt,
		},
	}
	return s
}

func (s *SocketDevice) String() string {
	s.SonoffDevice.String()
	return fmt.Sprintf(`{%s}`,
		s.SonoffDevice.String())
}
