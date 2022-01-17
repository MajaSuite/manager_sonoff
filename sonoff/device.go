package sonoff

import "fmt"

/*
{
	data1=WQimgeLh4QyLfmg4wCqjmxLKJPIGAn2xqwECJyIwYPVCKDw2HcGpjC3+FStTkx0wiR/r68Qkvutyz55LnLC/sQ==
	iv=MTkzNTAxNTc3MTYzODk3OQ==
	encrypt=true
	seq=1
	id=10010aaa26
	apivers=1
	type=light
	txtvers=1
}

{
	"txtvers=1",
	"id"="10010aaa26",
	"type"="switch?",
	"apivers"=1,
	"seq"=1,
	"data1"={"switch":"on","startup":"stay","pulse":"on","pulseWidth":2000,"ssid":"eWeLink","otaUnlock":true}
}

*/
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
