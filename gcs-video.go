package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"time"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
	objectName string
}

// Get Public address, make sure the bocket's ACL is set to public-read.
func (c *ClientUploader) GetPulicAddress() string {
	if len(c.objectName) == 0 {
		return ""
	}
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, c.uploadPath+c.objectName)
}

// Upload Image object
func (c *ClientUploader) UploadImage(file io.ReadCloser) error {
	c.objectName = buildFileName() + ".jpeg"
	return c.uploadFile(file, c.objectName)
}

// Upload video object
func (c *ClientUploader) UploadVideo(file io.ReadCloser) error {
	c.objectName = buildFileName() + ".mp4"
	return c.uploadFile(file, c.objectName)
}

// uploadFile uploads an object
func (c *ClientUploader) uploadFile(file io.ReadCloser, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + c.objectName).NewWriter(ctx)
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
