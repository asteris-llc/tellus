package storage

import (
	"encoding/base64"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os/user"
	"strings"
	"testing"
)

type VaultTestSuite struct {
	suite.Suite
	rootPath string
	client   *api.Client
}

func (s *VaultTestSuite) SetupSuite() {
	if testing.Short() {
		s.T().Skip("skipping integration test in short mode")
	}

	s.rootPath = "tellus-test"

	// read token
	usr, err := user.Current()
	s.Require().Nil(err)
	token, err := ioutil.ReadFile(usr.HomeDir + "/.vault-token")
	s.Require().Nil(err)

	// create config
	config := api.DefaultConfig()
	config.Address = strings.Replace(config.Address, "https", "http", 1)
	client, err := api.NewClient(config)
	client.SetToken(string(token))
	s.Require().Nil(err)
	s.client = client
}

func (s *VaultTestSuite) SetupTest() {
	err := s.client.Sys().Mount(s.rootPath, "generic", "Tellus testing")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *VaultTestSuite) TearDownTest() {
	err := s.client.Sys().Unmount(s.rootPath)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *VaultTestSuite) path(path string) string {
	return s.rootPath + "/" + path
}

///////////////////////////////////// TESTS /////////////////////////////////////

func (s *VaultTestSuite) TestGet() {
	path := s.path("x")
	data := Blob("abc")

	_, err := s.client.Logical().Write(path, VaultKV{"blob": data})
	s.Nil(err)

	store := NewVaultStore(s.client, s.rootPath)
	out, err := store.Get("x")
	s.Nil(err)
	s.Equal(out, data)
}

func (s *VaultTestSuite) TestSet() {
	path := s.path("x")
	data := Blob("abc")
	szd := base64.StdEncoding.EncodeToString(data)

	store := NewVaultStore(s.client, s.rootPath)
	err := store.Set("x", data)
	s.Nil(err)

	secret, err := s.client.Logical().Read(path)
	s.Nil(err)
	s.Equal(szd, secret.Data["blob"])
}

func (s *VaultTestSuite) TestDelete() {
	path := s.path("x")
	data := Blob("abc")

	_, err := s.client.Logical().Write(path, VaultKV{"blob": data})
	s.Nil(err)

	store := NewVaultStore(s.client, s.rootPath)
	err = store.Delete("x")
	s.Nil(err)

	secret, err := s.client.Logical().Read(path)
	s.Nil(err)
	s.Nil(secret)
}

/////////////////////////////////// END TESTS ///////////////////////////////////

func TestVaultTestSuite(t *testing.T) {
	suite.Run(t, new(VaultTestSuite))
}
