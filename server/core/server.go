package core

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/initialize"
	"jykj-cmbp-dev-platform/server/service/system"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.CMBP_CONFIG.System.UseMultipoint || global.CMBP_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	if global.CMBP_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.CMBP_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.CMBP_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.CMBP_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 jykj-cmbp-dev-platform
	煤矿大脑-人工智能基础平台
	当前版本:v2.6.0
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, address)
	global.CMBP_LOG.Error(s.ListenAndServe().Error())
}
