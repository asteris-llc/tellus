package main

import (
	"github.com/Sirupsen/logrus"
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

	// config defaults and args
	CmdServe.Flags().String("address", ":4000", "address to run server on")
	viper.BindPFlag("address", CmdServe.Flags().Lookup("address"))
}

func main() {
	CmdTellusRoot.AddCommand(CmdServe)
	CmdTellusRoot.Execute()
}

var (
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
			address, err := cmd.Flags().GetString("address")
			if err != nil {
				panic(err)
			}

			logrus.WithField("addr", address).Info("listening")
			web.Serve(address)
		},
	}
)
