package main

import (
	"github.com/spf13/cobra"

	"go-proxy-server/internal/config"
	"go-proxy-server/internal/server"
	"go-proxy-server/pkg/logger"
)

var rootCmd = &cobra.Command{
	Use:   "proxy-server",
	Short: "Proxy server with load balancing, caching, authentication and metrics",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewViperLoader().LoadConfig()
		log := logger.NewLogger(cfg.LogLevel)
		server.StartGRPCServer(cfg, log)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.NewLogger("error").Fatalf("Error executing command: %v", err)
	}
}
