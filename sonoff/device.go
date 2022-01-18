package sonoff

import "fmt"

type Device struct {
	Ip       string `json:"ip,omitempty"`
	Type     string `json:"type"`
	DeviceId string `json:"deviceid"`
	Name     string `json:"name"`
	TxtVers  int    `json:"txtvers"`
	ApiVers  int    `json:"apivers"`
	Seq      int    `json:"seq"`
	Encrypt  bool   `json:"encrypt"`
	IV       string `json:"iv,omitempty"`
	Data     Data   `json:"-"`
	Data1    string `json:"data1,omitempty"`
	Data2    string `json:"data2,omitempty"`
	Data3    string `json:"data3,omitempty"`
	Data4    string `json:"data4,omitempty"`
	Cmd      string `json:"cmd,omitempty"`
}

func (e *Device) String() string {
	return fmt.Sprintf(`{"ip":"%s","type":"%s","deviceid":"%s","txtvers":%d,"apivers":%d,"seq":%d,"encrypt":%v,"iv":"%s","data1":"%s","data2":"%s","data3":"%s","data4":"%s"}`,
		e.Ip, e.Type, e.DeviceId, e.TxtVers, e.ApiVers, e.Seq, e.Encrypt, e.IV, e.Data1, e.Data2, e.Data3, e.Data4)
}

func (d *Data) String() string {
	return fmt.Sprintf(`{"switch":"%s","startup":"%s","pulse":"%s","pulseWidth":%d,"ssid":"%s","otaUnlock":"%v","fwVersion":"%s","deviceid":"%s","bssid":"%s","signalStrength":%d"}`,
		d.Switch, d.Startup, d.Pulse, d.PulseWidth, d.Ssid, d.OtaUnlock, d.FwVersion, d.Deviceid, d.Bssid, d.Strength)
}
