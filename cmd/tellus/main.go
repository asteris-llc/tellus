package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/tellus/storage"
	"github.com/asteris-llc/tellus/tf"
	"github.com/asteris-llc/tellus/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// config paths
	viper.SetConfigName("tellus")
	viper.AddConfigPath("/etc/tellus/")
	viper.AddConfigPath("$HOME/.tellus/")
	viper.AddConfigPath(".tellus/")
	viper.ReadInConfig()

	// environment variable defaults
	viper.SetEnvPrefix("tellus")

	// config defaults and args
	CmdServe.Flags().String("address", "", "address to run server on")
	viper.BindEnv("address")

	CmdServe.Flags().String("port", "4000", "port to run server on")
	viper.BindEnv("port", "PORT") // compatible with 12-factor app approach

	CmdServe.Flags().String("storage", "memory", "how to store state given to the server")
	viper.BindEnv("storage")

	// vault storage
	CmdServe.Flags().String("vault-addr", "https://127.0.0.1:8200", "vault address")
	viper.BindEnv("vault-addr", "VAULT_ADDR") // compatible with vault's CLI

	CmdServe.Flags().String("vault-token", "", "vault token")
	viper.BindEnv("vault-token", "VAULT_TOKEN") // compatible with vault's CLI

	CmdServe.Flags().String("vault-mount", "tellus", "which mount to store secrets in")
	viper.BindEnv("vault-mount")

	// grab the flags we've just defined above
	viper.BindPFlags(CmdServe.Flags())

	// logging
	CmdTellusRoot.PersistentFlags().String("log-level", "info", "verbosity level for logs")
	viper.BindPFlag("log-level", CmdTellusRoot.PersistentFlags().Lookup("log-level"))
	viper.BindEnv("log-level")

	CmdTellusRoot.PersistentFlags().String("log-format", "text", "format of logged output")
	viper.BindPFlag("log-format", CmdTellusRoot.PersistentFlags().Lookup("log-format"))
	viper.BindEnv("log-format")
}

func main() {
	defer Recovery()

	// set log levels, etc
	err := setUpLogging()
	if err != nil {
		GracefullyFail(err.Error())
	}

	// run viper
	CmdTellusRoot.AddCommand(CmdServe, CmdVersion)
	CmdTellusRoot.Execute()
}

var (
	Version = "0.0.1"

	CmdTellusRoot = &cobra.Command{
		Use:   "tellus",
		Short: "Tellus is a collaboration toolkit for teams using Terraform",
		Long:  "Tellus allows sharing and saving .tfstate files and other bits of Terraform config for use in a team. Full documentation available at https://github.com/asteris-llc/tellus",
	}

	CmdServe = &cobra.Command{
		Use:   "serve",
		Short: "start and run the server",
		Long:  "start and run the Tellus server on a given port. This command will continue until interrupted.",
		Run: func(cmd *cobra.Command, args []string) {
			address := viper.GetString("address")
			if address == "" {
				address = "0.0.0.0"
			}

			port := viper.GetString("port")

			// set up storage
			var store storage.BlobStorer
			storeDest := viper.GetString("storage")
			switch storeDest {
			case "memory":
				store = storage.NewMemoryStore()
			case "vault":
				client, err := VaultClient()
				if err != nil {
					GracefullyFail(err.Error())
				}
				store = storage.NewVaultStore(client, viper.GetString("vault-mount"))
				logrus.WithField("addr", viper.GetString("vault-addr")).Info("created vault client")
			default:
				GracefullyFail("unsupported storage engine")
			}

			logrus.WithFields(logrus.Fields{
				"address": address,
				"port":    port,
				"storage": storeDest,
			}).Info("listening")
			web.Serve(fmt.Sprintf("%s:%s", address, port), tf.New(store))
		},
	}

	CmdVersion = &cobra.Command{
		Use:   "version",
		Short: "print the version and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}
)
