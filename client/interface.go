package client

import (
	"context"

	"github.com/llm_project/models"
)

type LLMClient interface {
	Generate(*models.LLMRequest) (*models.LLMResponse, error)
	StreamIt(context.Context, *models.LLMRequest, chan string, chan any)
}
