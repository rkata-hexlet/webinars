package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func consumer_worker_pool() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "order-events",
		GroupID:        "payments-service",
		CommitInterval: 0, // ручной коммит
	})
	defer r.Close()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer cancel()

	const numWorkers = 8
	msgs := make(chan kafka.Message, numWorkers) // буфер = backpressure

	// ─────────────────────────────────────────
	// Запускаем воркеры
	// ─────────────────────────────────────────
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for msg := range msgs {
				// Обрабатываем
				if err := process1(msg.Value); err != nil {
					log.Printf("worker %d: process error offset=%d: %v", id, msg.Offset, err)
					// Не коммитим — сообщение придёт снова после рестарта
					continue
				}

				// Коммитим только после успешной обработки
				if err := r.CommitMessages(ctx, msg); err != nil {
					log.Printf("worker %d: commit error offset=%d: %v", id, msg.Offset, err)
				}
			}
		}(i)
	}

	// ─────────────────────────────────────────
	// Читаем из Kafka и кладём в канал
	// ─────────────────────────────────────────
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			break // ctx отменён → выходим
		}
		msgs <- msg // блокирует если все воркеры заняты → backpressure
	}

	close(msgs) // сигнал воркерам — больше сообщений не будет
	wg.Wait()   // ждём пока все воркеры закончат текущие сообщения

	log.Println("consumer stopped gracefully")
}

func process1(value []byte) error {
	fmt.Printf("processing: %s\n", value)
	return nil
}
