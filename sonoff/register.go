package sonoff

import (
	"encoding/json"
	"fmt"
	"log"
)

type (
	DeviceInfo struct {
		Id     string `json:"deviceid"`
		Key    string `json:"apikey"`
		Chip   string `json:"chipid"`
		Method string `json:"accept"`
	}
	DeviceRegister struct {
		Version float64 `json:"version"`
		Sid     string  `json:"ssid"`
		Pass    string  `json:"password"`
		Ip      string  `json:"serverName"`
		Port    string  `json:"port"`
	}
)

func New(sid string, pass string, ip string, port string) error {
	// getting device info
	deviceInfo, err := httpRequest("GET", "http://10.10.7.1/device", nil)
	if err != nil {
		return err
	}

	log.Println("device info", string(deviceInfo))
	var info DeviceInfo
	if err := json.Unmarshal(deviceInfo, &info); err != nil {
		return err
	}

	// send register command
	register := &DeviceRegister{
		Version: 4,
		Sid:     sid,
		Pass:    pass,
		Ip:      ip,
		Port:    port,
	}

	registerRequest, err := json.Marshal(register)
	if err != nil {
		return err
	}

	deviceRegister, err := httpRequest("POST", "http://10.10.7.1/ap", registerRequest)
	log.Println("device register", string(deviceRegister))

	var resp interface{}
	if err := json.Unmarshal(deviceRegister, &resp); err != nil {
		return err
	}

	e := resp.(map[string]interface{})["error"]
	if e.(float64) == 0 {
		return nil
	}
	return fmt.Errorf("error code %d", e)
}
