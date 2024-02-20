package system

import "time"

type SysMenu struct {
	SysBaseMenu
	MenuId      string                 `json:"menuId" gorm:"comment:菜单ID"`
	AuthorityId uint                   `json:"-" gorm:"comment:角色ID"`
	Children    []SysMenu              `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuId"`
	Btns        map[string]uint        `json:"btns" gorm:"-"`
}

type SysAuthorityMenu struct {
	MenuId      string `json:"menuId" gorm:"comment:菜单ID;column:sys_base_menu_id"`
	AuthorityId string `json:"-" gorm:"comment:角色ID;column:sys_authority_authority_id"`
}

func (s SysAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}

type RoleMenus struct {
	Id         string    `json:"id"`
	RoleId     string    `json:"role_id"`
	MenuId     string    `json:"menu_id"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (RoleMenus) TableName() string {
	return "t_role_menus"
}
