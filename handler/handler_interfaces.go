package handler

import "net/http"

type LLMHandler interface {
	Generate(http.ResponseWriter, *http.Request)
	StreamIt(w http.ResponseWriter, r *http.Request)
}
