package Protocol

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Tls struct {
	Version      layers.TLSVersion
	Type         layers.TLSType
	EncryptedMsg []byte
}

var TlsChan chan gopacket.Packet

func (t *Tls) LayerTls(packet gopacket.Packet) {
	tls := layers.TLS{
		BaseLayer:        layers.BaseLayer{},
		ChangeCipherSpec: nil,
		Handshake:        nil,
		AppData:          nil,
		Alert:            nil,
	}
	if tlspk := packet.Layer(layers.LayerTypeTLS); tlspk != nil {
		err := tls.DecodeFromBytes(tlspk.(*layers.TLS).LayerPayload(), gopacket.NilDecodeFeedback)
		if err != nil {
			return
		}
		if tls.AppData != nil {
			t.Version = tls.AppData[0].Version
			t.Type = tls.AppData[0].ContentType
			t.EncryptedMsg = tls.AppData[0].Payload
			return
		}
		if tls.Alert != nil {
			t.Version = tls.Alert[0].Version
			t.Type = tls.Alert[0].ContentType
			t.EncryptedMsg = tls.Alert[0].EncryptedMsg
			return
		}
	}
}
