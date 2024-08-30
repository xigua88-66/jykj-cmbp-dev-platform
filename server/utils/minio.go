package utils

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"time"
)

type MinIO struct {
}

func (m *MinIO) MinioClient() *minio.Client {
	endpoint := "172.24.1.71:9000"
	accessKeyID := "admin123456"
	secretAccessKey := "admin123456"

	// 初始化Minio客户端
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

// StreamUpload 流式上传
func (m *MinIO) StreamUpload(bucketName, objectName, filePath string) error {
	if bucketName == "" {
		bucketName = "obs-isf"
	}
	client := m.MinioClient()
	// 流式上传文件
	_, err := client.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Download 下载至本地
func (m *MinIO) Download(bucketName, objectName, localFilePath string) error {
	client := m.MinioClient()
	// 下载文件到本地
	err := client.FGetObject(context.Background(), bucketName, objectName, localFilePath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Successfully downloaded %s to %s\n", objectName, localFilePath)
	return nil
}

// Delete 删除文件
func (m *MinIO) Delete(bucketName, objectName string) error {
	client := m.MinioClient()
	err := client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Successfully deleted %s\n", objectName)
	return nil
}

// DownloadUrl 获取下载链接
func (m *MinIO) DownloadUrl(bucketName, objectName string) (u *url.URL, err error) {
	client := m.MinioClient()
	t := 30 * time.Second
	fmt.Println(t)
	u, err = client.PresignedGetObject(context.Background(), bucketName, objectName, t, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return u, nil
}

// ListObject 列举文件列表
func (m *MinIO) ListObject(bucketName string) []string {
	client := m.MinioClient()
	objNames := []string{}
	// 获取文件列表
	for object := range client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{}) {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		objNames = append(objNames, object.Key)
		log.Printf("Found object: %s\n", object.Key)
	}
	return objNames
}

//func main() {
//	var m MinIO
//	//list := m.ListObject("obs-isf")
//	//fmt.Println(list)
//	u, err := m.DownloadUrl("obs-isf", "data/test/README.md")
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println(u)
//}
