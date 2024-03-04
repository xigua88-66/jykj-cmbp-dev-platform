package system

import (
	"jykj-cmbp-dev-platform/server/global"
)

type JwtBlacklist struct {
	global.CmbpModel
	Jwt string `gorm:"type:text;comment:jwt"`
}
