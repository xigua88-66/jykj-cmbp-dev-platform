package system

import (
	"time"
)

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     uint            `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId        *uint           `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}

type Roles struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	RoleName   string    `json:"role_name"` //角色描述
	Remarks    string    `json:"remarks"`   // 备注
	Flag       bool      `json:"flag"`      // 角色状态 0 禁用 1 启用
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (Roles) TableName() string {
	return "t_roles_info"
}

type UserRoles struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	RoleId     string    `json:"role_id"`
	MineCode   string    `json:"mine_code"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	User       string    `json:"user"`
	Role       string    `json:"role"`
}

func (UserRoles) TableName() string {
	return "t_user_roles"
}
