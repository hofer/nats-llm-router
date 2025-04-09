package main

import (
	"github.com/nats-io/nats.go"
	"runtime"
)



// main
func main() {
	const natsTailscale = "nats://100.93.123.116:4222"
	nc, _ := nats.Connect(natsTailscale)

	natsOllamaProxy := NewNatsOllamaProxy()
	natsOllamaProxy.Start(nc)

	runtime.Goexit()
}
