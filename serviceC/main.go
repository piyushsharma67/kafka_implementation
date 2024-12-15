package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9093",
		"client.id":         "serviceC",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("closing %s", err)
		os.Exit(1)
	}

	topic := "message"

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(fmt.Sprintf("yay i am called"))},
		nil, // delivery channel
	)

	if err != nil {
		fmt.Printf("closing %s", err)
		os.Exit(0)
	}

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received message: %s\n", msg)

		// Echo the message back to the client
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}

	defer conn.Close()
}

func main() {
	http.HandleFunc("/", Handler)

	err := http.ListenAndServe(fmt.Sprintf(":%d", 3000), nil)

	if err != nil {
		os.Exit(1)
	}
}
