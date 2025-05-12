# 🏁 NATS to LLM proxy/router
![workflow](https://github.com/hofer/nats-llm/actions/workflows/build.yml/badge.svg)

This cli tool makes Ollama or Gemini LLMs accessible via NATS microservices. Requests sent to this NATS microservice are
forwarded to the corresponding LLM.

> [!WARNING]
> 🚨 🚧 This tool is under active development 🚧 🚨
>
> This tool is very much work in progress. Expect almost daily breaking changes...


Run the following command to start an Ollama proxy:
```bash
./nats-llm proxy ollama --url="nats://localhost:4222"
```

Please check the [the examples folder](./examples) to see how a client can access an LLM exposed via NATS.

## Testing

Use the following command to manually send messages to this service:
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