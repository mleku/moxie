package main

import (
	"github.com/mleku/moxie/src/encoding/json"
	"github.com/mleku/moxie/src/fmt"
	"github.com/mleku/moxie/internal/net/http"
	"github.com/mleku/moxie/src/os"
	"github.com/mleku/moxie/internal/time"
)

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func main() {
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			Message:   "Moxie API is running",
			Timestamp: time.Now(),
			Version:   "1.0.0",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	port := ":8080"
	fmt.Printf("Starting Moxie JSON API on %s\n", port)
	fmt.Println("Try: curl http://localhost:8080/api/status")

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
