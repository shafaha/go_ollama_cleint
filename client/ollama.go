package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/llm_project/models"
)

const (
	url = "http://localhost:11434/api/generate"
)

type ollamaClient struct {
}

type Payload struct {
	Model  string `json:"model"`
	Prompt string `json:"Prompt"`
	Stream bool   `json:"stream"`
}

type ResponsePayload struct {
	Output string `json:"output"`
}

func NewOllamaClient() *ollamaClient {
	return &ollamaClient{}
}

func (client *ollamaClient) StreamIt(ctx context.Context, req *models.LLMRequest, ch chan string, errCh chan any) {
	fmt.Println("client:StreamIt request received")
	requestBody := Payload{
		Model:  req.Model,
		Prompt: req.Prompt,
		Stream: true, // Enable streaming from Ollama
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		errCh <- err
		return
	}

	llmRequest, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		errCh <- err
		return
	}

	llmRequest.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(llmRequest)
	if err != nil {
		fmt.Printf("Error making request to Ollama: %v\n", err)
		errCh <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("Error: Received non-OK status code %d from Ollama", resp.StatusCode)
		fmt.Println(errMsg)
		errCh <- errors.New(errMsg)
		return
	}

	// Instead of using scanner, use io.Reader to handle the stream
	fmt.Println("Reading the response now: 77")
	reader := bufio.NewReader(resp.Body)
	for {
		// Read the data into chunks
		bts, err := reader.ReadBytes('\n')
		if err != nil && err.Error() != "EOF" {
			fmt.Printf("Error reading stream: %v\n", err)
			errCh <- err
			break
		}

		// If we have data, push it to the channel
		if len(bts) != 0 {
			llmResponse := &Response{}
			json.Unmarshal(bts, llmResponse)
			fmt.Printf("data received: %s\n", bts)
			ch <- llmResponse.Message
		}

		// Break on EOF (end of stream)
		if err == io.EOF {
			fmt.Printf("err: EndofFile received: %s\n", bts)
			errCh <- errors.New("")
			break
		}
	}

	fmt.Println("Closing the data channel")
}

func (client *ollamaClient) Generate(req *models.LLMRequest) (*models.LLMResponse, error) {
	payload := Payload{
		Model:  req.Model,
		Prompt: req.Prompt,
	}

	fmt.Println("Calling the client now")

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("unable to connect to client")
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("unable to read the response")
	}

	// Parse the response JSON
	var response ResponsePayload
	if err := json.Unmarshal(b, &response); err != nil {
		fmt.Println(err)
		return nil, errors.New("unable to read the response")
	}

	return &models.LLMResponse{
		Response: response.Output,
	}, nil
}
