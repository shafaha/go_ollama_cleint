package main

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/llm_project/client"
	handler "github.com/llm_project/handler/llm"
	service "github.com/llm_project/service/llm"
)

func main() {
	fmt.Println("Hello world")
	router := mux.NewRouter()
	llmClient := client.NewOllamaClient()
	svc := service.NewLLMService("llama3.2", llmClient)
	h := handler.NewLLMHandler(svc)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, or specify the front-end URL like http://localhost:3000
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Allows credentials like cookies, authorization headers, etc.
	})
	router.HandleFunc("/fruits", h.Generate).Methods("POST")
	router.HandleFunc("/stream", h.StreamIt).Methods("POST")

	http.ListenAndServe(":8080", corsHandler.Handler(router))
}
