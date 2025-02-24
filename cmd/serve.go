/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/explore_service/cmd/serve"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "server for dating app",
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.AddCommand(serve.GrpcServerCmd)
}
