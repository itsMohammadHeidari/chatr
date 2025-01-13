package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/itsmohammadheidari/chatr/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run the Chatr client (TUI-based)",
	Long:  "Run the Chatr client, connecting to a TCP server.",

	Run: func(cmd *cobra.Command, args []string) {
		host := viper.GetString("host")
		if host == "" {
			host = "127.0.0.1"
		}

		portStr := viper.GetString("port")
		if portStr == "" {
			portStr = "8080"
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port: %s\n", portStr)
		}

		username := viper.GetString("username")
		if username == "" {
			username = "Guest"
		}

		log.Printf("Connecting to server at %s:%d as %s...\n", host, port, username)

		cl := client.NewClient(host, port, username)
		if err := cl.Start(); err != nil {
			fmt.Printf("Client error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// CLI flags
	clientCmd.Flags().StringP("host", "H", "", "Server host (overrides .env if set)")
	clientCmd.Flags().StringP("port", "P", "", "Server port (overrides .env if set)")
	clientCmd.Flags().StringP("username", "u", "", "Username to display in the chat (overrides .env if set)")
}
