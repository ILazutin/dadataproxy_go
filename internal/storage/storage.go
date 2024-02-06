package storage

import "errors"

type Storage interface {
	Save(string, interface{}) error
	Read(string) (interface{}, error)
	ReadAllKeys() ([]string, error)
}

var (
	ErrUrlNotFound = errors.New("url not found")
)
