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
