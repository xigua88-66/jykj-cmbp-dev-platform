package global

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"strings"
	"time"
)

type CmbpModel struct {
	ID         string     `gorm:"primaryKey;type:varchar(32)" json:"id"`                                   // 主键ID
	CreateTime *time.Time `gorm:"default:current_timestamp" json:"create_time"`                            // 创建时间
	UpdateTime *time.Time `gorm:"default:current_timestamp;onUpdate:current_timestamp" json:"update_time"` // 更新时间
}

// BeforeCreate GORM的模型钩子，在创建记录之前自动生成ID
func (m *CmbpModel) BeforeCreate(tx *gorm.DB) (err error) {
	// 生成一个UUID
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.New("生成UUID失败")
	}
	shortUUID := strings.ToUpper(strings.Join(strings.Split(uid.String(), "-"), ""))
	m.ID = shortUUID
	now := time.Now()
	if m.CreateTime == nil {
		m.CreateTime = &now
	}
	if m.UpdateTime == nil {
		m.UpdateTime = &now
	}
	return nil
}

func (m *CmbpModel) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdateTime = &now
	return nil
}
