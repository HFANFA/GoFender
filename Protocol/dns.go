package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

var DnsChan chan gopacket.Packet

type Dns struct {
	SrcIP    net.IP `json:"src_ip"`
	DstIP    net.IP `json:"dst_ip"`
	SrcPort  string `json:"src_port"`
	DstPort  string `json:"dst_port"`
	Question []DnsQuestion
	Answer   []DnsAnswer
}

type DnsQuestion struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type DnsAnswer struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (dns *Dns) LayerDNS(packetChan chan gopacket.Packet) {
	for {
		packet := <-packetChan
		ip := IPProtocol{Protocol: "UDP"}
		ip.LayerIP(packet)
		udp := TransportProtocol{Protocol: "UDP"}
		udp.LayerUDP(packet)
		dns.SrcIP = ip.SourceIp
		dns.DstIP = ip.DestIp
		dns.SrcPort = udp.SrcPort
		dns.DstPort = udp.DstPort
		dnsLayer := packet.Layer(layers.LayerTypeDNS)
		if dnsLayer != nil {
			dnsPacket := dnsLayer.(*layers.DNS)
			for _, question := range dnsPacket.Questions {
				dns.Question = append(dns.Question, DnsQuestion{
					Type: question.Type.String(),
					Name: string(question.Name),
				})
			}
			for _, answer := range dnsPacket.Answers {
				dns.Answer = append(dns.Answer, DnsAnswer{
					Type: answer.Type.String(),
					Name: string(answer.Name),
				})
			}
		}
	}
}
