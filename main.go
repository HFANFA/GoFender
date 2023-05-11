package main

import (
	"GoFender/Database"
	"GoFender/GetPacket"
	"GoFender/MessageQueue"
	"GoFender/SuricataMatch"
	"GoFender/WebServer"
	"GoFender/YamlConfig"
	"github.com/google/gopacket"
	_ "net/http/pprof"
)

func main() {
	go WebServer.WebStart()
	PacketChan := make(chan gopacket.Packet, 10240)
	defer close(PacketChan)
	go func() {
		err := GetPacket.CapturePackets(PacketChan)
		if err != nil {
			panic(err)
		}
	}()
	if PacketChan != nil {
		handel := MessageQueue.ConsumerGroupHandler{}
		go func() {
			for {
				Packet := <-PacketChan
				MessageQueue.ProducerPacket(Packet, "Packet_Data")
			}
		}()
		handel.ConsumerPacket("Packet_Data")
	}
}

func init() {
	YamlConfig.Myconfig = YamlConfig.ParseYaml("./config.yaml")
	MessageQueue.Producer = MessageQueue.InitKafka()
	Database.DataPool = Database.InitDatabase()
	SuricataMatch.RuleSet = SuricataMatch.RulesParse(YamlConfig.Myconfig.RulesPath)
	SuricataMatch.ACTrie = SuricataMatch.BulidTrie(SuricataMatch.RuleSet)

}
