package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type IPProtocol struct {
	SourceIp net.IP
	DestIp   net.IP
	Protocol string
}

func (ip *IPProtocol) LayerIP(Packet gopacket.Packet) {
	ipv4 := Packet.Layer(layers.LayerTypeIPv4)
	if ipv4 != nil {
		ipv, _ := ipv4.(*layers.IPv4)
		ip.SourceIp = ipv.SrcIP
		ip.DestIp = ipv.DstIP
		ip.Protocol = ipv.Protocol.String()
	}
}
