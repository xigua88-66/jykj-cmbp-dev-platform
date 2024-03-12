package utils

import (
	"os"
)

func Unzip(src interface{}, dest string, delete bool) error {
	if delete {
		_, err := os.Stat(dest)
		if os.IsExist(err) {
			err := os.RemoveAll(dest)
			if err != nil {
				return err
			}
		}
	}
	delete = false
	return nil
}
