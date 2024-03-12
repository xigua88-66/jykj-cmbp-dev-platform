package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type ModelsOptionRouter struct {
}

func (s *ModelsOptionRouter) InitModelsOptionRouter(Router *gin.RouterGroup) {

	modelsRouter := Router
	modelsRouterV19 := Router.Group("/v1.9")
	modelsRouterV15 := Router.Group("/v1.5")
	modelsApi := v1.ApiGroupApp.SystemApiGroup.ModelOptionApi
	{
		modelsRouter.GET("/v1.4/modelfield", modelsApi.GetModelField)      // 获取model小类
		modelsRouter.GET("/v1.12/model_list", modelsApi.GetModelList)      // 模型市场
		modelsRouter.GET("/v1.2/models", modelsApi.GetModelStore)          // 模型仓库
		modelsRouter.POST("/v1.2/model", modelsApi.UploadModel)            // 模型仓库
		modelsRouter.DELETE("/v1.6/cancel_upload", modelsApi.CancelUpload) // 模型仓库

		modelsRouterV15.GET("get_hardware_info", modelsApi.GetHardware)                  // 获取硬件分类
		modelsRouterV15.GET("get_uuid", modelsApi.GetModelOpsUuid)                       // 获取模型操作的uuid
		modelsRouterV15.DELETE("unpublish_model", modelsApi.UnPublishModel)              // 取消模型操作的uuid
		modelsRouterV15.GET("test_free_application", modelsApi.GetTestFreeApplication)   // 取消模型操作的uuid
		modelsRouterV15.POST("test_free_application", modelsApi.PostTestFreeApplication) // 取消模型操作的uuid
		//modelsRouterV15.PUT("test_free_application", modelsApi.PutTestFreeApplication)       // 取消模型操作的uuid
		//modelsRouterV15.DELETE("test_free_application", modelsApi.DeleteTestFreeApplication) // 取消模型操作的uuid

		modelsRouterV19.GET("algorithm", modelsApi.GetAlgorithm)               // 获取算法
		modelsRouterV19.GET("update_AIMoniterend", modelsApi.GetAutoUpdateEnd) // End自动更新任务
		modelsRouterV19.GET("model_kind", modelsApi.GetModelKind)              // 获取模型大类
		modelsRouterV19.GET("industry", modelsApi.GetIndustry)                 // 获取行业信息

	}
}
