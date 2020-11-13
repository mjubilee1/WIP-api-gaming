package config

import (
	"context"
	"log"
	"cloud.google.com/go/storage"	
	"api-gaming/internal/util"
)

var (
	bucketname = util.ViperEnvVariable("GOOGLE_BUCKET_NAME")
	storageBucket *storage.BucketHandle
)

// InitGoogle - Initializing google
func InitGoogle() {
	storageClient, err := storage.NewClient(context.Background())

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	storageBucket = storageClient.Bucket("playground.sleepless-gamers.com")
}

func StorageConn() (*storage.BucketHandle) {
	return storageBucket
}