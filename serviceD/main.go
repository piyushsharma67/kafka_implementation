package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

func main() {
	//creating a connection to service c using websocket
	conn, _, err := websocket.DefaultDialer.Dial("ws://service_c:3000/ws", nil)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	fmt.Println("Connected to the server c")
	defer conn.Close()

	//creating consomer config to consume message from Kafka
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9093",
		"group.id":          "serviceD",
		"auto.offset.reset": "smallest",
		"debug":             "generic,broker,consumer",
	})
	if err != nil {
		fmt.Printf("closing service_b %s", err)
		os.Exit(1)
	}

	messageChannel := make(chan *kafka.Message)
	err = consumer.SubscribeTopics([]string{"message"}, nil)
	var wg sync.WaitGroup
	wg.Add(1)

	if err != nil {
		fmt.Printf("Exiting %v\n", err)
		os.Exit(1)
	}
	go func() {
		defer wg.Done()
		for {
			ev := consumer.Poll(100) // Poll for messages every 100 ms
			switch e := ev.(type) {
			case *kafka.Message:
				// Send the received message to the channel
				messageChannel <- e
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error service_b: %v\n", e)
				// You can handle errors here, like logging or retries
			default:
			}
		}
	}()

	for msg := range messageChannel {
		// Process the received message
		fmt.Printf("Received message: %s\n", string(msg.Value))
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello, Serverc! i have got message %s", msg.Value)))
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
	}

	close(messageChannel)
	wg.Wait()
	consumer.Close()
}
