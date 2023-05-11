package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type Ethernets struct {
	SourceMAC net.HardwareAddr `json:"source_mac"`
	TargetMAC net.HardwareAddr `json:"target_mac"`
}

func (E *Ethernets) LayerEthernet(Packet gopacket.Packet) {
	ether := Packet.Layer(layers.LayerTypeEthernet)
	if ether != nil {
		ethernetLayer := ether.(*layers.Ethernet)
		E.SourceMAC = ethernetLayer.SrcMAC
		E.TargetMAC = ethernetLayer.DstMAC
	}
}
