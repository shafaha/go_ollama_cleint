package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/llm_project/models"
	"github.com/llm_project/service"
)

type llmHandler struct {
	svc service.LLMService
}

func NewLLMHandler(svc service.LLMService) *llmHandler {
	return &llmHandler{
		svc: svc,
	}
}

func (handler *llmHandler) Generate(w http.ResponseWriter, r *http.Request) {
	llmReq := &models.LLMRequest{}
	err := json.NewDecoder(r.Body).Decode(llmReq)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("{\"error\": \"Bad request\"}"))

		return
	}

	resp, err := handler.svc.Generate(llmReq)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error())))

		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", "Internal Error")))

		return
	}
	w.WriteHeader(200)
	w.Write(b)
}

func (handler *llmHandler) StreamIt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler:StreamIt request received")
	llmReq := &models.LLMRequest{}
	err := json.NewDecoder(r.Body).Decode(llmReq)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("{\"error\": \"Bad request\"}"))

		return
	}

	if llmReq.Prompt == "" {
		w.WriteHeader(400)
		w.Write([]byte("{\"error\": \"Prompt is required\"}"))

		return
	}

	ch := make(chan string)
	ech := make(chan any)
	ctx := context.WithoutCancel(r.Context())
	defer close(ch)
	defer close(ech)

	go handler.svc.StreamIt(ctx, llmReq, ch, ech)
	for {
		select {
		case err := <-ech:
			fmt.Printf("handler error recieved: %v\n", err)
			fmt.Fprintf(w, "%s", err)
			w.(http.Flusher).Flush()
			return
		case str := <-ch:
			fmt.Println("handler data recieved: " + str)
			fmt.Fprintf(w, "%s", str)
			w.(http.Flusher).Flush()
		}
	}
}
