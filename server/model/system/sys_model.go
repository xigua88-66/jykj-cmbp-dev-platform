package system

import (
	"jykj-cmbp-dev-platform/server/global"
	"time"
)

type ModelAll struct {
	global.CMBP_MODEL
	IndustryCode            int    `gorm:"not null;default:1"`
	ModelType               string `gorm:"size:7;not null"`
	FieldCode               string `gorm:"size:50;not null"`
	ModelName               string `gorm:"size:50;not null"`
	ModelChineseName        string `gorm:"size:255"`
	ModelVersion            string `gorm:"size:20;not null"`
	ModelDescription        string `gorm:"size:500;not null"`
	TechnicalDescription    string `gorm:"size:255"`
	PerformanceDescription  string `gorm:"size:255"`
	Advantage               string `gorm:"type:text"`
	MD5                     string `gorm:"size:50"`
	EncryptPassword         string `gorm:"size:255"`
	HardwareType            int    // 硬件类型 0为GPU 1为煤矿大脑盒子 2为华为泰山+A300
	IsImage                 bool
	CMD                     string `gorm:"size:500"` // 识别模型的运行命令
	JsonURL                 string `gorm:"size:200"` // 识别模型输出的json流url
	ImgURL                  string `gorm:"size:200"` // 识别模型出入的图片流url
	BusinessList            string `gorm:"size:200"`
	BusinessParams          string `gorm:"type:text"`
	ModelPurpose            int
	BusinessAPI             string  `gorm:"type:text"`
	AiModelAPI              string  `gorm:"type:text"`
	OnBoot                  bool    // 是否开机自启动
	NeedGPU                 bool    // 是否需要GPU卡 0为不需要1为需要
	AuditState              int     // 审核状态 0为待审核，1为审核通过、-1为拒绝
	DenyReason              string  `gorm:"size:500"` // 拒绝原因
	User                    string  `gorm:"size:36"`
	IsRealChannel           string  `gorm:"size:50;default:null"` // 是否逐帧识别 空为逐帧 非空为隔帧
	AiModelName             string  `gorm:"size:50;default:null"` // 模型库名称 唯一 默认为空（旧模型）
	AiModelID               string  `gorm:"size:32"`              // 算法id 唯一
	TestStatus              int     `gorm:"default:null"`         // 测试状态
	Accuracy                int     // 模型准确率
	TestDuration            float64 // 模型测试时长
	MetadataUpdateFlag      int     `gorm:"default:0"` // 模型属性更新
	AiModelUpdateFlag       int     `gorm:"default:0"` // 识别模型更新状态
	BusinessModelUpdateFlag int     `gorm:"default:0"` // 业务模型更新状态
	RuntimeUpdateFlag       int     `gorm:"default:0"` // 运行镜像环境更新状态
	NewModelFlag            bool    // 判断是否为新模型
	Edition                 int     `gorm:"default:0"` // 默认为0，更新一次递增1
	AlgorithmID             int     // 算法id
	BuildWay                int     // 模型构建方式，1为手动构建，2为自动构建
	ModelKind               int
	UUID                    string    `gorm:"size:32"`
	CreateTime              time.Time `gorm:"default:current_timestamp"`
	UpdateTime              time.Time `gorm:"default:current_timestamp on update current_timestamp"`
	BusinessType            string    `gorm:"type:text"`
	IsProcess               int       // 模型是否需要审核
}
