package cmd

import (
	"fmt"
	"os"

	"github.com/itsmohammadheidari/chatr/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Chatr",
	Short: "Chatr is a TCP-based TUI chat system.",
	Long:  "Chatr is a TCP-based TUI chat system that establishes client-server connections using raw TCP sockets.",
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := config.LoadConfig(); err != nil {
			return err
		}
		config.BindFlags(rootCmd)
		return nil
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
