package system

import (
	"jykj-cmbp-dev-platform/server/global"
	"time"
)

type SysAuthority struct {
	CreatedAt       time.Time       // 创建时间
	UpdatedAt       time.Time       // 更新时间
	DeletedAt       *time.Time      `sql:"index"`
	AuthorityId     string          `json:"authorityId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName   string          `json:"authorityName" gorm:"comment:角色名"`                                    // 角色名
	ParentId        *string         `json:"parentId" gorm:"comment:父角色ID"`                                       // 父角色ID
	DataAuthorityId []*SysAuthority `json:"dataAuthorityId" gorm:"many2many:sys_data_authority_id;"`
	Children        []SysAuthority  `json:"children" gorm:"-"`
	SysBaseMenus    []SysBaseMenu   `json:"menus" gorm:"many2many:sys_authority_menus;"`
	Users           []SysUser       `json:"-" gorm:"many2many:sys_user_authority;"`
	DefaultRouter   string          `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}

func (SysAuthority) TableName() string {
	return "sys_authorities"
}

// Roles 角色信息表
type Roles struct {
	global.CmbpModel
	Name      string      `gorm:"uniqueIndex;size:9;not null" json:"name"`
	RoleName  string      `gorm:"uniqueIndex;size:100" json:"role_name"` //角色描述
	Remarks   string      `gorm:"size:100" json:"remarks"`               // 备注
	Flag      bool        `gorm:"default:null" json:"flag"`              // 角色状态 0 禁用 1 启用
	UserRoles []UserRoles `gorm:"foreignKey:role_id;references:id"`
}

func (Roles) TableName() string {
	return "t_roles_info"
}

// UserRoles 用户角色关联表
type UserRoles struct {
	global.CmbpModel
	UserID   string `gorm:"index;size:32"`
	RoleID   string `gorm:"index;size:32"`
	MineCode string `gorm:"mine_code"`
	Role     Roles  `gorm:"foreignKey:RoleID"`
	User     Users  `gorm:"foreignKey:UserID"`
}

func (UserRoles) TableName() string {
	return "t_user_roles"
}
