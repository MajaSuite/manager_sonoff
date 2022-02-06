package sonoff

import "fmt"

const (
	NO_TYPE Type = iota
	BULB
	DIY
	SOCKET
)

type Type byte

func (t Type) String() string {
	switch t {
	case BULB:
		return "Bulb"
	case DIY:
		return "DIY device"
	case SOCKET:
		return "Power socket"
	default:
		return "n/a"
	}
}

type Device interface {
	Type() Type
	ID() string
	IP() string
	Name() string
	String() string
}

type SonoffDevice struct {
	deviceType Type   `json:"type"`
	deviceId   string `json:"deviceid"`
	deviceIp   string `json:"ip,omitempty"`
	deviceName string `json:"name"`
	txtVers    int    `json:"txtvers"`
	apiVers    int    `json:"apivers"`
	seq        int    `json:"seq"`
}

func (sd *SonoffDevice) Type() Type {
	return sd.deviceType
}

func (sd *SonoffDevice) ID() string {
	return sd.deviceId
}

func (sd *SonoffDevice) IP() string {
	return sd.deviceIp
}

func (sd *SonoffDevice) Name() string {
	return sd.deviceName
}

func (sd *SonoffDevice) String() string {
	return fmt.Sprintf(`"type":"%s","id":"%s","ip":"%s","name":"%s",txtvers":%d,"apivers":%d,"seq":%d`,
		sd.deviceType, sd.deviceId, sd.deviceIp, sd.deviceName, sd.txtVers, sd.apiVers, sd.seq)
}
