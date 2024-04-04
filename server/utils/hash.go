package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}

func FileMD5(filePath string) (md5sum string, err error) {
	// 打开ZIP文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建一个新的MD5散列
	hash := md5.New()

	// 创建一个用于读取文件内容的缓冲读取器
	reader := bufio.NewReader(file)

	// 逐块读取文件内容并写入MD5哈希
	chunk := make([]byte, 1024*1024*10) // 假设每次读取10MB的数据块
	for {
		n, err := reader.Read(chunk)
		if err == io.EOF {
			break // 如果已到达文件末尾，退出循环
		}
		if err != nil {
			//panic(err)
			return "", err
		}

		// 将读取的数据块写入MD5哈希
		_, err = hash.Write(chunk[:n])
		if err != nil {
			//panic(err)
			return "", err
		}
	}

	// 计算最终的MD5摘要，并转换为16进制形式
	md5sum = fmt.Sprintf("%x", hash.Sum(nil))
	return md5sum, nil
}
