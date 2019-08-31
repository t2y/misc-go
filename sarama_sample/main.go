package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

const (
	listTopicsCommand  = "listTopics"
	createTopicCommand = "createTopic"
	sendMessageCommand = "send"
	recvMessageCommand = "recv"
)

var (
	bootstrap string
	command   string
	topic     string
	message   string
)

func init() {
	flag.StringVar(
		&bootstrap, "bootstrap", "127.0.0.1:9092",
		"comma-separated bootstrap servers")
	flag.StringVar(&command, "command", listTopicsCommand, "command")
	flag.StringVar(&topic, "topic", "mytopic", "topic")
	flag.StringVar(&message, "message", "test", "message")
}

func listTopics(brokers []string, config *sarama.Config) {
	consumer := getConsumer(brokers, config)
	defer func() { _ = consumer.Close() }()
	topics, err := consumer.Topics()
	if err != nil {
		log.Fatal(err)
	}
	for i := range topics {
		fmt.Println(topics[i])
	}
}

func createTopic(config *sarama.Config, brokers []string, topic string) {
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	defer func() { _ = admin.Close() }()

	detail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = admin.CreateTopic(topic, &detail, false)
	if err != nil {
		log.Fatal("Error while creating topic: ", err.Error())
	}

	log.Printf("topic: %s is created", topic)
}

func sendMessage(config *sarama.Config, brokers []string, topic, message string) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() { _ = producer.Close() }()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println("send error: ", err.Error())
		return
	}

	fmt.Println("Partition: ", partition)
	fmt.Println("Offset: ", offset)
}

func recvMessage(config *sarama.Config, brokers []string, topic string) {
	consumer := getConsumer(brokers, config)
	defer func() { _ = consumer.Close() }()

	partitionConsumer, err := consumer.ConsumePartition(
		topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Println(err)
		return
	}
	defer partitionConsumer.Close()

	for {
		msg := <-partitionConsumer.Messages()
		log.Printf("Consumed message: [%s], offset: [%d]\n", msg.Value, msg.Offset)
	}
}

func getConsumer(brokers []string, config *sarama.Config) sarama.Consumer {
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}

func main() {
	flag.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_2_0_0

	brokers := strings.Split(bootstrap, ",")
	log.Printf("brokers: %v", brokers)

	log.Printf("command: %s", command)
	if command == listTopicsCommand {
		listTopics(brokers, config)
	} else if command == createTopicCommand {
		createTopic(config, brokers, topic)
	} else if command == sendMessageCommand {
		sendMessage(config, brokers, topic, message)
	} else if command == recvMessageCommand {
		recvMessage(config, brokers, topic)
	} else {
		log.Fatalf("unsupported: %s", command)
	}
}
