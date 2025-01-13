package cmd

import (
	"log"
	"strconv"

	"github.com/itsmohammadheidari/chatr/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Chatr server",
	Long:  `Run the Chatr server, accepting TCP connections and broadcasting messages among connected clients.`,
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

		log.Printf("Starting server on %s:%d...\n", host, port)

		srv := server.NewServer(host, port)
		if err := srv.Start(); err != nil {
			log.Fatalf("Server error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// CLI flags
	serverCmd.Flags().StringP("host", "H", "", "Server host (overrides .env if set)")
	serverCmd.Flags().StringP("port", "P", "", "Server port (overrides .env if set)")
}
