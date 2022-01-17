package sonoff

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

var defaultPort = 8081

type (
	Request struct {
		Deviceid string `json:"deviceid"`
		Data     Data   `json:"data"`
	}
	Response struct {
		Seq   int  `json:"seq"`
		Error int  `json:"error"`
		Data  Data `json:"data"`
	}
	Data struct {
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
)

func (e *Device) Run(cmd string, data1 string, data2 string) error {
	log.Printf("run command %s (%s %s) for device: %s", cmd, data1, data2, e.DeviceId)

	var run string
	var request = Request{Deviceid: e.DeviceId, Data: Data{}}
	switch cmd {
	case "toggle":
		if e.Data.Switch == "on" {
			request.Data.Switch = "off"
		} else {
			e.Data.Switch = "on"
		}
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

	res, err := httpRequest("POST", fmt.Sprintf("http://%s:%s/zeroconf/%s", e.Ip, defaultPort, run), req)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(res, &response); err != nil {
		return err
	}

	// todo parse error
	// { "seq": 2, "error": 0, "data": { } }
	/*
		error:
		- 0: successfully
		- 400: The operation failed and the request was formatted incorrectly. The request body is not a valid JSON format.
		- 401: The operation failed and the request was unauthorized. Device information encryption is enabled on the device, but the request is not encrypted.
		- 404: The operation failed and the device does not exist. The device does not support the requested deviceid.
		- 422: The operation failed and the request parameters are invalid. For example, the device does not support setting specific device information.
		- 403: The operation failed and the OTA function was not unlocked. The interface "3.2.6OTA function unlocking" must be successfully called first.
		- 408: The operation failed and the pre-download firmware timed out. You can try to call this interface again after optimizing the network environment or increasing the network speed.
		- 413: The operation failed and the request body size is too large. The size of the new OTA firmware exceeds the firmware size limit allowed by the device.
		- 424: The operation failed and the firmware could not be downloaded. The URL address is unreachable (IP address is unreachable, HTTP protocol is unreachable, firmware does not exist, server does not support Range request header, etc.)
		- 471: The operation failed and the firmware integrity check failed. The SHA256 checksum of the downloaded new firmware does not match the value of the request body's sha256sum field. Restarting the device will cause bricking issue.
		- 500: The operation failed and the device has errors. For example, the device ID or API Key error which is not authenticated by the vendor's OTA unlock service;
		- 503: The operation failed and the device is not able to request the vendor's OTA unlock service. For example, the device is not connected to WiFi, the device is not connected to the Internet, the manufacturer's OTA unlock service is down, etc.
	*/

	// todo send results to mqtt ??
	switch cmd {
	case "signal":
		// response.Data.Strength
	case "status":
		//"data": {
		//	"switch": "off",
		//	"startup": "off",
		//	"pulse": "off",
		//	"pulseWidth": 500,
		//	"ssid": "eWeLink",
		//	"otaUnlock": false,
		//	"fwVersion": "3.5.0",
		//	"deviceid": "100000140e",
		//	"bssid": "ec:17:2f:3d:15:e"
		//}
	}

	return nil
}
