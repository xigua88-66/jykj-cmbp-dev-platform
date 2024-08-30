package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	//"google.golang.org/genproto/googleapis/rpc/context"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/common/response"
	"jykj-cmbp-dev-platform/server/model/system"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	"jykj-cmbp-dev-platform/server/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ModelOptionApi struct {
}

func (m *ModelOptionApi) GetModelField(c *gin.Context) {
	var params systemReq.ModelFiled
	err := c.BindQuery(&params)
	if err != nil {
		return
	}
	err = utils.Verify(params, utils.ModelFieldVerify)
	if err != nil {
		return
	}
	rspData, err := modelService.GetModelField(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		//rspData := []interface{}{}
		//rspData = append(rspData, data)
		response.OkWithData(rspData, c)
		return
	}
}

func (m *ModelOptionApi) GetAlgorithm(c *gin.Context) {
	var params systemReq.AlgorithmRqe
	err := c.BindQuery(&params)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	algorithmList, err := modelService.GetAlgorithmLogic(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(algorithmList, c)
	}
}

func (m *ModelOptionApi) GetModelList(c *gin.Context) {
	var params systemReq.ModelListReq
	c.BindQuery(&params)
	err := utils.Verify(params, utils.GetModelListVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var userId = utils.GetUserID(c)
	modelList, err := modelService.GetModelList(params, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(modelList, c)
}

func (m *ModelOptionApi) GetModelStore(c *gin.Context) {
	var params systemReq.ModelStoreRqe
	c.BindQuery(&params)
	role := utils.GetUserRole(c)
	var user system.Users
	global.CMBP_DB.Where("id = ?", utils.GetUserID(c)).Find(&user)
	modelStore, err := modelService.GetModelStore(params, user, role)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(modelStore, c)
		return
	}
}

func (m *ModelOptionApi) UploadModel(c *gin.Context) {
	var params systemReq.UploadModelStoreReq

	err := c.ShouldBind(&params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(params, utils.UploadModelStoreVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	videoFile, err := c.FormFile("video_file")
	if err != nil {
		msg := "视频文件不能为空"
		if !strings.Contains(err.Error(), "http: no such file") {
			msg = err.Error()
		}
		response.FailWithMessage(msg, c)
		return
	}
	imgFile, err := c.FormFile("model_img")
	if err != nil {
		msg := "图片文件不能为空"
		if !strings.Contains(err.Error(), "http: no such file") {
			msg = err.Error()
		}
		response.FailWithMessage(msg, c)
		return
	}
	ext := path.Ext(videoFile.Filename)
	if ext == "" || strings.ToUpper(ext[1:]) != "MP4" {
		response.FailWithMessage("视频格式应为mp4", c)
		return
	}
	ext = path.Ext(imgFile.Filename)
	if ext == "" || strings.ToUpper(ext[1:]) != "JPG" {
		response.FailWithMessage("图片格式应为jpg", c)
		return
	}
	userId := utils.GetUserID(c)
	rspData, err := modelService.UploadModel(params, videoFile, imgFile, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(rspData, c)
		return
	}
}

func (m *ModelOptionApi) DeleteModel(c *gin.Context) {
	modelID := c.Param("model_id")
	if modelID == "" {
		response.FailWithMessage("路由中模型ID为必传项", c)
		return
	}
	resData, err := modelService.DeleteModel(modelID)
}

func (m *ModelOptionApi) GetAutoUpdateEnd(c *gin.Context) {
	userID := utils.GetUserID(c)
	rspData, err := modelService.GetAutoUpdateTask(userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) GetModelKind(c *gin.Context) {
	var params systemReq.GetModelKind
	c.BindQuery(&params)
	rspData, err := modelService.GetModelKind(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) GetHardware(c *gin.Context) {
	var params systemReq.GetHardWare
	c.BindQuery(&params)
	rspData, err := modelService.GetHardWare(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) GetModelOpsUuid(c *gin.Context) {
	userID := utils.GetUserID(c)
	rspData, err := modelService.GetModelOpsUuid(userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) GetIndustry(c *gin.Context) {
	var params systemReq.GetIndustry
	c.BindQuery(&params)
	rspData, err := modelService.GetIndustry(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) UnPublishModel(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		response.FailWithMessage("未检测到参数uuid", c)
	}
	token := utils.GetToken(c)
	err := modelService.UnPublishModel(uuid, token)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(nil, "删除成功", c)
	}
}

func (m *ModelOptionApi) CancelUpload(c *gin.Context) {
	uuid := c.Query("uuid")
	if uuid == "" {
		response.FailWithMessage("未检测到uuid参数", c)
	}
	userID := utils.GetUserID(c)
	err := modelService.CancelUpload(uuid, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.Ok(c)
	}
}

func (m *ModelOptionApi) GetTestFreeApplication(c *gin.Context) {
	var params systemReq.GetTestFreeApply
	c.BindQuery(&params)
	rspData, err := modelService.GetTestFreeApplication(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithData(rspData, c)
	}
}

func (m *ModelOptionApi) PostTestFreeApplication(c *gin.Context) {

}

func (m *ModelOptionApi) UploadFile(c *gin.Context) {
	var params systemReq.UploadFile
	c.ShouldBind(&params)
	userId := utils.GetUserID(c)
	taskId := params.TaskId
	chunk := params.Chunk
	if chunk == "" {
		chunk = "0"
	}
	filename := fmt.Sprintf("%s%s", taskId, chunk)
	fileDir := fmt.Sprintf("/home/models/fileSave/%s", userId)
	_, err := os.Stat(fileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(fileDir, 0755)
	}
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.SaveFile(file, filepath.Join(fileDir, filename))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
	return
}

func (m *ModelOptionApi) GetAIModelDirTree(c *gin.Context) {
	var params systemReq.GetModelDirTree
	c.ShouldBindQuery(&params)
	userId := utils.GetUserID(c)

	if params.Path != "" && params.Path != "[]" {
		modelDir := fmt.Sprintf("/home/tmp/%s/%s", userId, params.UUID)
		_, err := os.Stat(modelDir)
		if os.IsExist(err) {
			dirTree := utils.GetDirTree(modelDir, modelDir, 2)
			response.OkWithData(dirTree, c)
		} else {
			response.FailWithMessage("该文件夹不存在，请传递正确的路径", c)
		}
	} else if params.OBSPath != "" {
		modelDir := "/OBS/" + params.OBSPath
		_, err := os.Stat(modelDir)
		if os.IsExist(err) {
			dirTree := utils.GetDirTree(modelDir, modelDir, 0)
			response.OkWithData(dirTree, c)
		} else {
			response.FailWithMessage("该文件夹不存在，请传递正确的路径", c)
		}

	} else if params.WeightsID != "" {
		var w system.WeightsManagement
		global.CMBP_DB.Model(&system.WeightsManagement{}).Where("id = ?", params.WeightsID).First(&w)
		fileDir := filepath.Join("/OBS/WeightsLibrary", w.ID, "data")
		_, err := os.Stat(fileDir)
		if os.IsExist(err) {
			dirTree := utils.GetDirTree(fileDir, fileDir, 0)
			// TODO add url
			response.OkWithData(dirTree, c)
		} else {
			response.FailWithMessage("数据错误"+err.Error(), c)
		}
	} else if params.Offline != 0 {
		key := params.UUID + "_AIModel_tmp_dirs"
		data, _ := global.CMBP_REDIS.Get(context.Background(), key).Bytes()
		var dirTree []*utils.TreeNode
		json.Unmarshal(data, &dirTree)
		response.OkWithData(dirTree, c)
	} else {
		modelDir := fmt.Sprintf("/home/tmp/%s/%s", userId, params.UUID)
		_, err := os.Stat(modelDir)
		var dirTree []*utils.TreeNode
		if os.IsExist(err) {
			dirTree = utils.GetDirTree(modelDir, modelDir, 0)
		} else {
			dirTree = utils.GetDirTree(modelDir, modelDir, 1)
		}
		response.OkWithData(dirTree, c)
	}
	return
}

func (m *ModelOptionApi) CheckName(c *gin.Context) {
	var params systemReq.CheckName
	c.ShouldBind(&params)
	userId := utils.GetUserID(c)
	rspData, err := modelService.CheckName(params, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(rspData, c)
	return
}

func (m *ModelOptionApi) GetModelBusiness(c *gin.Context) {
	modelId := c.Query("model_id")
	if modelId == "" {
		response.FailWithMessage("参数不能为空", c)
		return
	}
	var model system.ModelAll
	global.CMBP_DB.Model(&system.ModelAll{}).Where("audit_state = 1 AND id = ?", modelId).First(&model)
	if model.ID == "" {
		response.FailWithMessage("模型不存在", c)
		return
	}
	zipName := model.ModelName + "V" + model.ModelVersion
	zipPath := fmt.Sprintf("/home/OBS/models/models/%s/%s.zip", zipName, zipName)
	_, err := os.Stat(zipPath)
	if os.IsNotExist(err) {
		response.FailWithMessage("OBS下业务模型包不存在", c)
		return
	}
	businessList, err := utils.GetModelBusinessList(zipPath)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(businessList, c)
	return

}

func (m *ModelOptionApi) ModelTrainRelate(c *gin.Context) { // TODO
	f := map[string]interface{}{
		"param_desc": "init.weights",
		"param_name": "init.weights"}
	files := []interface{}{}
	files = append(files, f)
	rspData := map[string]interface{}{
		"is_need_relate": true,
		"relation_files": files,
		"replace_weights": map[string]interface{}{
			"desc":        "需要升级的权重",
			"target_file": "",
			"type":        0,
		},
	}
	response.OkWithData(rspData, c)
}

func (m *ModelOptionApi) PutModelBusiness(c *gin.Context) {
	modelId := c.PostForm("model_id")
	businessDict := c.PostForm("business_dict")
	if modelId == "" {
		response.FailWithMessage("modelId不能为空", c)
		return
	}
	var modelAll system.ModelAll
	global.CMBP_DB.Model(&system.ModelAll{}).Where("id = ?", modelId).First(&modelAll)
	if modelAll.ID == "" {
		response.FailWithMessage("该模型不存在", c)
		return
	}
	var businessList []map[string]interface{}
	if businessDict != "" {
		err := json.Unmarshal([]byte(businessDict), &businessList)
		if err != nil || len(businessList) == 0 {
			response.FailWithMessage("business_dict不符合规则", c)
			return
		}
		busStr := make([]string, 0)
		busType := make(map[string]interface{})
		busDict := make(map[string]interface{})
		busApi := make(map[string]interface{})
		for _, dict := range businessList {
			busName := dict["business_name"].(string)
			busParams := dict["business_params"]
			busApis, ok := dict["business_apis"].([]interface{})
			if !ok {
				busApis = []interface{}{}
			}
			busType[busName] = dict["business_type"]
			busDict[busName] = busParams
			busApi[busName] = busApis
			busStr = append(busStr, busName)
		}

		modelAll.BusinessList = strings.Join(busStr, "|")
		params, _ := json.Marshal(busDict)
		api, _ := json.Marshal(busApi)
		bType, _ := json.Marshal(busType)
		modelAll.BusinessParams = string(params)
		modelAll.BusinessAPI = string(api)
		modelAll.BusinessType = string(bType)

		var modelInfo []system.Model
		global.CMBP_DB.Model(&system.Model{}).Where("model_all_id = ?", modelAll.ID).Find(&modelInfo)

		for _, m := range modelInfo {
			m.BusinessList = modelAll.BusinessList
			m.BusinessParams = modelAll.BusinessParams
			m.BusinessType = modelAll.BusinessType
			global.CMBP_DB.Save(&m)
		}
		global.CMBP_DB.Save(&modelAll)
	}
	response.Ok(c)
}

func (m *ModelOptionApi) NothingToDo(c *gin.Context) {
	response.Ok(c)
	return
}

func (m *ModelOptionApi) AddHotModule(c *gin.Context) {
	// TODO
	response.Ok(c)
	return
}
