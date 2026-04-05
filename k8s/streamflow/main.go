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

	// --- ОБЩИЕ ЭНДПОИНТЫ ДЛЯ ВСЕХ РОЛЕЙ ---
	// Теперь Liveness Probe всегда будет получать 200 OK
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	// Имитация поломки для тестов студентами
	http.HandleFunc("/fail", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("CRITICAL ERROR: Simulated service crash!")
	})

	// --- ЛОГИКА ПО РОЛЯМ ---

	if role == "transcoder" {
		log.Printf("--- [TRANSCODER] Started. Processing video... ---")
		// Воркеру не нужен HTTP сервер для работы, но он нужен для Liveness Probe!
		go func() {
			for {
				log.Printf("[TRANSCODER] Processing chunk %d...", time.Now().Unix()%100)
				time.Sleep(5 * time.Second)
			}
		}()
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
			// Обращаемся к сервису по DNS имени внутри K8s
			resp, err := http.Get("http://catalog-service/list")
			if err != nil {
				http.Error(w, "Catalog unreachable: "+err.Error(), 500)
				return
			}
			defer resp.Body.Close()

			var data map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&data)

			fmt.Fprintf(w, "<h1>StreamFlow Gateway</h1>")
			fmt.Fprintf(w, "<p>Movies: %v</p>", data["videos"])
			fmt.Fprintf(w, "<hr><small>API Pod: %s | Catalog Pod: %s</small>", os.Getenv("HOSTNAME"), data["node"])
		})
	}

	log.Printf("Service [%s] started on port %s", role, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
