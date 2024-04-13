// quiz/quiz.go
package quiz

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

var (
	questions     []QuestionAnswer
	usedQuestions = make(map[int]bool)
	mu            sync.Mutex
)

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func init() {
	// Load questions from CSV file
	err := loadQuestionsFromCSV("questions.csv")
	if err != nil {
		log.Fatal("Error loading questions:", err)
	}
}

func loadQuestionsFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		qa := QuestionAnswer{
			Question: record[0],
			Answer:   record[1],
		}
		questions = append(questions, qa)
	}

	return nil
}

func HandleQuestion(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Check if all questions have been used
	if len(usedQuestions) == len(questions) {
		// Respond with JSON error message
		errorMessage := struct {
			Error string `json:"error"`
		}{"No more questions available"}
		jsonResponse, _ := json.Marshal(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
		return
	}

	// Get a random unused question
	var qa QuestionAnswer
	for {
		idx := rand.Intn(len(questions))
		if !usedQuestions[idx] {
			qa = questions[idx]
			usedQuestions[idx] = true
			break
		}
	}

	// Marshal the QuestionAnswer struct to JSON
	jsonResponse, err := json.Marshal(qa)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error marshaling JSON:", err)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	w.Write(jsonResponse)
}
