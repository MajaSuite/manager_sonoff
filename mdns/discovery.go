package mdns

import (
	"fmt"
	"github.com/miekg/dns"
	"golang.org/x/net/ipv4"
	"log"
	"manager_sonoff/sonoff"
	"manager_sonoff/utils"
	"net"
	"strings"
	"time"
)

const (
	domain = "_ewelink._tcp.local."
)

var (
	discoveryAddr = "224.0.0.251"
	discoveryPort = 5353
)

type Discovery struct {
}

func NewDiscovery(debug bool, chanel chan sonoff.Device) error {
	d := &Discovery{}

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", discoveryAddr, discoveryPort))
	if err != nil {
		return err
	}

	mconn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return err
	}
	mpkt := ipv4.NewPacketConn(mconn)

	go d.notifier(mconn, addr)

	for {
		buffer := make([]byte, 65536)
		_, _, src, err := mpkt.ReadFrom(buffer)
		if err != nil {
			log.Printf("failed to read packet: %v", err)
			continue
		}

		msg := new(dns.Msg)
		if err := msg.Unpack(buffer); err != nil {
			log.Printf("failed to unpack packet: %v", err)
			continue
		}

		var data, iv, id, tp string
		var seq, api, txt int
		var enc bool
		var ip = strings.Split(src.String(), ":")[0]
		var service string
		for _, answer := range append(msg.Answer, msg.Extra...) {
			switch answer.Header().Rrtype {
			case dns.TypeTXT:
				service = answer.(*dns.TXT).Hdr.Name
				for _, pair := range answer.(*dns.TXT).Txt {
					s := strings.Split(pair, "=")
					if len(s) > 1 {
						v := pair[len(s[0])+1:]
						switch s[0] {
						case "data1":
							data = v
						case "iv":
							iv = v
						case "encrypt":
							enc = utils.ConvertBool(v)
						case "seq":
							seq = utils.ConvertInt(v)
						case "id":
							id = v
						case "apivers":
							api = utils.ConvertInt(v)
						case "txtvers":
							txt = utils.ConvertInt(v)
						case "type":
							tp = v
						}
					}
				}
			}
		}

		if strings.HasSuffix(service, domain) {
			switch tp {
			case "light":
				dev := sonoff.NewBulb(id, ip, api, txt, seq, data, iv, enc)
				if dev != nil {
					chanel <- dev
				}
			case "diy":
				dev := sonoff.NewDIY(id, ip, api, txt, seq, data)
				if dev != nil {
					chanel <- dev
				}
			case "socket":
				//...
			}
		}
	}

	return nil
}

func (d *Discovery) notifier(conn *net.UDPConn, addr *net.UDPAddr) {
	for {
		m := new(dns.Msg)
		m.SetQuestion(domain, dns.TypePTR)
		m.RecursionDesired = false

		buf, err := m.Pack()
		if err != nil {
			log.Println("error pack query", err)
		}

		_, err = conn.WriteToUDP(buf, addr)
		if err != nil {
			log.Println("error write udp", err)
		}

		time.Sleep(time.Second * 60)
	}
}
