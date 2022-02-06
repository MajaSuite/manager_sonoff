package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/client"
	"github.com/MajaSuite/mqtt/packet"
	"log"
	"manager_sonoff/mdns"
	"manager_sonoff/sonoff"
)

var (
	srv        = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid   = flag.String("clientid", "sonoff-1", "client id for mqtt server")
	keepalive  = flag.Int("keepalive", 30, "keepalive timeout for mqtt server")
	login      = flag.String("login", "", "login string for mqtt server")
	pass       = flag.String("pass", "", "password string for mqtt server")
	debug      = flag.Bool("debug", false, "print debuging hex dumps")
	reg        = flag.Bool("reg", false, "to sonoff new device")
	sid        = flag.String("sid", "myhome", "network name for registration")
	key        = flag.String("key", "mypass", "network key for registration")
	serverIp   = flag.String("ip", "192.168.1.1", "sonoff server ip")
	serverPort = flag.String("port", "8081", "sonoff server port")
)

func main() {
	flag.Parse()

	if *reg {
		log.Println("sonoff new device registration")
		sonoff.NewRegistration(*sid, *key, *serverIp, *serverPort)
		return
	}

	log.Println("starting manager_sonoff")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttId uint16 = 1
	mqtt, err := client.Connect(*srv, *clientid, uint16(*keepalive), false, *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server ")
	}

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttId
	sp.Topics = []packet.SubscribePayload{{Topic: "sonoff/#", QoS: 1}}
	mqtt.Send <- sp
	mqttId++

	devices := make(map[string]sonoff.Device)

	log.Println("start mDNS discovery")
	discovery := make(chan sonoff.Device)
	go mdns.NewDiscovery(*debug, discovery)
	if err != nil {
		panic(err)
	}

	// main cycle
	for {
		select {
		case pkt := <-mqtt.Receive:
			if pkt.Type() == packet.PUBLISH {
				var dev sonoff.Device
				//topics := strings.Split(pkt.(*packet.PublishPacket).Topic, "/")
				if err := json.Unmarshal([]byte(pkt.(*packet.PublishPacket).Payload), &dev); err == nil {
					// todo change logic
					//if dev.Cmd != "" {
					//	if devices[topics[1]] != nil {
					//		if err := devices[topics[1]].Run(dev.Cmd, dev.Data1, dev.Data2); err != nil {
					//			log.Println("error running command", err)
					//		} else {
					//			// send results back to mqtt server
					//		}
					//	}
					//} else {
					//	// restore devices from mqtt
					//	if devices[dev.ID()] == nil {
					//		log.Printf("restore device %s", dev.ID())
					//		devices[dev.ID()] = dev
					//	}
					//}
				} else {
					log.Println("error unmarshal request from mqtt", err)
				}
			}
		case dev := <-discovery:
			p := packet.NewPublish()
			p.Id = mqttId
			p.Topic = fmt.Sprintf("sonoff/%s", dev.ID())
			p.QoS = 1

			switch dev.Type() {
			case sonoff.BULB:
				// todo
			case sonoff.DIY:
				if devices[dev.ID()] == nil {
					log.Printf("new device: %s", dev)
					devices[dev.ID()] = dev
					p.Retain = true
				} else {
					devices[dev.ID()] = dev
				}
				p.Payload = dev.(*sonoff.DiyDevice).String()
			case sonoff.SOCKET:
				// todo waits for device
			}

			if p.Payload != "" {
				mqtt.Send <- p
				mqttId++
			}
		}
	}
}
