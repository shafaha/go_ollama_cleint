package client

import "time"

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"response"`
	Done      bool      `json:"done"`
}
