package main

import (
	"context"
	"github.com/hofer/nats-llm/pkq/llm"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/ollama/ollama/api"
	"os"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}


	// Example to chat with Ollama
	chatWithOllama(nc)

	// Example to chat with Gemini
	chatWithGemmini(nc)
}
func chatWithOllama(nc *nats.Conn) {
	ollamaLlm := llm.NewNatsOllamaLLM(nc)

	ctx := context.Background()
	ollamaRes, err := ollamaLlm.Chat(ctx, &api.ChatRequest{
		Model: "",
		Messages: []api.Message{
			{
				Role: "user",
				Content: "Who is the current president of the USA?",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info(ollamaRes.Message.Content)
}

func chatWithGemmini(nc *nats.Conn) {
	geminiLlm := llm.NewNatsGeminiLLM(nc)

	ctx := context.Background()
	geminiRes, err := geminiLlm.Chat(ctx, &llm.GeminiChatRequest{
		Model: "",
		Text: "Who is the current president of the USA?",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info(geminiRes.Response)
}
