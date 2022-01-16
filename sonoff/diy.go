package sonoff

import "log"

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
		OtaUnlock  bool   `json:"otaUnlock,omitempty"`
		FwVersion  string `json:"fwVersion,omitempty"`
		Deviceid   string `json:"deviceid,omitempty"`
		Bssid      string `json:"bssid,omitempty"`
	}
)

func (e *Entry) Run(cmd string, data1 string, data2 string) {
	log.Printf("run command %s (%s %s) for device: %s", cmd, data1, data2, e.DeviceId)
}

/*
http://[ip]:[port]/zeroconf/switch
{
    "deviceid": "",
    "data": {
        "switch": "on"
	}
}
return
{
	"seq": 2,
    "error": 0,
    "data": { }
}

error:
- 0: successfully
- 400: The operation failed and the request was formatted incorrectly. The request body is not a valid JSON format.
- 401: The operation failed and the request was unauthorized. Device information encryption is enabled on the device, but the request is not encrypted.
- 404: The operation failed and the device does not exist. The device does not support the requested deviceid.
- 422: The operation failed and the request parameters are invalid. For example, the device does not support setting specific device information.

*/
func Toggle() {
}

/*
http://[ip]:[port]/zeroconf/startup
{
    "deviceid": "",
    "data": {
        "startup": "stay"  // "on" "off"
	}
}

return
{
	"seq": 2,
    "error": 0,
    "data": { }
}
*/
func PowerOn() {
}

/*
http://[ip]:[port]/zeroconf/signal_strength
{
    "deviceid": "",
	"data": { }
}

return
{
	"seq": 2,
    "error": 0,
    "data": {
        "signalStrength": -67
    }
}

*/
func Signal() {

}

/*
http://[ip]:[port]/zeroconf/pulse
{
    "deviceid": "",
    "data": {
        "pulse": "on",
        "pulseWidth": 2000		// 500~36000000
    }
}

return
{
	"seq": 2,
    "error": 0,
    "data": { }
}
*/
func Pulse() {

}

/*
http://[ip]:[port]/zeroconf/wifi
{
    "deviceid": "",
    "data": {
        "ssid": "eWeLink",
        "password": "WeLoveIoT"
	}
}
return
{
	"seq": 2,
    "error": 0,
    "data": { }
}
*/
func Wifi() {

}

/*
http://[ip]:[port]/zeroconf/ota_unlock
{
    "deviceid": "",
	"data": { }
}

return
- 500: The operation failed and the device has errors. For example, the device ID or API Key error which is not authenticated by the vendor's OTA unlock service;
- 503: The operation failed and the device is not able to request the vendor's OTA unlock service. For example, the device is not connected to WiFi, the device is not connected to the Internet, the manufacturer's OTA unlock service is down, etc.
*/
func OtaUnlock() {

}

/*
http://[ip]:[port]/zeroconf/ota_flash
{
    "deviceid": "",
    "data": {
        "downloadUrl": "http://192.168.1.184/ota/new_rom.bin",
        "sha256sum": "3213b2c34cecbb3bb817030c7f025396b658634c0cf9c4435fc0b52ec9644667"
    }
}

return http error code

- 403: The operation failed and the OTA function was not unlocked. The interface "3.2.6OTA function unlocking" must be successfully called first.
- 408: The operation failed and the pre-download firmware timed out. You can try to call this interface again after optimizing the network environment or increasing the network speed.
- 413: The operation failed and the request body size is too large. The size of the new OTA firmware exceeds the firmware size limit allowed by the device.
- 424: The operation failed and the firmware could not be downloaded. The URL address is unreachable (IP address is unreachable, HTTP protocol is unreachable, firmware does not exist, server does not support Range request header, etc.)
- 471: The operation failed and the firmware integrity check failed. The SHA256 checksum of the downloaded new firmware does not match the value of the request body's sha256sum field. Restarting the device will cause bricking issue.
*/
func OtaFlash() {

}

/*
http://[ip]:[port]/zeroconf/info
{
    "deviceid": "",
	"data": { }
}

return json
{
	"seq": 2,
    "error": 0,
    "data": {
		"switch": "off",
		"startup": "off",
		"pulse": "off",
		"pulseWidth": 500,
		"ssid": "eWeLink",
		"otaUnlock": false,
		"fwVersion": "3.5.0",
		"deviceid": "100000140e",
		"bssid": "ec:17:2f:3d:15:e"
	}
}
*/
func Info() {
}
