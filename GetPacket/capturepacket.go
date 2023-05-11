package GetPacket

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

func CapturePackets(PacketChan chan gopacket.Packet) error {
	DeviceName := findDevice()
	inactive, err := pcap.NewInactiveHandle(DeviceName)
	if err != nil {
		log.Fatal(err)
	}
	defer inactive.CleanUp()
	if err = inactive.SetTimeout(time.Minute); err != nil {
		log.Fatal(err)
	} else if err = inactive.SetImmediateMode(true); err != nil {
		log.Fatal(err)
	}
	handle, err1 := inactive.Activate()
	if err1 != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	SourcePaket := gopacket.NewPacketSource(handle, handle.LinkType())
	defer close(PacketChan)
	for {
		packet, err := SourcePaket.NextPacket()
		if err == nil {
			PacketChan <- packet
			continue
		}
	}
}
