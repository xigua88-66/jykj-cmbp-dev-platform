package system

import "jykj-cmbp-dev-platform/server/global"

type SysBaseMenuBtn struct {
	global.CMBP_MODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID string `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}
