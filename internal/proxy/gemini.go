package proxy

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func gemini() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
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

	model := client.GenerativeModel("gemini-1.5-pro-latest")

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

	res, err := session.SendMessage(ctx, genai.Text("Who is the current president of the USA?"))
	//res, err := session.SendMessage(ctx, genai.Text("Which theaters in Mountain View show Barbie movie?"))
	if err != nil {
		log.Fatalf("session.SendMessage: %v", err)
	}

	printResponse(res)

	//
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
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
