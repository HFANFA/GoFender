package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

var ArpChan chan gopacket.Packet

type Arp struct {
	SenderMac net.HardwareAddr `json:"sender_mac"`
	TargetMac net.HardwareAddr `json:"target_mac"`
	SenderIp  net.IP           `json:"sender_ip"`
	TargetIp  net.IP           `json:"target_ip"`
}

func (A *Arp) LayerArp(packetChan chan gopacket.Packet) {
	for {
		packet := <-packetChan
		arpLayer := packet.Layer(layers.LayerTypeARP)
		if arpLayer != nil {
			arp := arpLayer.(*layers.ARP)
			A.SenderMac = arp.SourceHwAddress
			A.TargetMac = arp.DstHwAddress
			A.SenderIp = arp.SourceProtAddress
			A.TargetIp = arp.DstProtAddress
		}
	}
}
