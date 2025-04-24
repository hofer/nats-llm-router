package proxy

func GetGeminiSchemaChat() (string, error) {
	return marshalSchema(&GeminiChatRequest{}, &GeminiChatResponse{})
}