package proxy

import (
	"github.com/ollama/ollama/api"
)

func GetGeminiSchemaChat() (string, error) {
	return marshalSchema(&api.ChatRequest{}, &api.ChatResponse{})
}
