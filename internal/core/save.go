package core

import (
	"fmt"
	"io"
	"strings"
)

type SaveResponseSt struct {
	Id  string
	Url string
}

func (c *Core) Save(data io.Reader, size int64, contentType string) (*SaveResponseSt, error) {
	objectId := c.GenerateObjectName()

	err := c.ost.PutObject(c.bucketName, objectId, data, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("fail to ost.PutObject: %w", err)
	}

	return &SaveResponseSt{
		Id:  objectId,
		Url: strings.ReplaceAll(c.urlTemplate, "{id}", objectId),
	}, nil
}
