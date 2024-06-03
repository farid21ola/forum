package domain

import (
	"errors"
	"github.com/farid21ola/forum/storage"
)

var (
	ErrBadCredentials  = errors.New("invalid username or password")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("unauthorized")
)

type Domain struct {
	Storage storage.Storage
}

func NewDomain(storage storage.Storage) *Domain {
	return &Domain{Storage: storage}
}
