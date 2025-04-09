# Nats to LLM router

This component is a Nats Service listening to messages and forwards requests to a LLM (Ollama, Gemini).

Use the following command to send messages to this service:
```
# Generate Text
nats req --reply-timeout=10s ollama.generate '{"model": "gemma2:2b", "prompt": "What is atorvastatin? Respond in one sentence."}'

# Create an embedding:
nats req --reply-timeout=10s ollama.embed '{"model": "snowflake-arctic-embed2", "input": "What is atorvastatin? Respond in one sentence."}'
```

Limitation: nats does have a size limit for payload.

## Nats cli commands
Given the nats-llm-router is based on Nats Mirco, the following commands are useful:

List services:
```
nats micro ls
```

List service info:
```
nats micro info NatsOllama
```