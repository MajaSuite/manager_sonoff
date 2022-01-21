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

var (
	discoveryAddr = "224.0.0.251"
	discoveryPort = 5353
)

type Discovery struct {
	domain   string
	Reporter chan *sonoff.Device
}

func NewDiscovery(domain string) (*Discovery, error) {
	d := &Discovery{
		domain:   domain,
		Reporter: make(chan *sonoff.Device),
	}

	uconn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	upkt := ipv4.NewPacketConn(uconn)

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", discoveryAddr, discoveryPort))
	if err != nil {
		return nil, err
	}

	mconn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}
	mpkt := ipv4.NewPacketConn(mconn)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		ifAddr, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, ifaddr := range ifAddr {
			// skip ipv6
			if strings.Contains(ifaddr.String(), "::") {
				continue
			}

			// skip localhost
			if strings.Contains(ifaddr.String(), "127.0.0.1") {
				continue
			}

			ifAddr := strings.Split(ifaddr.String(), "/")

			log.Printf("join to %s interface with ip %s", iface.Name, ifAddr[0])

			if err := upkt.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(discoveryAddr)}); err != nil {
				log.Println("join error", err)
			}

			go d.notifier(mconn, addr)
			go d.listener(mpkt)
		}
	}

	return d, nil
}

func (d *Discovery) listener(mpkt *ipv4.PacketConn) error {
	buffer := make([]byte, 65536)
	for {
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

		device := sonoff.Device{}
		device.Ip = strings.Split(src.String(), ":")[0]
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
							device.Data1 = v
						case "iv":
							device.IV = v
						case "encrypt":
							device.Encrypt = utils.ConvertBool(v)
						case "seq":
							device.Seq = utils.ConvertInt(v)
						case "id":
							device.DeviceId = v
						case "apivers":
							device.ApiVers = utils.ConvertInt(v)
						case "txtvers":
							device.TxtVers = utils.ConvertInt(v)
						case "type":
							device.Type = v
						}
					}
				}
			}
		}

		if strings.HasSuffix(service, d.domain) {
			d.Reporter <- &device
		}
	}
}

func (d *Discovery) notifier(conn *net.UDPConn, addr *net.UDPAddr) {
	for {
		m := new(dns.Msg)
		m.SetQuestion(d.domain, dns.TypePTR)
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
