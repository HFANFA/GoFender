package Protocol

import (
	"github.com/google/gopacket"
	"io"
	"net/http"
)

var HttpChan chan gopacket.Packet

type HttpReq struct {
	SrcIP   string
	DstIP   string
	SrcPort string
	DstPort string
	Host    string
	Method  string
	Path    string
	Proto   string
	Header  http.Header
	Body    string
}

type HttpResp struct {
	SrcIP   string
	DstIP   string
	SrcPort string
	DstPort string
	Proto   string
	Status  string
	Header  http.Header
	Body    string
}

func NewHttpReq(req *http.Request, SIP string, DIP string, SPort string, DPort string) (httpReq *HttpReq, err error) {
	body := req.Body
	buff, _ := io.ReadAll(body)
	//bodycache, _ := utf8.DecodeRune(buff)
	return &HttpReq{
		SrcIP:   SIP,
		SrcPort: SPort,
		DstIP:   DIP,
		DstPort: DPort,
		Host:    req.Host,
		Method:  req.Method,
		Path:    req.URL.Path,
		Proto:   req.Proto,
		Header:  req.Header,
		Body:    string(buff),
	}, err
}

func NewHttpResp(resp *http.Response, SIP string, DIP string, SPort string, DPort string) (httpReq *HttpResp, err error) {
	body := resp.Body
	buff, _ := io.ReadAll(body)
	//bodycache, _ := utf8.DecodeRune(buff)
	return &HttpResp{
		SrcIP:   SIP,
		SrcPort: SPort,
		DstIP:   DIP,
		DstPort: DPort,
		Proto:   resp.Proto,
		Status:  resp.Status,
		Header:  resp.Header,
		Body:    string(buff),
	}, err
}
