// go_program_caller.go

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func main() {
	// Send a GET request to the API endpoint
	resp, err := http.Get("http://localhost:8080/questioneasy")
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Parse the JSON response into a QuestionAnswer struct
	var qa QuestionAnswer
	if err := json.Unmarshal(body, &qa); err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// Print the question and answer
	log.Println("Question:", qa.Question)
	log.Println("Answer:", qa.Answer)
}
