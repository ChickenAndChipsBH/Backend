package main

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
	questionsEasy       []QuestionAnswer
	questionsMedium     []QuestionAnswer
	questionsHard       []QuestionAnswer
	usedQuestionsEasy   = make(map[int]bool)
	usedQuestionsMedium = make(map[int]bool)
	usedQuestionsHard   = make(map[int]bool)
	mu                  sync.Mutex
)

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func init() {
	err := loadQuestionsFromCSV("Backend/questionseasy.csv", &questionsEasy)
	if err != nil {
		log.Fatal("Error loading easy questions:", err)
	}
	err = loadQuestionsFromCSV("Backend/questionsmedium.csv", &questionsMedium)
	if err != nil {
		log.Fatal("Error loading medium questions:", err)
	}
	err = loadQuestionsFromCSV("Backend/questionshard.csv", &questionsHard)
	if err != nil {
		log.Fatal("Error loading hard questions:", err)
	}
}

func loadQuestionsFromCSV(filePath string, questions *[]QuestionAnswer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|' // Set the delimiter to '|'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		qa := QuestionAnswer{
			Question: record[0],
			Answer:   record[1],
		}
		*questions = append(*questions, qa)
	}

	return nil
}

func HandleQuestionEasy(w http.ResponseWriter, r *http.Request) {
	handleQuestion(w, r, &questionsEasy, &usedQuestionsEasy)
}

func HandleQuestionMedium(w http.ResponseWriter, r *http.Request) {
	handleQuestion(w, r, &questionsMedium, &usedQuestionsMedium)
}

func HandleQuestionHard(w http.ResponseWriter, r *http.Request) {
	handleQuestion(w, r, &questionsHard, &usedQuestionsHard)
}

func handleQuestion(w http.ResponseWriter, r *http.Request, questions *[]QuestionAnswer, usedQuestions *map[int]bool) {
	mu.Lock()
	defer mu.Unlock()

	// Check if all questions have been used
	if len(*usedQuestions) == len(*questions) {
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
		idx := rand.Intn(len(*questions))
		if !(*usedQuestions)[idx] {
			qa = (*questions)[idx]
			(*usedQuestions)[idx] = true
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

func main() {
	log.Printf("Starting quiz server\n")

	log.Printf("Starting server on port 8080...\n")

	// Register the endpoints for different difficulty levels
	http.HandleFunc("/questioneasy", HandleQuestionEasy)
	http.HandleFunc("/questionmedium", HandleQuestionMedium)
	http.HandleFunc("/questionhard", HandleQuestionHard)

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
