package llm

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"time"
)

func NewNatsGeminiLLM(nc *nats.Conn) *NatsGeminiLLM {
	return &NatsGeminiLLM{
		client: nc,
	}
}

type NatsGeminiLLM struct {
	client *nats.Conn
}

func (n *NatsGeminiLLM) Chat(ctx context.Context, req *GeminiChatRequest) (GeminiChatResponse, error) {
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return GeminiChatResponse{}, err
	}

	deadline, _ := ctx.Deadline()
	remainingDuration := time.Until(deadline)

	msg, err := n.client.Request("gemini.chat", jsonStr, remainingDuration)
	if err != nil {
		return GeminiChatResponse{}, err
	}

	var chatResponse GeminiChatResponse
	err = json.Unmarshal(msg.Data, &chatResponse)
	if err != nil {
		return GeminiChatResponse{}, err
	}

	return chatResponse, nil
}

type GeminiChatRequest struct {
	Model   string `json:"model"`
	Text    string `json:"text"`
}

type GeminiChatResponse struct {
	Response string `json:"response"`
}
