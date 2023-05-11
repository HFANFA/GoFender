package GetPacket

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"regexp"
)

var DeviceName string

func findDevice() string {
	reg := regexp.MustCompile(`.*(Wireless Network|Family Controller|eth0|en0).*`)
	if reg == nil {
		log.Fatal("regexp error")
	}
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		find := reg.FindAllString(device.Description+" "+device.Name, -1)
		if len(find) > 0 {
			IP := device.Addresses[1].IP
			Ip := IP.To4()
			if !IP.IsLoopback() || !IP.IsLinkLocalMulticast() || !IP.IsLinkLocalUnicast() || Ip[0] == 169 {
				fmt.Println("[+] Find activeted device network interface: ", device.Description)
				fmt.Println("[+] Name: ", device.Name)
				fmt.Println("[+] MAC addresses: ", device.Addresses[0].IP.String())
				fmt.Println("[+] IP address: ", device.Addresses[1].IP.String())
				DeviceName = device.Name
			}
			break
		} else {
			continue
		}
	}
	if DeviceName == "" {
		fmt.Println("[-] Can't find the device Network Interface!")
		os.Exit(1)
	}
	return DeviceName
}
