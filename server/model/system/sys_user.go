package system

import (
	"github.com/gofrs/uuid/v5"
	"jykj-cmbp-dev-platform/server/global"
	"time"
)

type SysUser struct {
	global.CmbpModel
	UUID        uuid.UUID      `json:"uuid" gorm:"index;comment:用户UUID"`                                                     // 用户UUID
	Username    string         `json:"userName" gorm:"index;comment:用户登录名"`                                                  // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                             // 用户登录密码
	NickName    string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                            // 用户昵称
	SideMode    string         `json:"sideMode" gorm:"default:dark;comment:用户侧边主题"`                                          // 用户侧边主题
	HeaderImg   string         `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	BaseColor   string         `json:"baseColor" gorm:"default:#fff;comment:基础颜色"`                                           // 基础颜色
	ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:活跃颜色"`                                      // 活跃颜色
	AuthorityId string         `json:"authorityId" gorm:"default:888;comment:用户角色ID"`                                        // 用户角色ID
	Authority   SysAuthority   `json:"authority" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
	Phone       string         `json:"phone"  gorm:"comment:用户手机号"`                     // 用户手机号
	Email       string         `json:"email"  gorm:"comment:用户邮箱"`                      // 用户邮箱
	Enable      int            `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
}

func (SysUser) TableName() string {
	return "sys_users"
}

type Users struct {
	global.CmbpModel
	Account        int       `json:"account" gorm:"default:NULL;comment:账号"`
	Password       string    `json:"password"`
	MineCode       string    `json:"mine_code"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Phone          int       `json:"phone"`
	Token          string    `json:"token"`
	IsActive       bool      `json:"is_active"`
	RootDisable    bool      `json:"root_disable"`
	MoveFlag       int       `json:"move_flag"`
	ExpireTime     time.Time `json:"expire_time"`
	ExpireLoginNum int       `json:"expire_login_num"`
	DingAccount    int       `json:"ding_account"`
}

func (Users) TableName() string {
	return "t_user_info"
}

type MineRegistry struct {
	MineCode        string    `json:"mine_code"`
	MineFullname    string    `json:"mine_fullname"`
	MineShortname   string    `json:"mine_shortname"`
	MineCapacity    int       `json:"mine_capacity"`
	MinePersonTotal int       `json:"mine_person_total"`
	UserFlag        bool      `json:"user_flag"`
	VerifyFlag      int       `json:"verify_flag"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Longitude       float32   `json:"longitude"`
	Latitude        float32   `json:"latitude"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}

func (MineRegistry) TableName() string {
	return "t_mine_register"
}
