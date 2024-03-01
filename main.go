package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	dotenv "github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)


// Schemas

// Schemas for /open_api_completion
type Request struct {
	Question string `json:"question"`
}

type Response struct {
    ApiCompletion string `json:"api_completion"`
}

// Schemas for /add
type Add struct {
	A string `json:"a"`
	B string `json:"b"`
}

type Res struct {
	Total int `json:"total"`
}

//Schemas for /total
type Freq struct {
	String string `json:"string"`
}

type Total map[string]int


func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// My Endpoint here
	http.HandleFunc("/", handleRootDir)
	http.HandleFunc("/open_api_completion", handleOpenApiEndPoint)
	http.HandleFunc("/add", handleAdd)
	http.HandleFunc("/total",hadleTotal)

	// My web host location
	err2 := http.ListenAndServe(":8080", nil)

	// Error handling
	if err2 != nil {
		log.Fatal(err2)
	}
}

func handleRootDir(w http.ResponseWriter, r *http.Request) {
	
	response := map[string]string{
									"message": "Hello World",
								}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleOpenApiEndPoint(w http.ResponseWriter, r *http.Request) {
	log.Println("open-api-end point os being called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowd",http.StatusForbidden)
	}

	var reqBody Request

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg := reqBody.Question
	apiKey := os.Getenv("openapi")
	response := Response{
		ApiCompletion: openApiMessage(apiKey, msg),
	}

	jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	log.Println("Trying to call the add endpoint")

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong Method", 404)
	}

	var reqBody Add

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := Res {
		Total: addTwoString(reqBody.A, reqBody.B),
	}

	jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func hadleTotal(w http.ResponseWriter, r *http.Request) {
	if (r.Method != http.MethodPost) {
		http.Error(w, "Wrong Method", 404)
	}

	var reqBody Freq
	
	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	x := reqBody.String
	mp := make(map[rune]int)
	for _,v := range x {
		mp[v] += 1
	}

	response := make(Total)

	for key, value := range mp {
		response[string(key)] = value
	}

	jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}


func addTwoString(a, b string) int {
	first, err1 := strconv.Atoi(a)
	second, err2 := strconv.Atoi(b)

	if err1 != nil && err2 != nil {
		panic("Not a valid number")
	}
	
	result := first + second


	return result
}

func openApiMessage(secretKey, content string) string {
	client := openai.NewClient(secretKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "Fix the error"
	}

	fmt.Println(resp.Choices[0].Message.Content)

	return resp.Choices[0].Message.Content
}
