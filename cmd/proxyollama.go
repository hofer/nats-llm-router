package cmd

import (
	"github.com/hofer/nats-llm/internal/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var proxyOllamaUrl string

// proxyollamaCmd represents the proxyollama command
var proxyollamaCmd = &cobra.Command{
	Use:   "ollama",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Connecting to the Nats.io Server: %s", proxyNatsUrl)
		log.Infof("Connecting to Ollama on url: %s", proxyOllamaUrl)
		err := proxy.StartOllamaProxy(proxyNatsUrl, proxyOllamaUrl)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	proxyCmd.AddCommand(proxyollamaCmd)
	proxyollamaCmd.PersistentFlags().StringVarP(&proxyOllamaUrl, "ollamaUrl", "o", "http://localhost:11434", "URL to the Nats.io server")
}
