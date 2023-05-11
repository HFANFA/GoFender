package Utils

import (
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/oschwald/geoip2-golang"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Location struct {
	Locationinfo string
	LocX         float64
	LocY         float64
}

type IpToLocation struct {
	SrcIP           string
	SrcPhyLocation  string
	SrcLocX         float64
	SrcLocY         float64
	DestIP          string
	DestPhyLocation string
	DestLocX        float64
	DestLocY        float64
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsMulticast() || IP.IsLinkLocalUnicast() || IP.IsUnspecified() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func GetPublicIp() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal("Can't get public IP, Internet error !", err)
		}
	}(resp.Body)
	content, _ := io.ReadAll(resp.Body)
	return string(content)
}

func (l *Location) GetPhyInfo(Ip string) {
	db, err := geoip2.Open("./etc/LocationDatabase/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	ip := make(net.IP, 1)
	ip = net.ParseIP(Ip)
	//backup ip to location
	if !IsPublicIP(ip) {
		db, err := ip2location.OpenDB("./etc/LocationDatabase/IP2LOCATION-LITE-DB5.BIN")
		if err != nil {
			log.Fatal(err)
			return
		}
		Result, _ := db.Get_all(GetPublicIp())
		l.Locationinfo = "本地局域网"
		l.LocX, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(Result.Longitude)), 64)
		l.LocY, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", float64(Result.Latitude)), 64)
		return
	}
	Results, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	if len(Results.Country.IsoCode) > 1 {
		l.Locationinfo = Results.Country.Names["zh-CN"] + "(" + Results.Country.IsoCode + ")" //国家
		if len(Results.Subdivisions) > 0 {
			if len(Results.Subdivisions[0].Names["zh-CN"]) > 0 && Results.Country.IsoCode != "CN" {
				l.Locationinfo += " " + Results.Subdivisions[0].Names["zh-CN"] + "(" + Results.Subdivisions[0].Names["en"] + ")"
			} else if len(Results.Subdivisions[0].Names["zh-CN"]) > 0 && Results.Country.IsoCode == "CN" {
				l.Locationinfo += " " + Results.Subdivisions[0].Names["zh-CN"]
			} else {
				l.Locationinfo += " " + Results.Subdivisions[0].Names["en"]
			}
		}
		if len(Results.City.Names["zh-CN"]) > 0 {
			l.Locationinfo += " " + Results.City.Names["zh-CN"]
		} else if len(Results.City.Names["en"]) > 0 {
			l.Locationinfo += " " + Results.City.Names["en"]
		}
		l.LocX = Results.Location.Longitude
		l.LocY = Results.Location.Latitude
	}
}

func (ipl *IpToLocation) GetIpInfo(Srcip, Destip string) {
	//source ip infomaion
	ipl.SrcIP = Srcip
	srcLocinfo := Location{}
	srcLocinfo.GetPhyInfo(Srcip)
	ipl.SrcPhyLocation = srcLocinfo.Locationinfo
	ipl.SrcLocX = srcLocinfo.LocX
	ipl.SrcLocY = srcLocinfo.LocY
	//destination ip infomaion
	ipl.DestIP = Destip
	destLocinfo := Location{}
	destLocinfo.GetPhyInfo(Destip)
	ipl.DestPhyLocation = destLocinfo.Locationinfo
	ipl.DestLocX = destLocinfo.LocX
	ipl.DestLocY = destLocinfo.LocY
}
