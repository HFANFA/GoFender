package MessageQueue

import (
	"GoFender/ProcessPacket"
	"GoFender/Utils"
	"GoFender/YamlConfig"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"time"
)

type ConsumerGroupHandler struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	compaket := Utils.CommonPacket{}
	for msg := range claim.Messages() {
		err := json.Unmarshal(msg.Value, &compaket)
		if err != nil {
			log.Println(err)
		}
		ProcessPacket.PacketProcess(compaket)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (h ConsumerGroupHandler) ConsumerPacket(Topic string) {
	//create consumer config
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	//create consumer group
	group, err := sarama.NewConsumerGroup([]string{YamlConfig.Myconfig.KafkaServer}, "PacketProcessGroup", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()
	go func() {
		for err := range group.Errors() {
			log.Println("Consumer Start Error", err)
		}
	}()
	//迭代消费者session
	ctx := context.Background()
	for {
		err := group.Consume(ctx, []string{Topic}, h) //consume message
		if err != nil {
			panic(err)
		}
	}
}
