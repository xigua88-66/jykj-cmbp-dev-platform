package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type JwtRouter struct{}

func (s *JwtRouter) InitJwtRouter(Router *gin.RouterGroup) {
	//jwtRouter := Router.Group("jwt")
	jwtRouter := Router
	jwtApi := v1.ApiGroupApp.SystemApiGroup.JwtApi
	{
		jwtRouter.DELETE("/v1.0/logout", jwtApi.JsonInBlacklist) // jwt加入黑名单
	}
}
