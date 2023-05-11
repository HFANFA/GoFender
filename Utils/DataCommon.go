package Utils

import (
	"GoFender/GetPacketLayer"
	"GoFender/Protocol"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

type CommonPacket struct {
	ComTime       time.Time
	ComDesIp      string
	ComSrcIp      string
	ComDesPort    string
	ComSrcPort    string
	ComProtocol   string
	ComPacketData []byte
}

type NomPacket struct {
	CommInfo CommonPacket
	Type     string
}

type EvilPacket struct {
	CommInfo   CommonPacket
	Type       string
	AttackType string
}

func (com *CommonPacket) PacketLayer(packet gopacket.Packet) {
	ip := Protocol.IPProtocol{}
	ip.LayerIP(packet)
	com.ComSrcIp = ip.SourceIp.String()
	com.ComDesIp = ip.DestIp.String()
	trans := Protocol.TransportProtocol{}
	if packet.Layer(layers.LayerTypeTCP) != nil {
		trans.LayerTCP(packet)
	} else if packet.Layer(layers.LayerTypeUDP) != nil {
		trans.LayerUDP(packet)
	}
	com.ComDesPort = trans.DstPort
	com.ComSrcPort = trans.SrcPort
	com.ComTime = packet.Metadata().Timestamp
	layertype := &GetPacketLayer.LayerType{}
	layertype.IdentificationProtocol(packet)
	if layertype.ProtocolType != "" {
		com.ComProtocol = layertype.ProtocolType
	} else if layerall := packet.Layers(); layerall[len(layerall)-1].LayerType().String() != "Payload" {
		com.ComProtocol = layerall[len(layerall)-1].LayerType().String()
	} else {
		com.ComProtocol = trans.Protocol
	}
	com.ComPacketData = packet.Data()
}
