package system

import (
	"jykj-cmbp-dev-platform/server/global"
)

type JwtBlacklist struct {
	global.CMBP_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
