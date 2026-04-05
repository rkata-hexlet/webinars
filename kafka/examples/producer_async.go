package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

func producer_async() {
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "order-events",
		Balancer:               &kafka.LeastBytes{},
		Async:                  true,
		BatchSize:              100,
		BatchTimeout:           10_000_000, // 10ms в наносекундах
		AllowAutoTopicCreation: true,
	}
	defer w.Close()

	messages := []kafka.Message{
		{Key: []byte("u1"), Value: []byte(`{"event":"order"}`)},
		{Key: []byte("u2"), Value: []byte(`{"event":"payment"}`)},
	}

	var wg sync.WaitGroup
	for _, msg := range messages {
		wg.Add(1)
		go func(m kafka.Message) {
			defer wg.Done()
			if err := w.WriteMessages(context.Background(), m); err != nil {
				fmt.Printf("write error: %v\n", err)
			}
		}(msg)
	}

	wg.Wait()
	w.Close()
}
