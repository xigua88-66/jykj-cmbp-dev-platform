package middleware

import (
	"github.com/gin-gonic/gin"
	"jykj-cmbp-dev-platform/server/service"
)

var casbinService = service.ServiceGroupApp.SystemServiceGroup.CasbinService

// CasbinHandler 拦截器
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//if global.CMBP_CONFIG.System.Env != "develop" {
		//	waitUse, _ := utils.GetClaims(c)
		//	//获取请求的PATH
		//	path := c.Request.URL.Path
		//	obj := strings.TrimPrefix(path, global.CMBP_CONFIG.System.RouterPrefix)
		//	// 获取请求方法
		//	act := c.Request.Method
		//	// 获取用户的角色
		//	sub := waitUse.AuthorityId
		//	e := casbinService.Casbin() // 判断策略中是否存在
		//	success, _ := e.Enforce(sub, obj, act)
		//	if !success {
		//		response.FailWithDetailed(gin.H{}, "权限不足", c)
		//		c.Abort()
		//		return
		//	}
		//}
		c.Next()
	}
}
