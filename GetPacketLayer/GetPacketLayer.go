package GetPacketLayer

import (
	godpi "github.com/GreyNoise-Intelligence/go-dpi"
	"github.com/google/gopacket"
	"log"
)

type LayerType struct {
	ProtocolType string
}

func (lt *LayerType) IdentificationProtocol(packet gopacket.Packet) {
	godpi.Initialize()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			lt.ProtocolType = ""
		}
	}()
	defer godpi.Destroy()
	flow, _ := godpi.GetPacketFlow(packet)
	result := godpi.ClassifyFlow(flow)
	if result.Protocol != "" {
		lt.ProtocolType = string(result.Protocol)
	}

}
