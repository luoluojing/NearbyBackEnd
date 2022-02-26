package main

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
)

const (
	BUCKET_NAME = "aroundworld"
)

// bucket权限是如何识别的？
// google cloud 账号。程序运行在GCE机器上，GCS通过google cloud账号识别

func saveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	//client, err := storage.NewClient(cts).withCredentail(xxx)
	//如果程序在自己local电脑上， 需要在上面建立credential链接
	//IAM 机器账号 权限设置
	if err != nil {
		return "", err
	}

	object := client.Bucket(BUCKET_NAME).Object(objectName)
	wc := object.NewWriter(ctx)

	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	//allow all users to have acces; ACL: across control list
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {

		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	return attrs.MediaLink, nil
}
