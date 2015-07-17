package storage

import (
	"errors"
)

var (
	ErrNoBlob = errors.New("blob not found")
)

type Blob []byte

type BlobGetter interface {
	Get(string) (Blob, error)
}

type BlobSetter interface {
	Set(string, Blob) error
}

type BlobDeleter interface {
	Delete(string) error
}
