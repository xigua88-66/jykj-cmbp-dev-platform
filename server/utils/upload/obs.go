package upload

import (
	"mime/multipart"

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
