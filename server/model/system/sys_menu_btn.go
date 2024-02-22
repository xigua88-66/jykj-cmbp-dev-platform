package system

import (
	"jykj-cmbp-dev-platform/server/global"
	"time"
)

type SysBaseMenuBtn struct {
	global.CMBP_MODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID string `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}

type Menus struct {
	global.CMBP_MODEL
	ID          string    `json:"id" gorm:"primary_key;not null;unique"`
	Type        int       `json:"type" gorm:"not bull"`
	Name        string    `json:"name"`
	Level       int       `json:"level"`
	OrderID     int       `json:"order_id"`
	Status      int       `json:"status"`
	LastMenu    string    `json:"last_menu"`
	Url         string    `json:"url"`
	RoleList    string    `json:"role_list"`
	AssemblyUrl string    `json:"assembly_url"`
	Icon        string    `json:"icon"`
	IsRouting   bool      `json:"is_routing"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
	Roles       []string  `gorm:"foreignKey:role_id;references:id"`
}

func (Menus) TableName() string {
	return "t_menus_info"
}
