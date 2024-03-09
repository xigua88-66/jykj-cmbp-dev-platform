// 自动生成模板SysOperationRecord
package system

import (
	"time"

	"jykj-cmbp-dev-platform/server/global"
)

// 如果含有time.Time 请自行import time包
type SysOperationRecord struct {
	global.CmbpModel
	Ip           string        `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"`                                   // 请求ip
	Method       string        `json:"method" form:"method" gorm:"column:method;comment:请求方法"`                       // 请求方法
	Path         string        `json:"path" form:"path" gorm:"column:path;comment:请求路径"`                             // 请求路径
	Status       int           `json:"status" form:"status" gorm:"column:status;comment:请求状态"`                       // 请求状态
	Latency      time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:延迟" swaggertype:"string"` // 延迟
	Agent        string        `json:"agent" form:"agent" gorm:"column:agent;comment:代理"`                            // 代理
	ErrorMessage string        `json:"error_message" form:"error_message" gorm:"column:error_message;comment:错误信息"`  // 错误信息
	Body         string        `json:"body" form:"body" gorm:"type:text;column:body;comment:请求Body"`                 // 请求Body
	Resp         string        `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:响应Body"`                 // 响应Body
	UserID       string        `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id"`                    // 用户id
	User         SysUser       `json:"user"`
}

type FrontLog struct {
	ID         uint      `gorm:"primarykey" json:"ID"` // 主键ID
	MineCode   string    `gorm:"column:mine_code;size:9;not null"`
	Username   string    `gorm:"column:username;size:20;not null"`
	SubSystem  string    `gorm:"column:sub_system;size:20;not null;default:'cmbp'"`
	PageName   string    `gorm:"column:page_name;size:32;not null"`
	Mark       string    `gorm:"column:mark;size:32"`
	EnterTime  time.Time `gorm:"column:enter_time;not null"`
	LeaveTime  time.Time `gorm:"column:leave_time;not null"`
	CreateTime time.Time `gorm:"column:create_time;default:current_timestamp;autoCreateTime"`
}

func (FrontLog) TableName() string {
	return "sys_ops_log"
}
