package kafka

import(
    //"github.com/Shopify/sarama"
)
type U struct{
    Name string `json:"name"`
    Age int `json:"age"`
    encoded []byte
    err error
}
func test(){
    // sarama.AsyncProducer.Input() <- &sarama.ProducerMessage{
    //     Topic: "test",
    //     Key: nil,
         
    //     Value: &U{"ygq",18}, 
    // }
}
// func sendMsg(){
// 	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer func() {
// 		if err := producer.Close(); err != nil {
// 			log.Fatalln(err)
// 		}
// 	}()

// 	msg := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder("testing 123")}
// 	partition, offset, err := producer.SendMessage(msg)
// 	if err != nil {
// 		log.Printf("FAILED to send message: %s\n", err)
// 	} else {
// 		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
// 	}
// }

func getMsg(){
	// sarama.Logger = logger
	// // 连接kafka消息服务器
    // consumer, err := sarama.NewConsumer(strings.Split(kafka, ","), nil)
    // if err != nil {
    //     logger.Printf("Failed to start consumer: %s", err)
    // }
	// // consumer.Partitions 用户获取Topic上所有的Partitions. 消息服务器上已经创建了test这个topic,所以,在这里指定参数为test.
	// partitionList, err := consumer.Partitions("test")
	// if err != nil {
	// 	logger.Println("Failed to get the list of partitions: ", err)
	// }

	// for partition := range partitionList {
    //     pc, err := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
    //     if err != nil {
    //         logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
    //     }
    //     defer pc.AsyncClose()

    //     wg.Add(1)

    //     go func(sarama.PartitionConsumer) {
    //         defer wg.Done()
    //         for msg := range pc.Messages() {
    //             fmt.Println("message is :", msg)
    //             fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
    //             fmt.Println()
    //         }
    //     }(pc)
    // }
    //wg.Wait()
}