package main

import (
	"github.com/hofer/nats-llm/pkq/llm"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_SERVER_URL"))
	if err != nil {
		log.Fatal(err)
	}

	llm.NewNatsOllamaLLM(nc)
}
