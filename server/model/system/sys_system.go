package system

import (
	"jykj-cmbp-dev-platform/server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
