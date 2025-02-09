package models

type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type LLMResponse struct {
	Response string `json:"response"`
}
