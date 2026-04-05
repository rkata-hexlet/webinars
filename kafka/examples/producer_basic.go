package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func producer_basic() {
	w := &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  "order-events",
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
	}
	defer w.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userID := "user-42"
	payload := `{"event":"order_created","amount":1500}`

	err := w.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(userID),
			Value: []byte(payload),
		},
	)
	if err != nil {
		fmt.Printf("write error: %v\n", err)
		return
	}
	fmt.Println("message sent")
}
