package utils

import (
	"fmt"
	"jykj-cmbp-dev-platform/server/global"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@%^*"
)

// UniqueRandomStr 生成元素唯一的字符串
func UniqueRandomStr() string {
	rand.Seed(time.Now().UnixNano())

	charSet := []byte(letterBytes)
	num := RandomInt(18, 30)
	if num > len(charSet) {
		return ""
	}

	shuffled := make([]string, 0, num)
	for len(shuffled) < num {
		idx := rand.Intn(len(charSet))
		ch := charSet[idx]
		charSet = append(charSet[:idx], charSet[idx+1:]...)
		shuffled = append(shuffled, string(ch))
	}
	return strings.Join(shuffled, "")
}

// Encrypt 字符创串进行加密
func Encrypt(s string, move bool) []byte {

	// 将字符串s转换为UTF-8编码的字节切片
	b := []byte(s)
	n := len(b)

	// 创建一个双倍长度的新字节切片c用于存放加密后的数据
	c := make([]byte, n*2)

	temp := make([]byte, 20) // 临时缓冲区
	j := 0

	for i := 0; i < n; i++ {
		b1 := b[i]
		b2 := b1 ^ 64 // 对b1进行异或操作
		c1 := b2 % 16
		c2 := b2 / 16 // Go中除法默认就是整数除法，不需要特别指定floor操作
		c1 += 65      // 将c1映射到A-P的ASCII区间
		c2 += 65      // 将c2映射到A-P的ASCII区间
		c[j] = c1
		c[j+1] = c2
		j += 2
	}

	if !move { // 如果不进行移动操作，则直接返回加密后的c
		return c
	}

	// 生成随机移位数组
	var randInt []int
	for i := 0; i < 10; i++ {
		randInt = append(randInt, rand.Intn(11)+10) // 生成10到20之间的随机数
	}

	// 进行随机移位操作
	for _, i := range randInt {
		copy(temp[:i], c[:i])
		copy(c[:n*2-i], c[i:])
		copy(c[n*2-i:], temp[:i])
		ReverseByte(c) // 自定义一个反转字节切片的函数
	}

	// 将随机移位数组转化为字符串并加密
	appendStr := ""
	for _, v := range randInt {
		appendStr += strconv.Itoa(v) + "|"
	}
	appendStr = appendStr[:len(appendStr)-1] // 去除最后一个'|'
	appended := Encrypt(appendStr, false)    // 递归加密随机移位数组字符串

	// 合并最终加密结果并返回
	bX := []byte("X")
	sumByte := append(c, bX...)
	sumByte = append(sumByte, appended...)
	return sumByte
}

// ReverseByte 对字节数组进行翻转
func ReverseByte(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 复原字符串的移位
func _decryptCore(s string) string {
	c := []byte(s)
	b := make([]byte, len(c)/2)
	j := 0
	for i := 0; i < len(b); i++ {
		c1 := c[j]
		c2 := c[j+1]
		j += 2
		c1 -= 65
		c2 -= 65
		b2 := c2*16 + c1
		b1 := b2 ^ 64
		b[i] = b1
	}
	return string(b)
}

// Decrypt 字符串解密
func Decrypt(s string) (string, error) {
	activeCodeList := strings.Split(s, "X")
	if len(activeCodeList) == 2 && activeCodeList[1] != "" {
		// 解密移位码
		coreRandStr := _decryptCore(activeCodeList[1])
		randInt := make([]int, 0)
		for _, iStr := range strings.Split(coreRandStr, "|") {
			i, err := strconv.Atoi(iStr)
			if err != nil {
				return "", err
			}
			randInt = append(randInt, i)
		}
		reversedIntRandInt := ReverseInt(randInt)
		temp := make([]byte, 20)
		activeCode := []byte(activeCodeList[0])
		for _, i := range reversedIntRandInt {
			ReverseByte(activeCode)
			copy(temp[:i], activeCode[len(activeCode)-i:])
			copy(activeCode[i:], activeCode[:len(activeCode)-i])
			copy(activeCode[:i], temp[:i])
		}

		unShifted := _decryptCore(string(activeCode))
		return unShifted, nil
	}
	return "", fmt.Errorf("无效的输入格式")
}

// ReverseInt 数字数组反转
func ReverseInt(a []int) []int {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return a
}

func Py2So(dest, encryptFile, processor string) error {
	py2Ccmd := fmt.Sprintf("/cmbp/Python370/bin/cython -3 %s", encryptFile)
	if processor == "arm" {
		py2Ccmd = fmt.Sprintf("cython -3 %s", encryptFile)
	}
	pythonDesc := "/cmbp/Python352/include/python3.5m"
	if processor == "arm" {
		pythonDesc = "/usr/local/include/python3.5m"
	}
	c2soCmd := fmt.Sprintf("gcc -shared -pthread -fPIC  -I %s -o %s.so %s.c", pythonDesc, encryptFile[:len(encryptFile)-3], encryptFile[:len(encryptFile)-3])

	rmPySoCmd := fmt.Sprintf("rm -f %s %s.c", encryptFile, encryptFile[:len(encryptFile)-3])

	Cmd := fmt.Sprintf("cd %s && %s && %s && %s", dest, py2Ccmd, c2soCmd, rmPySoCmd)

	if processor == "arm" {
		Cmd = fmt.Sprintf("/bin/bash -c 'cd %s && %s && %s && %s'", dest, py2Ccmd, c2soCmd, rmPySoCmd)
	}

	if processor == "x86" {
		cmd := exec.Command("/bin/bash", "-c", Cmd)
		out, err := cmd.Output()
		if err != nil {
			global.CMBP_LOG.Fatal(err.Error())
			return err
		}
		global.CMBP_LOG.Info("py转so成功" + string(out))
		return nil
	} else if processor == "arm" {
		cmd := exec.Command("docker", "exec", "python-arrch64", "/bin/bash", "-c", Cmd)
		_, err := cmd.CombinedOutput()
		if err != nil {
			global.CMBP_LOG.Fatal(err.Error())
			return err
		} else {
			global.CMBP_LOG.Info("arm-py转so成功")
			return nil
		}
	}
	return nil
}

func AddEncryptFile(dest, zipPasswd, processor string) (string, error) {
	encrypted := Encrypt(zipPasswd, true)
	encryptFile := "encrypt.py"
	file, _ := os.Create(path.Join(dest, encryptFile))
	_, err := file.WriteString(fmt.Sprintf("ciphertext = '%s'", encrypted))
	if err != nil {
		return "", err
	}
	err = Py2So(dest, encryptFile, processor)
	if err != nil {
		return "", err
	}

	license, _ := os.Create(path.Join(dest, "license.init"))
	_, err = license.WriteString("HOST_ID|APP_ID")
	if err != nil {
		return "", err
	}
	//Cmd := fmt.Sprintf("cd %s && zip -r -1 -P %s license.zip license.init && rm -f license.init", dest, zipPasswd)
	Cmd := fmt.Sprintf(`cd "%s" && zip -r -1 -P %s license.zip license.init && rm -f license.init`, dest, zipPasswd)
	cmd := exec.Command("/bin/bash", "-c", Cmd)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}
