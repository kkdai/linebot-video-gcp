package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"log"
	"time"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file io.ReadCloser, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func buildFileName() string {
	return time.Now().Format("20060102150405")
}

func (c *ClientUploader) MakePublic(object string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	acl := client.Bucket(c.bucketName).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}
	log.Printf("Blob %v is now publicly accessible.\n", object)
	return nil
}
