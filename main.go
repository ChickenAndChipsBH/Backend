// main.go
package main

import (
	"log"
	"net/http"

	"github.com/ChickenAndChipsBH/Backend/quiz"
)

func main() {
	log.Printf("Starting quiz server\n")

	log.Printf("Starting server on port 8080...\n")

	// Register the /question endpoint
	http.HandleFunc("/question", quiz.HandleQuestion)

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
