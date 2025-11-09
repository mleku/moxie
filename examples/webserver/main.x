package main

import (
	"github.com/mleku/moxie/src/fmt"
	"github.com/mleku/moxie/internal/net/http"
	"github.com/mleku/moxie/src/os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Moxie web server!\n")
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
	})

	port := ":8080"
	fmt.Printf("Starting Moxie web server on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
