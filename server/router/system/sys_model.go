package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type ModelsOptionRouter struct {
}

func (s *ModelsOptionRouter) InitModelsOptionRouter(Router *gin.RouterGroup) {

	modelsRouter := Router
	modelsRouterV9 := Router.Group("/v1.9")
	modelsRouterV5 := Router.Group("/v1.5")
	modelsRouterV7 := Router.Group("/v1.7")
	modelsRouterV2 := Router.Group("/v1.2")
	modelsRouterV12 := Router.Group("/v1.12")
	modelsApi := v1.ApiGroupApp.SystemApiGroup.ModelOptionApi
	{
		modelsRouter.GET("/v1.4/modelfield", modelsApi.GetModelField)      // 获取model小类
		modelsRouter.DELETE("/v1.6/cancel_upload", modelsApi.CancelUpload) // 取消上传
		modelsRouter.GET("/v1.0/upload_model", modelsApi.NothingToDo)      //前端无意义接口
		//modelsRouter.GET("/v1.7/model_check", modelsApi.ModelCheck)
		modelsRouter.GET("/v1.10/jupyter_notebook", modelsApi.JupyterNoteBook)

		modelsRouterV2.GET("models", modelsApi.GetModelStore)           // 模型仓库
		modelsRouterV2.POST("model", modelsApi.UploadModel)             // 线下构建上传模型
		modelsRouterV2.DELETE("model/:model_id", modelsApi.DeleteModel) // 仓库模型删除

		modelsRouterV5.GET("get_hardware_info", modelsApi.GetHardware)     // 获取硬件分类
		modelsRouterV5.GET("get_uuid", modelsApi.GetModelOpsUuid)          // 获取模型操作的uuid
		modelsRouterV5.DELETE("unpublish_model", modelsApi.UnPublishModel) // 取消模型操作的uuid
		modelsRouterV5.GET("test_free_application", modelsApi.GetTestFreeApplication)
		modelsRouterV5.POST("test_free_application", modelsApi.PostTestFreeApplication)
		//modelsRouterV5.PUT("test_free_application", modelsApi.PutTestFreeApplication)       // 取消模型操作的uuid
		//modelsRouterV5.DELETE("test_free_application", modelsApi.DeleteTestFreeApplication) // 取消模型操作的uuid
		modelsRouterV5.POST("upload_file", modelsApi.UploadFile)        // 模型仓库上传模型文件
		modelsRouterV5.GET("get_new_dirs", modelsApi.GetAIModelDirTree) // 模型修改获取目录树
		modelsRouterV5.GET("runtime", modelsApi.RunTime)                //镜像管理

		modelsRouterV7.GET("model_check", modelsApi.CheckName)                // 模型名称重复性校验
		modelsRouterV7.PUT("model", modelsApi.PutModelBusiness)               // 模型业务模型信息修改
		modelsRouterV7.GET("model_business_parm", modelsApi.GetModelBusiness) // 线下模型新增-自动获取业务模型参数及类型

		modelsRouterV9.GET("algorithm", modelsApi.GetAlgorithm)               // 获取算法
		modelsRouterV9.GET("update_AIMoniterend", modelsApi.GetAutoUpdateEnd) // End自动更新任务
		modelsRouterV9.GET("model_kind", modelsApi.GetModelKind)              // 获取模型大类
		modelsRouterV9.GET("industry", modelsApi.GetIndustry)                 // 获取行业信息
		modelsRouterV9.POST("hot_module", modelsApi.AddHotModule)             // 获取行业信息

		modelsRouterV12.GET("model_list", modelsApi.GetModelList)             // 模型市场
		modelsRouterV12.GET("relate_model_train", modelsApi.ModelTrainRelate) // 模型关联训练数据信息

	}
}
