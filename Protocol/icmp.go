package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var ICMPChan chan gopacket.Packet

type ICMP struct {
	SrcMac   string `json:"src_mac"`
	DstMac   string `json:"dst_mac"`
	SrcIP    string `json:"src_ip"`
	DstIP    string `json:"dst_ip"`
	TypeCode string `json:"type_code"`
	Data     string `json:"data"`
}

func (I *ICMP) LayerIcmp(packet gopacket.Packet) {
	icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
	if icmpLayer != nil {
		I.SrcMac = packet.LinkLayer().LinkFlow().Src().String()
		I.DstMac = packet.LinkLayer().LinkFlow().Dst().String()
		I.SrcIP = packet.NetworkLayer().NetworkFlow().Src().String()
		I.DstIP = packet.NetworkLayer().NetworkFlow().Dst().String()
		icmp := icmpLayer.(*layers.ICMPv4)
		I.TypeCode = icmp.TypeCode.String()
		I.Data = string(icmp.Payload)
	}
}
