package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	kafka "github.com/segmentio/kafka-go"
)

func main() {
	e := echo.New()
	var bindAddress = os.Getenv("BIND_ADDRESS")
	if bindAddress == "" {
		bindAddress = ":8090"
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", sendMessages)

	// Start server
	e.Logger.Fatal(e.Start(bindAddress))
	// trap sigterm or interupt and gracefully shutdown the server
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	// Block until a signal is received.
	sig := <-ch
	fmt.Printf("Got signal: %s", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	e.Shutdown(ctx)

}

func sendMessages(c echo.Context) error {
	var kafkaBroker = os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}
	var kafkaTopic = os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "my-topic"
	}
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	var wg sync.WaitGroup
	var messageCount int
	messageCount, err := strconv.Atoi(c.QueryParam("count"))
	if err != nil {
		messageCount = 100
	}
	for i := 0; i < messageCount; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			writeMessage(w, i)
		}()
	}

	wg.Wait()

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
		return err
	}
	return nil
}

func writeMessage(w *kafka.Writer, i int) {
	message := kafka.Message{Key: []byte(fmt.Sprintf("Key-%d", i)), Value: []byte(fmt.Sprintf("Message No: %d, Time: %s", i, time.Now()))}
	err := w.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Printf("failed to write messages: %s\n", err)
	} else {
		fmt.Printf("message delivered %d\n", i)
	}
}
