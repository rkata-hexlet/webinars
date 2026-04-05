package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func consumer_basic() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "order-events",
		GroupID: "payments-service",
	})
	defer r.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	for {
		msg, err := r.ReadMessage(ctx) // automatically commit the offset, returns io.EOF on close reader
		if err != nil {
			break // ctx отменён → выходим из loop
		}
		log.Printf("offset=%d key=%s value=%s",
			msg.Offset, msg.Key, msg.Value)
		fmt.Printf("processing: %s\n", msg.Value)
	}

	log.Println("consumer stopped gracefully")
}
