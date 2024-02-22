package global

import (
	"errors"
	"github.com/gofrs/uuid/v5"
	"strings"
	"time"

	"gorm.io/gorm"
)

type CMBP_MODEL struct {
	ID        string         `gorm:"primarykey" json:"ID"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

// BeforeCreate GORM的模型钩子，在创建记录之前自动生成ID
func (m *CMBP_MODEL) BeforeCreate(tx *gorm.DB) (err error) {
	// 生成一个UUID
	uid, err := uuid.NewV4()
	if err != nil {
		return errors.New("生成UUID失败")
	}
	shortUUID := strings.ToUpper(strings.Join(strings.Split(uid.String(), "-"), ""))

	m.ID = shortUUID
	return nil
}
