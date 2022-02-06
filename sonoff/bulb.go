package sonoff

import "fmt"

type BulbDevice struct {
	SonoffDevice
	Encrypt bool   `json:"encrypt"`
	Data1   string `json:"data1,omitempty"`
	IV      string `json:"iv,omitempty"`
}

// [
//   data1=4Tp/SNAMhzqOdUaRY5baGfQx1MQA7Q615K5lOk2+csHwnVOelBpXMUjWb2tpCQ+HWyxnjv5yCNrwGUlQg/BOmw==
//   iv=NTc0NjU2MzIwMTg0NTc5Mg==
//   encrypt=true
//   seq=1
//   id=10010ac611
//   apivers=1
//   txtvers=1
// ]

func NewBulb(id string, ip string, api int, txt int, seq int, data1 string, iv string, enc bool) Device {
	d := &BulbDevice{
		SonoffDevice: SonoffDevice{
			deviceType: BULB,
			deviceId:   id,
			deviceIp:   ip,
			seq:        seq,
			apiVers:    api,
			txtVers:    txt,
		},
		Data1:   data1,
		IV:      iv,
		Encrypt: enc,
	}
	return d
}

func (b *BulbDevice) String() string {
	return fmt.Sprintf(`{%s, "data1":"%s","iv":"%s","encypt":%v}`,
		b.SonoffDevice.String(), b.Data1, b.IV, b.Encrypt)
}
