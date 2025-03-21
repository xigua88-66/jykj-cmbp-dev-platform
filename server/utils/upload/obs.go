package upload

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/pkg/errors"
	"jykj-cmbp-dev-platform/server/global"
)

var HuaWeiObs = new(Obs)

type Obs struct{}

func NewHuaWeiObsClient() (client *obs.ObsClient, err error) {
	return obs.New(global.CMBP_CONFIG.HuaWeiObs.AccessKey, global.CMBP_CONFIG.HuaWeiObs.SecretKey, global.CMBP_CONFIG.HuaWeiObs.Endpoint)
}

func (o *Obs) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// var open multipart.File
	open, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer open.Close()
	filename := file.Filename
	input := &obs.PutObjectInput{
		PutObjectBasicInput: obs.PutObjectBasicInput{
			ObjectOperationInput: obs.ObjectOperationInput{
				Bucket: global.CMBP_CONFIG.HuaWeiObs.Bucket,
				Key:    filename,
			},
			HttpHeader: obs.HttpHeader{
				ContentType: file.Header.Get("content-type"),
			},
		},
		Body: open,
	}

	var client *obs.ObsClient
	client, err = NewHuaWeiObsClient()
	if err != nil {
		return "", "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}

	_, err = client.PutObject(input)
	if err != nil {
		return "", "", errors.Wrap(err, "文件上传失败!")
	}
	filepath := global.CMBP_CONFIG.HuaWeiObs.Path + "/" + filename
	return filepath, filename, err
}

func (o *Obs) DeleteFile(key string) error {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	input := &obs.DeleteObjectInput{
		Bucket: global.CMBP_CONFIG.HuaWeiObs.Bucket,
		Key:    key,
	}
	var output *obs.DeleteObjectOutput
	output, err = client.DeleteObject(input)
	if err != nil {
		return errors.Wrapf(err, "删除对象(%s)失败!, output: %v", key, output)
	}
	return nil
}

func (o *Obs) SignDownUrl(path string, expire int) (url string, err error) {
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return "", errors.Wrap(err, "获取华为对象存储对象失败!")
	}
	//type CreateSignedUrlInput struct {
	//	Method      HttpMethodType
	//	Bucket      string
	//	Key         string
	//	Policy      string
	//	SubResource SubResourceType
	//	Expires     int
	//	Headers     map[string]string
	//	QueryParams map[string]string
	//
	input := obs.CreateSignedUrlInput{
		Bucket:  global.CMBP_CONFIG.HuaWeiObs.Bucket,
		Key:     path,
		Method:  obs.HttpMethodGet,
		Expires: global.CMBP_CONFIG.CMBPBase.OssExpireSeconds,
	}

	signedUrl, err := client.CreateSignedUrl(&input)
	if err != nil {
		return "", err
	}
	return signedUrl.SignedUrl, nil
}

func (o *Obs) Download(bucket, objKey, localPath string) (err error) {
	if !strings.HasSuffix(objKey, "/") {
		objKey = "/" + objKey
	}
	client, err := NewHuaWeiObsClient()
	if err != nil {
		return err
	}

	if bucket == "" {
		bucket = global.CMBP_CONFIG.HuaWeiObs.Bucket
	}

	input := new(obs.DownloadFileInput)
	input.Bucket = bucket
	input.Key = objKey
	input.DownloadFile = localPath
	input.PartSize = 10 * 1024 * 1024
	input.EnableCheckpoint = true

	// 构造下载参数
	//input := &obs.DownloadFileInput{
	//	Bucket:       bucket,        // 必须指定存储桶
	//	Key:          objKey,        // 对象键（保持原始值）
	//	SaveAsStream: false,         // 明确指定保存方式
	//	DownloadFile: localPath,     // 本地保存路径
	//	PartSize:     10 * 1024 * 1024, // 分段下载大小（10MB）
	//	EnableCheckpoint: true,      // 启用断点续传
	//}

	resp, err := client.DownloadFile(input)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("OBS下载：%v到本地失败", objKey))
	}
	return nil
}
