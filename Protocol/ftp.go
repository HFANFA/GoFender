package Protocol

import (
	"github.com/google/gopacket"
	"net"
)

var FTPChan chan gopacket.Packet

type FTP struct {
	SrcIP   net.IP `json:"src_ip"`
	DstIP   net.IP `json:"dst_ip"`
	SrcPort string `json:"src_port"`
	DstPort string `json:"dst_port"`
	Data    string `json:"data"`
}
type FTPData struct {
	SrcIP   net.IP `json:"src_ip"`
	DstIP   net.IP `json:"dst_ip"`
	SrcPort string `json:"src_port"`
	DstPort string `json:"dst_port"`
	FTPData string `json:"ftp_data"`
}

func (ftp *FTP) LayerFtp(packetChan chan gopacket.Packet) {
	for {
		packet := <-packetChan
		tcp := TransportProtocol{}
		tcp.LayerTCP(packet)
		ip := IPProtocol{Protocol: "TCP"}
		ip.LayerIP(packet)
		ftp.SrcIP = ip.SourceIp
		ftp.DstIP = ip.DestIp
		ftp.SrcPort = tcp.SrcPort
		ftp.DstPort = tcp.DstPort
		ftp.Data = tcp.Payload
	}
}

func (ftp *FTPData) LayerFtpData(packetChan chan gopacket.Packet) {
	for {
		packet := <-packetChan
		ip := IPProtocol{Protocol: "TCP"}
		ip.LayerIP(packet)
		tcp := TransportProtocol{}
		tcp.LayerTCP(packet)
		applications := packet.ApplicationLayer()
		if applications != nil {
			payload := applications.Payload()
			if payload != nil || len(payload) > 0 {
				ftp.SrcIP = ip.SourceIp
				ftp.DstIP = ip.DestIp
				ftp.SrcPort = tcp.SrcPort
				ftp.DstPort = tcp.DstPort
				ftp.FTPData = string(applications.Payload())
			}
			return
		}
	}

}
