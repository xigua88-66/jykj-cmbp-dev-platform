package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if strings.Index(f.Name, "..") > -1 {
			return []string{}, fmt.Errorf("%s 文件名不合法", f.Name)
		}
		fpath := filepath.Join(destDir, f.Name)
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return []string{}, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return []string{}, err
			}
		}
	}
	return paths, nil
}

func ZipFiles(filename string, files []string, oldForm, newForm string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = newZipFile.Close()
	}()

	zipWriter := zip.NewWriter(newZipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	// 把files添加到zip中
	for _, file := range files {

		err = func(file string) error {
			zipFile, err := os.Open(file)
			if err != nil {
				return err
			}
			defer zipFile.Close()
			// 获取file的基础信息
			info, err := zipFile.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// 使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
			header.Name = strings.Replace(file, oldForm, newForm, -1)

			// 优化压缩
			// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			if _, err = io.Copy(writer, zipFile); err != nil {
				return err
			}
			return nil
		}(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func CompressZip(sourcePath, targetPath, targetFile string, delSourcePath bool) error {
	// 检查源目录是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("源目录不存在")
	}

	// 创建目标目录（如果不存在的话）
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	target := filepath.Join(targetPath, targetFile+".zip")

	// 创建ZIP文件
	zf, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("创建ZIP文件时出错: %v", err)
	}
	defer zf.Close()

	zipWriter := zip.NewWriter(zf)
	defer zipWriter.Close()

	// 遍历源目录并添加到ZIP
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 创建ZIP中的目录条目
			fh, err := zip.FileInfoHeader(info)
			fh.Name = relPath + "/"
			if err != nil {
				return err
			}
			_, err = zipWriter.CreateHeader(fh)
			if err != nil {
				return err
			}
		} else {
			// 添加文件到ZIP
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			fh, err := zip.FileInfoHeader(info)
			fh.Name = relPath
			if err != nil {
				return err
			}

			w, err := zipWriter.CreateHeader(fh)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("文件压缩失败: %v", err)
	}

	// 如果需要，删除源文件目录
	if delSourcePath {
		if err := os.RemoveAll(sourcePath); err != nil {
			fmt.Printf("删除源目录时出错: %v\n", err)
		}
	}
	return nil
}
