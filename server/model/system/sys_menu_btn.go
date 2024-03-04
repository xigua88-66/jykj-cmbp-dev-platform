package system

import (
	"jykj-cmbp-dev-platform/server/global"
	"time"
)

type SysBaseMenuBtn struct {
	global.CmbpModel
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID string `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}

type Menus struct {
	//global.CmbpModel
	ID          string     `json:"id" gorm:"primary_key;not null;unique"`
	Type        int        `json:"type,omitempty" gorm:"not bull"`
	Name        string     `json:"name"`
	Level       int        `json:"level,omitempty"`
	OrderID     int        `json:"order_id"`
	Status      int        `json:"status"`
	LastMenu    string     `json:"last_menu,omitempty"`
	Url         string     `json:"url,omitempty"`
	RoleList    string     `json:"role_list,omitempty"`
	AssemblyUrl string     `json:"assembly_url,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	IsRouting   int        `json:"is_routing,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
	UpdateTime  *time.Time `json:"update_time,omitempty"`
	MenuID      string     `json:"menu_id" gorm:"-"`
	Children    []Menus    `json:"children" gorm:"-"`
	//Roles       []string  `gorm:"foreignKey:role_id;references:id"`
}

func (Menus) TableName() string {
	return "t_menus_info"
}
