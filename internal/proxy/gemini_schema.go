package proxy

import "github.com/hofer/nats-llm/pkq/llm"

func GetGeminiSchemaChat() (string, error) {
	return marshalSchema(&llm.GeminiChatRequest{}, &llm.GeminiChatResponse{})
}