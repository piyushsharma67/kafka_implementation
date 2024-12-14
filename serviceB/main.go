package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9093",
		"group.id":          "serviceB",
		"auto.offset.reset": "smallest",
		"debug":             "generic,broker,consumer",
	})
	if err != nil {
		fmt.Printf("closing service_b %s", err)
		os.Exit(1)
	}
	// fmt.Println("service_b")
	// err = consumer.SubscribeTopics([]string{"serviceA"}, nil)
	// run := true
	// for run {
	// 	ev := consumer.Poll(100)
	// 	switch e := ev.(type) {
	// 	case *kafka.Message:
	// 		// application-specific processing
	// 		fmt.Printf("Received message: %s\n", string(e.Value))
	// 		run = false
	// 	case kafka.Error:
	// 		fmt.Fprintf(os.Stderr, "%% Error service_b: %v\n", e)
	// 		// run = false
	// 	default:
	// 		fmt.Printf("Ignored %v\n", e)
	// 	}
	// }

	messageChannel := make(chan *kafka.Message)
	err = consumer.SubscribeTopics([]string{"serviceA"}, nil)
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
	}

	close(messageChannel)
	wg.Wait()
	consumer.Close()
}
