package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	//baseRouter := Router.Group("base")
	baseRouter := Router
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		//baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("/v1.0/login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
	}
	return baseRouter
}
