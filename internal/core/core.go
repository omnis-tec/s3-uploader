package core

import (
	"fmt"

	"github.com/google/uuid"
)

type Core struct {
	bucketName  string
	urlTemplate string
	ost         ostI
}

func NewCore(bucketName, urlTemplate string, ost ostI) (*Core, error) {
	err := ost.CreateBucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("fail to ost.CreateBucket: %w", err)
	}

	return &Core{
		bucketName:  bucketName,
		urlTemplate: urlTemplate,
		ost:         ost,
	}, nil
}

func (c *Core) GenerateObjectName() string {
	return uuid.NewString()
}
