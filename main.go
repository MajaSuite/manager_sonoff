package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/packet"
	"github.com/MajaSuite/mqtt/transport"
	"log"
	"manager_sonoff/mdns"
	"manager_sonoff/sonoff"
	"strings"
	"time"
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
		sonoff.New(*sid, *key, *serverIp, *serverPort)
		return
	}

	log.Println("starting manager_sonoff ...")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttId uint16 = 1
	mqtt, err := transport.Connect(*srv, *clientid, uint16(*keepalive), *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}
	go mqtt.Start()

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttId
	sp.Topics = []packet.SubscribePayload{{Topic: "sonoff/#", QoS: 1}}
	mqtt.Sendout <- sp
	mqttId++

	devices := make(map[string]*sonoff.Device)

	// fetch command data from mqtt server
	go func() {
		for {
			for pkt := range mqtt.Broker {
				if pkt.Type() == packet.PUBLISH {
					var entry sonoff.Device
					topics := strings.Split(pkt.(*packet.PublishPacket).Topic, "/")
					if err := json.Unmarshal([]byte(pkt.(*packet.PublishPacket).Payload), &entry); err == nil {
						if entry.Cmd != "" {
							if devices[topics[1]] != nil {
								if err := devices[topics[1]].Run(entry.Cmd, entry.Data1, entry.Data2); err != nil {
									log.Println("error running command", err)
								} else {
									// send results back to mqtt server
								}
							}
						} else {
							// restore devices from mqtt
							if devices[entry.DeviceId] == nil {
								log.Printf("restore device %s", entry.DeviceId)
								devices[entry.DeviceId] = &entry
							}
						}
					} else {
						log.Println(err)
					}
				}
			}
		}
	}()

	time.Sleep(3 * time.Second)

	log.Println("start mDNS discovery")
	discovery, err := mdns.NewDiscovery("_ewelink._tcp.local.")
	if err != nil {
		panic(err)
	}

	// receive updates from devices
	for entry := range discovery.Reporter {
		if entry.Encrypt == true {
			// todo ??? are you kidding?
			// [data1=4Tp/SNAMhzqOdUaRY5baGfQx1MQA7Q615K5lOk2+csHwnVOelBpXMUjWb2tpCQ+HWyxnjv5yCNrwGUlQg/BOmw==
			//		iv=NTc0NjU2MzIwMTg0NTc5Mg== encrypt=true seq=1 id=10010ac611 apivers=1 type=light txtvers=1]
		} else {
			var data1 sonoff.Data
			if err := json.Unmarshal([]byte(entry.Data1), &data1); err == nil {
				entry.Data = data1
			}
		}

		p := packet.NewPublish()
		p.Id = mqttId
		p.Topic = fmt.Sprintf("sonoff/%s", entry.DeviceId)
		p.QoS = 1
		p.Payload = entry.String()

		if devices[entry.DeviceId] == nil {
			log.Printf("new device: %s", entry.String())
			p.Retain = true
			mqtt.Sendout <- p
			mqttId++
		} else {
			if entry.Seq > devices[entry.DeviceId].Seq {
				log.Println("update device", entry.String())
				devices[entry.DeviceId] = entry

				mqtt.Sendout <- p
				mqttId++
			}
		}
	}
}
