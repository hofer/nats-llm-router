package main

import (
	"context"
	"github.com/hofer/nats-llm/pkq/llm"
	"github.com/nats-io/nats.go"
	"github.com/ollama/ollama/api"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Example to chat with Ollama
	chatWithOllama(nc)

	// Example to chat with Gemini
	chatWithGemini(nc)
}

func chatWithOllama(nc *nats.Conn) {
	ollamaLlm := llm.NewNatsOllamaLLM(nc)

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*5)
	firstMessage := api.Message{
		Role:    "user",
		Content: "How many people are in this image?",
		Images: []api.ImageData{
			readImageData(),
		},
	}
	ollamaRes, err := ollamaLlm.Chat(ctx, &api.ChatRequest{
		Model: "gemma3:27b",
		Messages: []api.Message{
			firstMessage,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Info(ollamaRes.Message.Content)
}

func chatWithGemini(nc *nats.Conn) {
	geminiLlm := llm.NewNatsGeminiLLM(nc)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	firstMessage := api.Message{
		Role:    "user",
		Content: "What is the man in the middle holding in his hands?",
		Images: []api.ImageData{
			readImageData(),
		},
	}
	geminiRes, err := geminiLlm.Chat(ctx, &api.ChatRequest{
		Model: "gemini-2.5-pro-exp-03-25",
		Messages: []api.Message{
			firstMessage,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("First response from LLM: %s", geminiRes.Message.Content)
	secondMessage := api.Message{
		Role:    "user",
		Content: "Is the person on the left holding flowers as well?",
	}
	geminiRes, err = geminiLlm.Chat(ctx, &api.ChatRequest{
		Model: "gemini-2.5-pro-exp-03-25",
		Messages: []api.Message{
			firstMessage,
			geminiRes.Message,
			secondMessage,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Second response from LLM: %s", geminiRes.Message.Content)
}

func readImageData() api.ImageData {
	filePath := "example.png"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", filePath, err)
	}
	// 2. Ensure the file is closed when the function exits
	defer file.Close()

	// 3. Read all bytes from the file
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}
	return data
}
