package sonoff

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type DiyDevice struct {
	SonoffDevice
	Data
	Cmd   string `json:"cmd,omitempty"`
	Data1 string `json:"data1,omitempty"`
	Data2 string `json:"data2,omitempty"`
}

func (d *DiyDevice) String() string {
	return fmt.Sprintf(`{%s, %s}`,
		d.SonoffDevice.String(), d.Data.String())
}

type Data struct {
	Switch     string `json:"switch,omitempty"`
	Startup    string `json:"startup,omitempty"`
	Stay       string `json:"stay,omitempty"`
	Pulse      string `json:"pulse,omitempty"`
	PulseWidth int    `json:"pulseWidth,omitempty"`
	Ssid       string `json:"ssid,omitempty"`
	Password   string `json:"password,omitempty"`
	OtaUnlock  bool   `json:"otaUnlock,omitempty"`
	FwVersion  string `json:"fwVersion,omitempty"`
	Deviceid   string `json:"deviceid,omitempty"`
	Bssid      string `json:"bssid,omitempty"`
	Strength   int    `json:"signalStrength,omitempty"`
	Url        string `json:"downloadUrl,omitempty"`
	Sha256sum  string `json:"sha256sum,omitempty"`
}

func (d *Data) String() string {
	return fmt.Sprintf(`"switch":"%s","startup":"%s","pulse":"%s","pulseWidth":%d,"ssid":"%s","otaUnlock":"%v","fwVersion":"%s","deviceid":"%s","bssid":"%s","signalStrength":%d"`,
		d.Switch, d.Startup, d.Pulse, d.PulseWidth, d.Ssid, d.OtaUnlock, d.FwVersion, d.Deviceid, d.Bssid, d.Strength)
}

type Request struct {
	Deviceid string `json:"deviceid"`
	Data     Data   `json:"data"`
}

type Response struct {
	Seq   int  `json:"seq"`
	Error int  `json:"error"`
	Data  Data `json:"data"`
}

func NewDIY(id string, ip string, api int, txt int, seq int, data1 string) Device {
	var data Data
	if err := json.Unmarshal([]byte(data1), &data); err != nil {
		return nil
	}

	d := &DiyDevice{
		SonoffDevice: SonoffDevice{
			deviceType: BULB,
			deviceId:   id,
			deviceIp:   ip,
			seq:        seq,
			apiVers:    api,
			txtVers:    txt,
		},
		Data: data,
	}
	return d
}

/*
 Run command on this device

 Valid commands are:
   toggle
   power (status)
   onstart (state)
   signal
   pulse (stat, width)
   setup (sid, pass)
   prepare
   update (url, sha256)
   status
*/
func (d *DiyDevice) Run(cmd string, data1 string, data2 string) error {
	log.Printf("run command %s (%s %s) on %s", cmd, data1, data2, d.deviceId)

	var run string
	var request = Request{Deviceid: d.deviceId, Data: Data{}}
	switch cmd {
	case "toggle":
		if d.Data.Switch == "on" {
			request.Data.Switch = "off"
		} else {
			request.Data.Switch = "on"
		}
		run = "switch"
	case "power":
		request.Data.Switch = data1
		run = "switch"
	case "onstart":
		request.Data.Startup = data1 // stay / on / off
		run = "startup"
	case "signal":
		run = "signal_strength"
	case "pulse":
		if n, err := strconv.Atoi(data2); err != nil {
			request.Data.Pulse = data1 // on / off
			request.Data.PulseWidth = n
			run = "pulse"
		} else {
			return fmt.Errorf("wrong value PulseWidth %s", data2)
		}
	case "setup":
		request.Data.Ssid = data1
		request.Data.Password = data2
		run = "wifi"
	case "prepare":
		run = "ota_unlock"
	case "update":
		request.Data.Url = data1       // "http://192.168.1.184/ota/new_rom.bin"
		request.Data.Sha256sum = data2 // "3213b2c34cecbb3bb817030c7f025396b658634c0cf9c4435fc0b52ec9644667"
		run = "ota_flash"
	case "status":
		run = "info"
	}

	req, err := json.Marshal(request)
	if err != nil {
		return err
	}

	res, err := httpRequest("POST", fmt.Sprintf("http://%s:8081/zeroconf/%s", d.deviceIp, run), req)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(res, &response); err != nil {
		return err
	}

	switch response.Error {
	case 400:
		return fmt.Errorf("the request was formatted incorrectly")
	case 401:
		return fmt.Errorf("the request was unauthorized. encryption is enabled, but request is not encrypted")
	case 403:
		return fmt.Errorf("the OTA function was not unlocked")
	case 404:
		return fmt.Errorf("the device does not exist")
	case 408:
		return fmt.Errorf("the pre-download firmware timed out")
	case 413:
		return fmt.Errorf("the request body size is too large")
	case 422:
		return fmt.Errorf("the request parameters are invalid")
	case 424:
		return fmt.Errorf("the firmware could not be downloaded")
	case 471:
		return fmt.Errorf("the firmware integrity check failed")
	case 500:
		return fmt.Errorf("the device has errors")
	case 503:
		return fmt.Errorf("the device is not able to request the vendor's OTA unlock service")
	}

	switch cmd {
	case "signal":
		d.Data.Strength = response.Data.Strength
	case "status":
		d.Data = response.Data
	}

	return nil
}
