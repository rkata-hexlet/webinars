package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

// ─────────────────────────────────────────────
// 1. AUTO COMMIT — CommitInterval > 0
// ─────────────────────────────────────────────

func runAutoCommit() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "order-events",
		GroupID:        "payments-service",
		CommitInterval: time.Second, // коммит каждую секунду автоматически
	})
	defer r.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for {
		// ReadMessage = FetchMessage + автоматический CommitMessages
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		// ⚠️ Если упадём здесь — offset уже закоммичен, сообщение потеряно
		process(msg.Value)
	}
}

// ─────────────────────────────────────────────
// 2. MANUAL COMMIT — FetchMessage + CommitMessages
// ─────────────────────────────────────────────

func runManualCommit() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "order-events",
		GroupID:        "payments-service",
		CommitInterval: 0, // ручной режим — автокоммит отключён
	})
	defer r.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for {
		// Шаг 1: получаем сообщение (offset НЕ коммитится)
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}

		// Шаг 2: обрабатываем
		if err := process(msg.Value); err != nil {
			// Ошибка обработки — НЕ коммитим, сообщение придёт снова
			log.Printf("process error offset=%d: %v", msg.Offset, err)
			continue
		}

		// Шаг 3: только после успешной обработки коммитим offset
		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("commit error offset=%d: %v", msg.Offset, err)
		}
	}
}

// ─────────────────────────────────────────────
// 3. ПАНИКА ДО КОММИТА — poison pill сценарий
// ─────────────────────────────────────────────

func runPanicDemo() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "order-events",
		GroupID:        "payments-service",
		CommitInterval: 0,
	})
	defer r.Close()

	ctx := context.Background()

	msg, err := r.FetchMessage(ctx)
	if err != nil {
		return
	}

	// 💥 Паника здесь — offset не закоммичен
	process(msg.Value)

	// Эта строка не выполнится при панике
	r.CommitMessages(ctx, msg)

	// → После рестарта consumer прочитает то же сообщение снова
	// → Это и есть At least once семантика
}

// ─────────────────────────────────────────────
// 4. БАТЧ КОММИТ — несколько сообщений за раз
// ─────────────────────────────────────────────

func runBatchCommit() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          "order-events",
		GroupID:        "payments-service",
		CommitInterval: 0,
	})
	defer r.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	batch := make([]kafka.Message, 0, 10)

	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			break
		}

		if err := process(msg.Value); err != nil {
			log.Printf("process error: %v", err)
			continue
		}

		batch = append(batch, msg)

		// Коммитим каждые 10 сообщений
		if len(batch) >= 10 {
			if err := r.CommitMessages(ctx, batch...); err != nil {
				log.Printf("batch commit error: %v", err)
			}
			batch = batch[:0]
		}
	}

	// Коммитим остаток при выходе
	if len(batch) > 0 {
		r.CommitMessages(context.Background(), batch...)
	}
}

// ─────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────

func process(value []byte) error {
	fmt.Printf("processing: %s\n", value)
	// симуляция работы
	time.Sleep(10 * time.Millisecond)
	return nil
}

func consumer_commit() {
	mode := "manual"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "auto":
		log.Println("starting auto commit consumer")
		runAutoCommit()
	case "manual":
		log.Println("starting manual commit consumer")
		runManualCommit()
	case "batch":
		log.Println("starting batch commit consumer")
		runBatchCommit()
	default:
		fmt.Println("usage: go run main.go [auto|manual|batch]")
	}
}
