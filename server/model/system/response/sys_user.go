package response

import (
	"jykj-cmbp-dev-platform/server/model/system"
	"time"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	//User      system.Users `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	IsCloud   bool   `json:"is_cloud"`
	IsExpire  int    `json:"is_expire"`
	Role      string `json:"role"`
}

type DataFactoryUserListRsp struct {
	ID         string      `json:"id"`
	MineCode   string      `json:"mine_code"`
	Username   string      `json:"username"`
	Phone      string      `json:"phone"`
	IsActive   int         `json:"is_active"`
	Roles      interface{} `json:"roles"`
	Email      string      `json:"email"`
	CreateTime string      `json:"create_time"`
}

type AdminGetUserList struct {
	ID            string     `json:"id"`
	MineCode      string     `json:"mine_code"`
	MineShortname string     `json:"mine_shortname"`
	Username      string     `json:"username"`
	Phone         string     `json:"phone"`
	RootDisable   int        `json:"root_disable"` // 需要自定义UnreadItem类型
	IsActive      int        `json:"is_active"`
	Roles         string     `json:"roles"`
	RoleID        string     `json:"role_id"` // 需要在GetRoleID方法中实现逻辑
	MoveFlag      bool       `json:"move_flag"`
	Email         string     `json:"email"`
	CreateTime    string     `json:"create_time"` // 若需要与原API保持一致，此处也可以存储为字符串
	ExpireTime    string     `json:"expire_time"`
	CreateAt      time.Time  `json:"-"`
	ExpireAt      *time.Time `json:"-"`
}

// RootDisableInt  返回处理过的RootDisable字段
func (u AdminGetUserList) RootDisableInt() int {
	if u.RootDisable == 1 {
		return 1
	} else {
		return -1
	}
}

// FormatCreateTime 格式化CreateTime字段为字符串
func (u AdminGetUserList) FormatCreateTime() string {
	return u.CreateAt.Format("2006-01-02 05:04:05")
}

// FormatExpireTime 格式化ExpireTime字段为字符串
func (u AdminGetUserList) FormatExpireTime() string {
	if u.ExpireAt != nil {
		return u.ExpireAt.Format("2006-01-02 05:04:05")
	} else {
		return ""
	}
}
