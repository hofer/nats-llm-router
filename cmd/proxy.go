package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var proxyNatsUrl string

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.PersistentFlags().StringVarP(&proxyNatsUrl, "url", "u", os.Getenv("NATS_URL"), "URL to the Nats.io server")
	proxyCmd.MarkFlagRequired("url")
}
