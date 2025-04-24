package proxy

import (
	"context"
	"fmt"
	"github.com/hofer/nats-llm/pkq/llm"
	log "github.com/sirupsen/logrus"
	"os"
	"encoding/json"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"runtime"
)

func StartNatsGeminiProxy(natsUrl string, ollamaUrl string) error {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}

	natsGeminiProxy := NewNatsGeminiProxy()
	natsGeminiProxy.Start(nc)

	runtime.Goexit()
	return nil
}

type NatsGeminiProxy struct {
	client *genai.Client
}

func NewNatsGeminiProxy() *NatsGeminiProxy {
	return &NatsGeminiProxy{}
}

func (n *NatsGeminiProxy) Start(nc *nats.Conn) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Close()
	n.client = client

	srv, err := micro.AddService(nc, micro.Config{
		Name:        "NatsGemini",
		Version:     "0.0.1",
		Description: "Nats microservice acting as a proxy for Gemini.",
	})
	if err != nil {
		log.Fatal(err)
	}
	//defer srv.Stop()

	root := srv.AddGroup("gemini")


	// Chat
	chatSchema, err := GetGeminiSchemaChat()
	if err != nil {
		log.Fatal(err)
	}
	err = root.AddEndpoint("chat", micro.HandlerFunc(n.chatHandler), micro.WithEndpointMetadata(map[string]string{
		"schema": chatSchema,
	}))
	if err != nil {
		log.Fatal(err)
	}
}

func (n *NatsGeminiProxy) chatHandler(req micro.Request) {

	var reqData llm.GeminiChatRequest
	err := json.Unmarshal(req.Data(), &reqData)
	if err != nil {
		req.Error("400", err.Error(), nil)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Error(err)
		req.Error("500", err.Error(), nil)
		return
	}
	defer client.Close()
	//
	//	// To use functions / tools, we have to first define a schema that describes
	//	// the function to the model. The schema is similar to OpenAPI 3.0.
	//	schema := &genai.Schema{
	//		Type: genai.TypeObject,
	//		Properties: map[string]*genai.Schema{
	//			"location": {
	//				Type:        genai.TypeString,
	//				Description: "The city and state, e.g. San Francisco, CA or a zip code e.g. 95616",
	//			},
	//			"title": {
	//				Type:        genai.TypeString,
	//				Description: "Any movie title",
	//			},
	//		},
	//		Required: []string{"location"},
	//	}
	//
	//	movieTool := &genai.Tool{
	//		FunctionDeclarations: []*genai.FunctionDeclaration{{
	//			Name:        "find_theaters",
	//			Description: "find theaters based on location and optionally movie title which is currently playing in theaters",
	//			Parameters:  schema,
	//		}},
	//	}

	model := client.GenerativeModel(reqData.Model) // "gemini-1.5-pro-latest"

	//	// Before initiating a conversation, we tell the model which tools it has
	//	// at its disposal.
	//	model.Tools = []*genai.Tool{movieTool}

	// For using tools, the chat mode is useful because it provides the required
	// chat context. A model needs to have tools supplied to it in the chat
	// history so it can use them in subsequent conversations.
	//
	// The flow of message expected here is:
	//
	// 1. We send a question to the model
	// 2. The model recognizes that it needs to use a tool to answer the question,
	//    an returns a FunctionCall response asking to use the tool.
	// 3. We send a FunctionResponse message, simulating the return value of
	//    the tool for the model's query.
	// 4. The model provides its text answer in response to this message.
	session := model.StartChat()
	//session.History = reqData.History

	res, err := session.SendMessage(ctx, genai.Text(reqData.Text))
	if err != nil {
		log.Error("session.SendMessage: %v", err)
		req.Error("500", err.Error(), nil)
		return
	}

	//	part := res.Candidates[0].Content.Parts[0]
	//	funcall, ok := part.(genai.FunctionCall)
	//	if !ok || funcall.Name != "find_theaters" {
	//		log.Fatalf("expected FunctionCall to find_theaters: %v", part)
	//	}
	//
	//	// Expect the model to pass a proper string "location" argument to the tool.
	//	if _, ok := funcall.Args["location"].(string); !ok {
	//		log.Fatalf("expected string: %v", funcall.Args["location"])
	//	}
	//
	//	// Provide the model with a hard-coded reply.
	//	res, err = session.SendMessage(ctx, genai.FunctionResponse{
	//		Name: movieTool.FunctionDeclarations[0].Name,
	//		Response: map[string]any{
	//			"theater": "AMC16",
	//		},
	//	})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	printResponse(res)

	resData := llm.GeminiChatResponse{
		Response: getResponseText(res),
		//History = session.History
	}

	responseData, err := json.Marshal(resData)
	if err != nil {
		log.Error("Error marshalling response:", err)
		req.Error("500", err.Error(), nil)
		return
	}
	err = req.Respond(responseData)
}

func getResponseText(resp *genai.GenerateContentResponse) string{
	response := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response += fmt.Sprint(part) //fmt.Println(part)
			}
		}
	}
	return response
}