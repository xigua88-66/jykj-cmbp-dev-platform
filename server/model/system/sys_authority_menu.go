package system

import "time"

type SysMenu struct {
	SysBaseMenu
	MenuId      string                 `json:"menuId" gorm:"comment:菜单ID"`
	AuthorityId string                 `json:"-" gorm:"comment:角色ID"`
	Children    []SysMenu              `json:"children" gorm:"-"`
	Parameters  []SysBaseMenuParameter `json:"parameters" gorm:"foreignKey:SysBaseMenuID;references:MenuId"`
	Btns        map[string]string      `json:"btns" gorm:"-"`
}

type SysAuthorityMenu struct {
	MenuId      string `json:"menuId" gorm:"comment:菜单ID;column:sys_base_menu_id"`
	AuthorityId string `json:"-" gorm:"comment:角色ID;column:sys_authority_authority_id"`
}

func (s SysAuthorityMenu) TableName() string {
	return "sys_authority_menus"
}

type MenusByID struct {
	Type        int    `json:"type"`
	Icon        string `json:"icon"`
	MenuID      string `json:"menu_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	LastMenu    string `json:"last_menu"`
	RoleList    string `json:"role_list"`
	AssemblyUrl string `json:"assembly_url"`
	IsRouting   int    `json:"is_routing"`
}

type MenusItem struct {
	Children    []MenusItem `json:"children"`
	Icon        string      `json:"icon,omitempty"` // omitempty 表示如果字段为空，则在JSON中省略
	ID          int         `json:"id"`
	MenuID      string      `json:"menu_id"`
	Name        string      `json:"name"`
	OrderID     int         `json:"order_id"`
	Status      int         `json:"status"`
	URL         string      `json:"url,omitempty"`
	Type        int         `json:"type,omitempty"`
	RoleList    string      `json:"role_list,omitempty"`
	AssemblyUrl string      `json:"assembly_url,omitempty"`
	IsRouting   int         `json:"is_routing,omitempty"`
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
