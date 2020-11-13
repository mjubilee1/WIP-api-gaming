package models

import (
	"fmt"
	"cloud.google.com/go/storage"
	"context"
	"api-gaming/internal/config"
)

// GetVideo - Return storage bucket handle data.
func GetVideo() *storage.BucketHandle {
	return config.StorageConn()
}

// ReadFile - Read a file from Google Cloud storage.
func ReadFile(fileName string) *storage.Reader {
	storage := config.StorageConn()

	obj := storage.Object(fileName).ReadCompressed(true)

	rc, err := obj.NewReader(context.Background())

	if err != nil {
		fmt.Println("readFile: unable to open file", err)
	}

	defer rc.Close()

	return rc
}