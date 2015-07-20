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
	viper.BindPFlag("address", CmdServe.Flags().Lookup("address"))
	viper.BindEnv("address")

	CmdServe.Flags().String("port", "4000", "port to run server on")
	viper.BindPFlag("port", CmdServe.Flags().Lookup("port"))
	viper.BindEnv("port", "PORT")

	CmdServe.Flags().String("storage", "memory", "how to store state given to the server")
	viper.BindPFlag("storage", CmdServe.Flags().Lookup("storage"))
	viper.BindEnv("storage")
}

func main() {
	defer Recovery()
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
			// get simple configs
			address := viper.GetString("address")
			port := viper.GetString("port")

			// set up storage
			var store storage.BlobStorer
			switch viper.GetString("storage") {
			case "memory":
				store = storage.NewMemoryStore()
			default:
				GracefullyFail("unsupported storage engine")
			}

			logrus.WithField("addr", address).Info("listening")
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
