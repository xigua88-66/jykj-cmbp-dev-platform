package system

import (
	"github.com/gin-gonic/gin"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/common/response"
	"jykj-cmbp-dev-platform/server/model/system"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	"jykj-cmbp-dev-platform/server/utils"
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
	data, err := modelService.GetModelField(params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	} else {
		response.OkWithData(data, c)
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
