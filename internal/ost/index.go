package ost

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type St struct {
	Client *minio.Client
}

func New(uri, keyId, key string, secure bool) (*St, error) {
	client, err := minio.New(uri, &minio.Options{
		Creds:  credentials.NewStaticV4(keyId, key, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New error: %w", err)
	}

	return &St{
		Client: client,
	}, nil
}

func (o *St) CreateBucket(name string) error {
	ctx := context.Background()

	// check if already exists
	exists, err := o.Client.BucketExists(ctx, name)
	if err != nil {
		return fmt.Errorf("minio.BucketExists error: %w", err)
	}
	if exists {
		return nil
	}

	err = o.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("minio.MakeBucket error: %w", err)
	}

	return nil
}

func (o *St) PutObject(bucketName, name string, data io.Reader, size int64, contentType string) error {
	if size <= 0 {
		dataContent, err := io.ReadAll(data)
		if err != nil {
			return fmt.Errorf("io.ReadAll error: %w", err)
		}
		size = int64(len(dataContent))
		data = bytes.NewReader(dataContent)
	}

	_, err := o.Client.PutObject(context.Background(), bucketName, name, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("minio.FPutObject error: %w", err)
	}

	return nil
}

func (o *St) GetObject(bucketName, name string) (io.ReadCloser, error) {
	result, err := o.Client.GetObject(context.Background(), bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio.GetObject error: %w", err)
	}

	return result, nil
}

func (o *St) RemoveObject(bucketName, name string) error {
	err := o.Client.RemoveObject(context.Background(), bucketName, name, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio.RemoveObject error: %w", err)
	}

	return nil
}
