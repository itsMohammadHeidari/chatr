package config

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Printf("Warning: unable to load .env file: %v", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func BindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = viper.BindPFlag(f.Name, f)
	})

	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = viper.BindPFlag(f.Name, f)
	})

	for _, c := range cmd.Commands() {
		BindFlags(c)
	}
}
