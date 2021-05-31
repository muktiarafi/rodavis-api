package utils

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/muktiarafi/rodavis-api/internal/api"
)

type ImagePersistence interface {
	Save(fileName string, image multipart.File) (string, error)
}

type CloudImagePersistence struct {
	context.Context
	*storage.BucketHandle
	BucketName string
}

func NewCloudImagePersistence(
	ctx context.Context,
	bucket *storage.BucketHandle,
	bucketName string,
) ImagePersistence {
	return &CloudImagePersistence{
		Context:      ctx,
		BucketHandle: bucket,
		BucketName:   bucketName,
	}
}

func (i *CloudImagePersistence) Save(fileName string, image multipart.File) (string, error) {
	defer image.Close()
	wc := i.BucketHandle.Object(fileName).NewWriter(i.Context)
	defer wc.Close()
	if _, err := io.Copy(wc, image); err != nil {
		return "", api.NewExceptionWithSourceLocation(
			"CloudImagePersistence.Save",
			"io.Copy",
			err,
		)
	}

	imageURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", i.BucketName, wc.Name)
	return imageURL, nil
}

type LocalImagePersistence struct {
	SaveDirectory string
}

func NewLocalImagePersistence(saveDirectory string) ImagePersistence {
	return &LocalImagePersistence{
		SaveDirectory: saveDirectory,
	}
}

func (i *LocalImagePersistence) Save(fileName string, image multipart.File) (string, error) {
	defer image.Close()
	savePath := filepath.Join(i.SaveDirectory, fileName)
	f, err := os.OpenFile(
		savePath,
		os.O_WRONLY|os.O_CREATE,
		os.ModePerm,
	)
	const op = "LocalImagePersistence.Save"
	if err != nil {
		return "", api.NewExceptionWithSourceLocation(
			op,
			"os.OpenFile",
			err,
		)
	}
	defer f.Close()

	if _, err := io.Copy(f, image); err != nil {
		return "", api.NewExceptionWithSourceLocation(
			op,
			"io.Copy",
			err,
		)
	}

	return savePath, nil
}
