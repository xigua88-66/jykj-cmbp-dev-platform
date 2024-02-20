package main

import (
	"go.uber.org/zap"

	"jykj-cmbp-dev-platform/server/core"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.6.0
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.CMBP_VP = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.CMBP_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.CMBP_LOG)
	global.CMBP_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.CMBP_DB != nil {
		//initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.CMBP_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
