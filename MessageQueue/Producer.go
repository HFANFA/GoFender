package MessageQueue

import (
	"GoFender/Utils"
	"GoFender/YamlConfig"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/gopacket"
	"log"
)

var Producer sarama.SyncProducer

func InitKafka() sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.NoResponse
	config.Producer.Return.Successes = true
	Producer, err := sarama.NewSyncProducer([]string{YamlConfig.Myconfig.KafkaServer}, config)
	/*defer func() {
		_ = Producer.Close()
	}()*/
	if err != nil {
		panic(err.Error())
	}
	return Producer
}

func ProducerPacket(packet gopacket.Packet, Topic string) {
	comm := &Utils.CommonPacket{}
	comm.PacketLayer(packet)
	packetmsg, _ := json.Marshal(comm)
	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.StringEncoder(packetmsg),
	}
	_, _, err := Producer.SendMessage(msg)
	if err != nil {
		log.Println("send msg failed, err:", err)
		return
	}
}
