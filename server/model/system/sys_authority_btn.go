package system

type SysAuthorityBtn struct {
	AuthorityId      string         `gorm:"comment:角色ID"`
	SysMenuID        string         `gorm:"comment:菜单ID"`
	SysBaseMenuBtnID string         `gorm:"comment:菜单按钮ID"`
	SysBaseMenuBtn   SysBaseMenuBtn ` gorm:"comment:按钮详情"`
}
