package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/llm_project/client"
	"github.com/llm_project/models"
)

type llmService struct {
	clientMap map[string]client.LLMClient
}

func NewLLMService(name string, c client.LLMClient) *llmService {
	s := &llmService{
		clientMap: make(map[string]client.LLMClient),
	}
	s.clientMap[name] = c
	return s
}

func (svc *llmService) Generate(req *models.LLMRequest) (r *models.LLMResponse, e error) {
	if c, ok := svc.clientMap[req.Model]; c == nil || !ok {
		return nil, errors.New("client not configured yet")
	} else {
		r, e = c.Generate(req)
		fmt.Println(e)
		if e != nil {
			return nil, errors.New("unexpected turn of event")
		}
	}

	return r, e
}

func (svc *llmService) StreamIt(ctx context.Context, req *models.LLMRequest, ch chan string, ech chan any) {
	fmt.Println("service:StreamIt request received")
	errCh := make(chan any)
	respCh := make(chan string)
	defer close(errCh)
	defer close(respCh)
	if c, ok := svc.clientMap[req.Model]; c == nil || !ok {
		ech <- errors.New("client not configured yet")
		return
	} else {
		go c.StreamIt(ctx, req, respCh, errCh)
		for {
			select {
			case err := <-errCh:
				fmt.Printf("service error recieved: %v\n", err)
				ech <- err
				return
			case str := <-respCh:
				fmt.Println("service data recieved: " + str)
				ch <- str
			}
		}
	}
}
