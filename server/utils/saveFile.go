package utils

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveFile(f *multipart.FileHeader, dest string) error {
	obj, err := f.Open()
	if err != nil {
		return err
	}
	defer obj.Close()

	target, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer target.Close()
	_, err = io.Copy(target, obj)
	if err != nil {
		return err
	}
	return nil
}

func GetFileSize(name string) int64 {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return 0
	}
	file, err := os.Open(name)
	if err != nil {
		return 0
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}
