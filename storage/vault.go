package storage

import (
	"encoding/base64"
	"errors"
	"github.com/hashicorp/vault/api"
)

type VaultKV map[string]interface{}

var (
	ErrNoVaultKey = errors.New("no key in vault secret data")
)

type VaultStore struct {
	client   *api.Client
	blobKey  string
	PathRoot string
}

func NewVaultStore(client *api.Client, pathRoot string) *VaultStore {
	return &VaultStore{client, "blob", pathRoot}
}

func (vs *VaultStore) path(key string) string {
	return vs.PathRoot + "/" + key
}

func (vs *VaultStore) Get(key string) (Blob, error) {
	op := vs.client.Logical()

	secret, err := op.Read(vs.path(key))
	if err != nil {
		return Blob{}, err
	}
	if secret == nil {
		return Blob{}, ErrNoBlob
	}

	value, ok := secret.Data[vs.blobKey]
	if !ok {
		return Blob{}, ErrNoVaultKey
	}

	blob, ok := value.(string)
	if !ok {
		return Blob{}, ErrBadValue
	}

	result, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return Blob{}, ErrBadValue
	}

	return Blob(result), nil
}

func (vs *VaultStore) Set(key string, blob Blob) error {
	op := vs.client.Logical()

	_, err := op.Write(vs.path(key), VaultKV{vs.blobKey: blob})
	if err != nil {
		return err
	}

	return nil
}

func (vs *VaultStore) Delete(key string) error {
	_, err := vs.client.Logical().Delete(vs.path(key))
	return err
}
