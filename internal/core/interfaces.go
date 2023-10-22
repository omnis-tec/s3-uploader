package core

import (
	"io"
)

type ostI interface {
	CreateBucket(name string) error
	PutObject(bucketName, name string, data io.Reader, size int64, contentType string) error
	GetObject(bucketName, name string) (io.ReadCloser, error)
	RemoveObject(bucketName, name string) error
}
