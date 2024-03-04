package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type ModelsOptionRouter struct {
}

func (s *ModelsOptionRouter) InitModelsOptionRouter(Router *gin.RouterGroup) {

	modelsRouter := Router
	modelsApi := v1.ApiGroupApp.SystemApiGroup.ModelOptionApi
	{
		modelsRouter.GET("/v1.4/modelfield", modelsApi.GetModelField)             // 获取model小类
		modelsRouter.GET("/v1.9/algorithm", modelsApi.GetAlgorithm)               // 获取算法
		modelsRouter.GET("/v1.12/model_list", modelsApi.GetModelList)             // 模型市场
		modelsRouter.GET("/v1.2/models", modelsApi.GetModelStore)                 // 模型仓库
		modelsRouter.GET("/v1.9/update_AIMoniterend", modelsApi.GetAutoUpdateEnd) // End自动更新任务
		modelsRouter.GET("/v1.9/model_kind", modelsApi.GetModelKind)              // 获取模型大类
		modelsRouter.GET("/v1.5/get_hardware_info", modelsApi.GetHardware)        // 获取硬件分类
		modelsRouter.GET("/v1.5/get_uuid", modelsApi.GetModelOpsUuid)             // 获取模型操作的uuid
		modelsRouter.GET("/v1.9/industry", modelsApi.GetIndustry)                 // 获取行业信息

	}
}
