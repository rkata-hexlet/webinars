#!/bin/bash

# 1. Создаем исходный код Go (если его еще нет)
cat <<EOF > main.go
package main
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	role := os.Getenv("SERVICE_NAME")
	port := "8080"

	if role == "transcoder" {
		log.Printf("--- [TRANSCODER] Started ---")
		for {
			log.Printf("[TRANSCODER] Processing chunk ID:%d...", time.Now().Unix()%1000)
			time.Sleep(7 * time.Second)
		}
	}

	if role == "catalog" {
		http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"videos": []string{"Interstellar", "Inception", "Tenet"},
				"node":   os.Getenv("HOSTNAME"),
			})
		})
	}

	if role == "api" {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			resp, err := http.Get("http://catalog-service/list")
			if err != nil {
				http.Error(w, "Catalog unreachable", 500)
				return
			}
			defer resp.Body.Close()
			var data map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&data)
			fmt.Fprintf(w, "<h1>StreamFlow Gateway</h1><p>Movies: %v</p><small>Pod: %s</small>", data["videos"], os.Getenv("HOSTNAME"))
		})
	}

	log.Printf("Service [%s] starting on %s", role, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
EOF

# 2. Создаем Dockerfile
cat <<EOF > Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
EOF

echo "--- Сборка образа внутри Minikube ---"

# 3. МАГИЯ: Переключаем контекст Docker на Minikube
# Это позволяет не пушить образ в интернет, а класть его сразу в кластер
eval $(minikube docker-env)

# 4. Сборка
docker build -t streamflow-mock:v1 .

echo "--- Готово! Образ streamflow-mock:v1 доступен для Kubernetes ---"