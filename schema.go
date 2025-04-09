package main

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
	"github.com/ollama/ollama/api"
)

type OllamaProxySchema struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

func schema(request any, response any) (*OllamaProxySchema, error) {
	reflector := jsonschema.Reflector{DoNotReference: true}
	reqSchema, err := reflector.Reflect(request).MarshalJSON()
	if err != nil {
		return nil, err
	}

	resSchema, err := reflector.Reflect(response).MarshalJSON()
	if err != nil {
		return nil, err
	}
	return &OllamaProxySchema{
		Request:  string(reqSchema),
		Response: string(resSchema),
	}, nil
}

func marshalSchema(request any, response any) (string, error) {
	serviceSchema, _ := schema(request, response)
	schemaData, err := json.Marshal(serviceSchema)
	if err != nil {
		return "", err
	}
	return string(schemaData), nil
}

func GetSchemaGenerate() (string, error) {
	return marshalSchema(&api.GenerateRequest{}, &api.GenerateResponse{})
}

func GetSchemaEmbed() (string, error) {
	return marshalSchema(&api.EmbedRequest{}, &api.EmbedResponse{})
}

func GetSchemaEmbedding() (string, error) {
	return marshalSchema(&api.EmbeddingRequest{}, &api.EmbeddingResponse{})
}

func GetSchemaChat() (string, error) {
	return marshalSchema(&api.ChatRequest{}, &api.ChatResponse{})
}
