package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"io/ioutil"
	"os/user"
)

func VaultClient() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = viper.GetString("vault-addr")

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	token := viper.GetString("vault-token")
	if token == "" {
		token = tokenFromFilesystem()
		logrus.Debug("read vault token from disk")
	} else {
		logrus.Debug("read vault token from config")
	}
	client.SetToken(token)

	return client, nil
}

func tokenFromFilesystem() string {
	usr, err := user.Current()
	if err != nil { // no current user, no home directory!
		logrus.Warn("no current user set, cannot read token from filesystem")
		return ""
	}

	token, err := ioutil.ReadFile(usr.HomeDir + "/.vault-token")
	if err != nil { // couldn't read it, oh well
		logrus.WithField("error", err).Warn("could not read vault token from disk")
		return ""
	}

	if len(token) == 0 {
		logrus.Warn("empty vault token loaded")
	}

	return string(token)
}
