package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var blue = "\033[34m"
var purple = "\033[35m"
var cyan = "\033[36m"
var gray = "\033[37m"
var white = "\033[97m"

func main() {

	var kafkaBroker = os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	var kafkaTopic = os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "my-topic"
	}
	var consumerGroup = os.Getenv("KAFKA_CONSUMER_GROUP")
	if consumerGroup == "" {
		consumerGroup = "my-topic-consumers"
	}
	var sleep = os.Getenv("SLEEP")
	if sleep == "" {
		sleep = "0"
	}
	sleepSeconds, err := strconv.Atoi(sleep)
	if err != nil {
		sleepSeconds = 0
	}
	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		GroupID:  consumerGroup,
		Topic:    kafkaTopic,
		MinBytes: 10e1, // 1KB
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()
	fmt.Printf("%sWaiting for messages on %sbroker: %v %sgroup: %v %stopic: %v%s\n", yellow, blue, r.Config().Brokers, purple, r.Config().GroupID, cyan, r.Config().Topic, reset)
	go func() {
		for {
			fmt.Printf("%sWaiting for a new message. %sCurrent Offset: %d, %sLag: %d %s\n ", yellow, blue, r.Offset(), purple, r.Lag(), reset)
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				fmt.Printf("Error occured while waiting for message %s%s", red, err)
				break
			}
			fmt.Printf("%smessage at %stopic/%spartition/%soffset %s%v/%s%v/%s%v: %s%v = %v%s\n", red, blue, purple, cyan, blue, m.Topic, purple, m.Partition, cyan, m.Offset, green, string(m.Key), string(m.Value), reset)
			time.Sleep(time.Second * time.Duration(sleepSeconds))
			fmt.Printf("%sSlept for %s%d seconds%s\n", yellow, blue, sleepSeconds, reset)
		}
	}()
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	fmt.Printf("%sGot signal: %s%s\n", red, sig, reset)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	} else {
		fmt.Printf("%sClosed the reader%s\n", green, reset)
	}
}
