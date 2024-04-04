package utils

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"jykj-cmbp-dev-platform/server/global"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateDir
//@description: 批量创建文件夹
//@param: dirs ...string
//@return: err error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			global.CMBP_LOG.Debug("create directory" + v)
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				global.CMBP_LOG.Error("create directory"+v, zap.Any(" error:", err))
				return err
			}
		}
	}
	return err
}

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: FileMove
//@description: 文件移动供外部调用
//@param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
//@return: err error

func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	revoke := false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

func FileCopy(src, dst string) error {
	// 读取源文件内容
	srcData, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// 将文件内容写入目标文件
	err = ioutil.WriteFile(dst, srcData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func DeLFile(filePath string) error {
	return os.RemoveAll(filePath)
}

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: TrimSpace
//@description: 去除结构体空格
//@param: target interface (target: 目标结构体,传入必须是指针类型)
//@return: null

func TrimSpace(target interface{}) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
}

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

// CopyDir 拷贝A目录下的所有东西到B目录下
func CopyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err = os.MkdirAll(dst, info.Mode())
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}

		for _, file := range files {
			srcPath := filepath.Join(src, file.Name())
			dstPath := filepath.Join(dst, file.Name())
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	} else {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyEnd(dest, zipPasswd, processor string) error {
	filter := []string{"file_operation_86.so", "file_operation_arm.so"}
	for _, f := range filter {
		fObj := path.Join(dest, f)
		_, err := os.Stat(fObj)
		if os.IsExist(err) {
			err = os.Remove(fObj)
			if err != nil {
				return err
			}
		}
	}
	end := "/home/models/AIMonitorEnd"
	err := CopyDir(end, dest)
	if err != nil {
		return err
	}
	if zipPasswd != "" {
		_, err := AddEncryptFile(dest, zipPasswd, processor)
		if err != nil {
			return err
		}
	}
	return nil
}

func Plugins2So(dir, processor string) bool {
	dirList, err := os.ReadDir(dir)
	if err != nil {
		global.CMBP_LOG.Fatal(err.Error())
		return false
	}
	res := true
	for _, fs := range dirList {
		busPath := path.Join(dir, fs.Name())
		filter := []string{".idea", "__pycache__", ".ipynb_checkpoints"}
		found := false
		for _, ft := range filter {
			if fs.Name() == ft {
				found = true
				break
			}
		}
		if !found {
			res = ActionPlugins2So(busPath, processor, 0)
		}
	}
	return res
}

func ActionPlugins2So(dir, processor string, level int) bool {
	igNorFile := []string{"__init__.py"}
	if level == 0 {
		igNorFile = append(igNorFile, "config.py", "plugin_debug.py")
	}
	res := false
	entryPoint := "BusinessModel"
	fs, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return res
	}
	filter := []string{".idea", "__pycache__", ".ipynb_checkpoints"}
	for index, f := range fs {
		fObj := path.Join(dir, f.Name())

		// 判断为业务模型文件夹
		if f.IsDir() && !(f.Name()[:1] == ".") {
			found := false
			for _, ft := range filter {
				if f.Name() == ft {
					found = true
					break
				}
			}
			if !found {
				res = ActionPlugins2So(fObj, processor, level+1)
			}
			// 判断为业务模型入口文件
		} else if !(f.Name()[:1] == ".") && len(f.Name()) > 3 && f.Name()[len(f.Name())-3:] == ".py" {
			found := false
			for _, ft := range igNorFile {
				if f.Name() == ft {
					found = true
					break
				}
			}
			if !found {
				existsEntryPoint := false
				fContent, _ := os.ReadFile(fObj)
				if strings.Contains(string(fContent), entryPoint) {
					existsEntryPoint = true
				}
				if existsEntryPoint {
					res = ActionPlugins2So(dir, processor, index)

				}
			}
		}
		if !res {
			return res
		}
	}
	return true
}

type TreeNode struct {
	Label      string      `json:"label"`
	ID         string      `json:"id"`
	Children   []*TreeNode `json:"children,omitempty"`
	RelPath    string      `json:"relpath"`
	FirstMatch bool        `json:"first,omitempty"`
	Name       string      `json:"__name__"`
}

var idCounter int

func GetDirTree(path, start string, includePlugins int) []*TreeNode {
	dirList, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	nodes := make([]*TreeNode, 0)
	for _, entry := range dirList {
		name := entry.Name()
		entryPath := filepath.Join(path, name)
		rel, err := filepath.Rel(start, entryPath)
		if err != nil {
			return nil
		}
		node := &TreeNode{
			Label:   name,
			ID:      entryPath,
			RelPath: rel,
			Name:    entryPath,
		}

		switch {
		case entry.IsDir():
			if includePlugins > 0 {
				switch includePlugins {
				case 1:
					if !strings.Contains(entryPath, "AIModel") {
						continue
					}
				case 2:
					if !strings.Contains(entryPath, "plugins") {
						continue
					}
				}
				if (includePlugins == 1 && name == "AIModel" && !strings.Contains(path, "AIModel")) ||
					(includePlugins == 2 && name == "plugins" && !strings.Contains(path, "plugins")) {
					node.FirstMatch = true
				}
				node.Children = GetDirTree(entryPath, start, includePlugins)
			} else {
				node.Children = GetDirTree(entryPath, start, 0)
				if name == "app" && !strings.Contains(path, "app") {
					node.FirstMatch = true
				}
			}
			nodes = append(nodes, node)
		default:
			idCounter++
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func TreeToJson(treeNodes []*TreeNode) ([]byte, error) {
	return json.MarshalIndent(treeNodes, "", "  ")
}
