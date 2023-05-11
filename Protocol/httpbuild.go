package Protocol

import (
	"bufio"
	"encoding/json"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
	"net/http"
	"strings"
	"time"
)

var StreamPool *tcpassembly.StreamPool

type HttpStreamFactory struct{}

type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *HttpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hStream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hStream.runRequest()
	//go hStream.runResponse()
	return &hStream.r
}

func (h *httpStream) runRequest() {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			return
		} else if err != nil || err == io.ErrUnexpectedEOF {
			go h.runResponse(req)
			return
		} else {
			srcIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)
			Httprep, _ := NewHttpReq(req, srcIp, dstIp, srcPort, dstPort)
			_, _ = json.Marshal(Httprep)
			req.Body.Close()
		}
	}
}

func (h *httpStream) runResponse(req *http.Request) {
	buf := bufio.NewReader(&h.r)
	for {
		resp, err := http.ReadResponse(buf, req)
		if err == io.EOF {
			return
		} else if err != nil || err == io.ErrUnexpectedEOF {
			return
		} else {
			srcIp, dstIp := SplitNet2Ips(h.net)
			srcPort, dstPort := Transport2Ports(h.transport)
			Httprep, _ := NewHttpResp(resp, srcIp, dstIp, srcPort, dstPort)
			_, _ = json.Marshal(Httprep)
			resp.Body.Close()
		}
	}
}

func SplitNet2Ips(net gopacket.Flow) (srcip, dstip string) {
	ips := strings.Split(net.String(), "->")
	if len(ips) > 1 {
		srcip = ips[0]
		dstip = ips[1]
	}
	return srcip, dstip
}

func Transport2Ports(transport gopacket.Flow) (srcport, dstport string) {
	ports := strings.Split(transport.String(), "->")
	if len(ports) > 1 {
		srcport = ports[0]
		dstport = ports[1]
	}
	return srcport, dstport
}

func LayerHttp(packetchan chan gopacket.Packet) {
	streamFactory := &HttpStreamFactory{}
	StreamPool = tcpassembly.NewStreamPool(streamFactory)
	Assembler := tcpassembly.NewAssembler(StreamPool)
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packetchan:
			tcp := packet.TransportLayer().(*layers.TCP)
			Assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			Assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}

}
