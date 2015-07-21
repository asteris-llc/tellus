package storage

import (
	"errors"
)

var (
	ErrNoBlob   = errors.New("blob not found")
	ErrBadValue = errors.New("bad data in value")
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

type BlobStorer interface {
	BlobGetter
	BlobSetter
	BlobDeleter
}
