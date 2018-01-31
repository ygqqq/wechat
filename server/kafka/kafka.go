package kafka

import(
    "github.com/Shopify/sarama"
    "encoding/json"
    "fmt"
)

var (
    brokers = []string{"127.0.0.1:9092"}
    topic   = "test"
    topics  = []string{topic}
    producer sarama.SyncProducer = nil
    consumer sarama.PartitionConsumer = nil
)

func getConf() *sarama.Config {
    conf := sarama.NewConfig()
    conf.Producer.RequiredAcks = sarama.WaitForAll
    conf.Producer.Return.Successes = true
    conf.ChannelBufferSize = 1
    conf.Version = sarama.V1_0_0_0
    return conf
}
func getProducer() sarama.SyncProducer{
    if producer == nil {
        producer, _ = sarama.NewSyncProducer(brokers, getConf())
    }
    return producer
}


func GetConsumer() sarama.PartitionConsumer {
    if consumer == nil {
        cs, _ := sarama.NewConsumer(brokers, getConf())
        consumer, _ = cs.ConsumePartition(topic, 0, sarama.OffsetNewest)
    }
    return consumer
}
func SendToKafka(msg interface{}) error {
    pro := getProducer()
    json, err := json.Marshal(msg)

    if err != nil {
        return err
    }

    msgLog := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(string(json)),
    }
    _, _, err = pro.SendMessage(msgLog)
    if err != nil {
        fmt.Printf("Kafka error: %s\n", err)
    }
    return nil
}
