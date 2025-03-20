package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"jykj-cmbp-dev-platform/server/global"
	"os"
	"os/exec"
	"path/filepath"
)

func GetModelBusinessList(zipObj string) (resData interface{}, err error) {
	dir, err := os.Getwd()
	cmd := exec.Command("python", filepath.Join(dir, "utils", "get_business_list.py"), zipObj)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	output, err := cmd.Output()
	if err != nil {
		global.CMBP_LOG.Error("获取业务模型参数错误：" + err.Error() + stderr.String())
		return nil, errors.New("获取业务模型参数错误：" + err.Error() + stderr.String())
	}
	var j interface{}
	err = json.Unmarshal(output, &j)
	if err != nil {
		global.CMBP_LOG.Error("获取业务模型参数错误：" + err.Error())
		return nil, errors.New("获取业务模型参数错误：" + err.Error())
	}
	return j, nil
}
