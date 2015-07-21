package main

import (
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
	}

	return client, nil
}

func tokenFromFilesystem() string {
	usr, err := user.Current()
	if err != nil { // no current user, no home directory!
		return ""
	}

	token, err := ioutil.ReadFile(usr.HomeDir + "/.vault-token")
	if err != nil { // couldn't read it, oh well
		return ""
	}

	return string(token)
}
