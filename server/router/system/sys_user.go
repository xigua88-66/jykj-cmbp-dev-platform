package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type UserRouter struct{}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	//userRouter := Router.Group("user").Use(middleware.OperationRecord())
	//userRouter := Router.Use(middleware.OperationRecord())
	userRouter := Router
	userRouterWithoutRecord := Router.Group("user")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("admin_register", baseApi.Register)               // 管理员注册账号
		userRouter.PUT("/v1.4/user_info", baseApi.ChangePassword)         // 用户修改密码
		userRouter.POST("setUserAuthority", baseApi.SetUserAuthority)     // 设置用户权限
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               // 删除用户
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)                // 设置用户信息
		userRouter.PUT("setSelfInfo", baseApi.SetSelfInfo)                // 设置自身信息
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) // 设置用户权限组
		userRouter.POST("resetPassword", baseApi.ResetPassword)           // 设置用户权限组

		userRouter.GET("v1.4/user_info", baseApi.GetUserInfo)             // 获取用户信息
		userRouter.GET("v1.4/message", baseApi.GetUserMsg)                // 获取用户信息
		userRouter.GET("/v1.0/users/:phone", baseApi.CMBPDataGetUserList) // 数据工厂获取用户列表
		userRouter.PUT("/v1.0/user", baseApi.EnableUser)                  // 启用和禁用用户
		//userRouter.GET("/v1.0/getusername", baseApi.GetUserName)          // 数据工厂获取用户名
		userRouter.GET("/v1.4/userlist", baseApi.GetUserList) // 启用和禁用用户

	}
	{
		//userRouterWithoutRecord.POST("getUserList", baseApi.GetUserList) // 分页获取用户列表
		userRouterWithoutRecord.GET("getUserInfo", baseApi.GetUserInfo) // 获取自身信息
	}
}

func (s *UserRouter) InitNotAuthRouter(Router *gin.RouterGroup) {
	userRouter := Router
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.GET("/v1.0/getusername", baseApi.GetUserName) // 数据工厂获取用户名
	}
}
