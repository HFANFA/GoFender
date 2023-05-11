package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var TransLayerChan chan gopacket.Packet

type TransportProtocol struct {
	Protocol string `json:"protocol"`
	SrcPort  string `json:"src_port"`
	DstPort  string `json:"dst_port"`
	Payload  string `json:"payload"`
}

func (pro *TransportProtocol) LayerTrans(PacketChan chan gopacket.Packet) {
	Packet := <-PacketChan
	if Packet.Layer(layers.LayerTypeTCP) != nil {
		tcp := TransportProtocol{}
		tcp.LayerTCP(Packet)
	} else if Packet.Layer(layers.LayerTypeUDP) != nil {
		udp := TransportProtocol{}
		udp.LayerUDP(Packet)
	}
}

func (pro *TransportProtocol) LayerTCP(Packet gopacket.Packet) {
	if Packet.Layer(layers.LayerTypeTCP) != nil {
		pro.Protocol = "TCP"
		tcp := Packet.Layer(layers.LayerTypeTCP)
		if tcp != nil {
			tcpLayer := tcp.(*layers.TCP)
			pro.SrcPort = tcpLayer.SrcPort.String()
			pro.DstPort = tcpLayer.DstPort.String()
		}
		payload := Packet.ApplicationLayer()
		if payload != nil {
			pro.Payload = string(payload.Payload())
		}
	}
}

func (pro *TransportProtocol) LayerUDP(Packet gopacket.Packet) {
	if Packet.Layer(layers.LayerTypeUDP) != nil {
		pro.Protocol = "UDP"
		udp := Packet.Layer(layers.LayerTypeUDP)
		if udp != nil {
			udpLayer := udp.(*layers.UDP)
			pro.SrcPort = udpLayer.SrcPort.String()
			pro.DstPort = udpLayer.DstPort.String()
		}
		payload := Packet.ApplicationLayer()
		if payload != nil {
			pro.Payload = string(payload.Payload())
		}
	}
}
