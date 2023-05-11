package WebServer

type MyInfo struct {
	Time     string `json:"time"`
	Type     string `json:"type"`
	DestIp   string `json:"destIp"`
	DestName string `json:"destName"`
	DestLocX string `json:"destLocX"`
	DestLocY string `json:"destLocY"`
	SrcIp    string `json:"srcIp"`
	SrcName  string `json:"srcName"`
	SrcLocX  string `json:"srcLocX"`
	SrcLocY  string `json:"srcLocY"`
}
