package system

import (
	"fmt"
	"gorm.io/gorm"
	"jykj-cmbp-dev-platform/server/global"
	"path/filepath"
	"strconv"
	"time"
)

type ModelAll struct {
	global.CmbpModel
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
	BusinessClass           string `gorm:"type:text"`
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
	TestStatus              *int    `gorm:"default:null"`         // 测试状态
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
	UUID                    string                  `gorm:"size:32"`
	BusinessType            string                  `gorm:"type:text"`
	IsProcess               int                     // 模型是否需要审核
	BusinessModelAllProxy   []ModelAllBusinessModel `gorm:"foreignKey:ModelAllID" json:"business_model_all_proxy"`
}

func (ModelAll) TableName() string {
	return "t_model_all"
}

// ModelNameAndVersion 模型的名称前缀
func (m *ModelAll) ModelNameAndVersion() string {
	if m.ModelKind != 1 {
		return m.ModelName + "V" + m.ModelVersion
	} else {
		return m.ModelName + "V" + m.ModelVersion + "." + strconv.Itoa(m.Edition)
	}
}

// ModelZipFile 模型的zip路径
func (m *ModelAll) ModelZipFile() string {
	return filepath.Join(global.CMBP_CONFIG.CMBPBase.ModelPath, m.ModelNameAndVersion(), ".zip")
}

// 模型的图片路径
func (m *ModelAll) ModelImgFile() string {
	return filepath.Join(global.CMBP_CONFIG.CMBPBase.OssModelMedia, m.ModelNameAndVersion(), ".jpg")
}

func (m *ModelAll) ModelVideoFile() string {
	return filepath.Join(global.CMBP_CONFIG.CMBPBase.OssModelMedia, m.ModelNameAndVersion(), ".mp4")
}

type ModelDetails struct {
	global.CmbpModel
	IndustryCode            *int    `gorm:"type:int"`
	IndustryName            *string `gorm:"type:varchar(32)"`
	Developer               *string `gorm:"type:varchar(20)"`
	DeveloperPhone          *string `gorm:"type:varchar(11)"`
	Edition                 *int    `gorm:"type:int"`
	ModelName               *string `gorm:"type:varchar(50)"`
	ModelChineseName        *string `gorm:"type:varchar(255)"`
	ModelDescription        *string `gorm:"type:varchar(500)"`
	TechnicalDescription    *string `gorm:"type:varchar(255)"`
	PerformanceDescription  *string `gorm:"type:varchar(255)"`
	ModelVersion            *string `gorm:"type:varchar(20)"`
	ModelKind               *int    `gorm:"type:int"`
	ModelKindName           *string `gorm:"type:varchar(32)"`
	AlgorithmID             *int    `gorm:"type:int"`
	AlgorithmName           *string `gorm:"type:varchar(32)"`
	BuildWay                *int    `gorm:"type:int"`
	ModelType               *string `gorm:"type:varchar(7)"`
	ModelTypeName           *string `gorm:"type:varchar(50)"`
	ModelField              *string `gorm:"type:varchar(50)"`
	ModelFieldName          *string `gorm:"type:varchar(50)"`
	HardwareType            *int    `gorm:"type:int"`
	HardwareTypeName        *string `gorm:"type:varchar(200)"`
	HardwareTypeDesc        *string `gorm:"type:varchar(500)"`
	TestStatus              *int    `gorm:"type:int"`
	Accuracy                *int    `gorm:"type:int"`
	TestDuration            *int    `gorm:"type:int"`
	MetadataUpdateFlag      *int    `gorm:"type:int"`
	AiModelUpdateFlag       *int    `gorm:"type:int"`
	BusinessModelUpdateFlag *int    `gorm:"type:int"`
	RuntimeUpdateFlag       *int    `gorm:"type:int"`
	DownLoadCount           *int    `gorm:"type:int"`
	ViewCount               *int    `gorm:"type:int"`
	UserViewCount           *int    `gorm:"type:int"`
	VisitorViewCount        *int    `gorm:"type:int"`
	UsedCount               *int    `gorm:"type:int"`
	CollectionCount         *int    `gorm:"type:int"`
	TotalNum                *int    `gorm:"type:int"`
	AnnotatedSampleCount    *int    `gorm:"type:int"`
	BusinessType            *string `gorm:"type:text"`
	BusinessParams          *string `gorm:"type:text"`
	IsImage                 *int    `gorm:"type:int"`
	RuntimeID               *string `gorm:"type:varchar(255)"` // 假设runtime_id是字符串类型，根据实际情况调整
	ModelPurpose            *int    `gorm:"type:int"`
}

func (ModelDetails) TableName() string {
	return "v_new_model_list"
}

// ModelType 模型类型表结构体
type ModelType struct {
	global.CmbpModel
	ModelType      string  `gorm:"not null;size:7" json:"model_type"`
	ModelTypeDesc  string  `gorm:"not null;size:50" json:"model_type_desc"`
	UserID         *string `gorm:"size:32" json:"user_id"`
	UpdateUserID   *string `gorm:"size:32" json:"update_user_id"`
	IsImage        bool    `gorm:"default:false" json:"is_image"`
	WithWeights    bool    `gorm:"default:false" json:"with_weights"`
	IsCV           bool    `gorm:"default:false" json:"is_CV"`
	Databases      string  `gorm:"size:255" json:"databases"`
	Cmd            string  `gorm:"size:500" json:"cmd"`
	Params         string  `gorm:"size:200" json:"params"`
	JsonPort       int     `json:"json_port"`
	ImgPort        int     `json:"img_port"`
	JsonPortAppend string  `gorm:"size:20" json:"json_port_append"`
	ImgPortAppend  string  `gorm:"size:20" json:"img_port_append"`
	ImgURL         string  `gorm:"size:50" json:"img_url"`
	IsGPU          bool    `gorm:"default:false" json:"is_GPU"`
	ModelFieldID   string  `gorm:"ForeignKey:ModelFieldID;references:id on delete:CASCADE" json:"-"`

	// 外键关联，定义模型领域与模型类型的关系
	ModelField ModelField `gorm:"foreignKey:ModelFieldID;references:ID" json:"model_field,omitempty"`
}

func (ModelType) TableName() string {
	return "t_model_type"
}

// ModelField 模型领域表结构体
type ModelField struct {
	global.CmbpModel
	Name         string `gorm:"size:50" json:"name"`
	FieldNameEn  string `gorm:"size:32" json:"field_name_en"`
	IndustryCode int    `gorm:"not null;default:1" json:"industry_code"`
	Code         string `gorm:"size:50" json:"code"`
	FieldImgPath string `gorm:"type:text" json:"field_img_path"`

	// 定义一个反向一对多关系（一对多的反向）
	ModelTypes []ModelType `gorm:"foreignKey:ModelFieldID"`
}

func (ModelField) TableName() string {
	return "t_model_field"
}

type Industry struct {
	IndustryCode     int64     `gorm:"primaryKey;not null"`                                  // 行业编号，主键
	IndustryName     string    `gorm:"not null;size:32"`                                     // 行业名称
	IndustryDesc     *string   `gorm:"type:text"`                                            // 行业描述
	IndustryNameEN   string    `gorm:"not null;size:32"`                                     // 行业英文名称
	IndustryAbbrName string    `gorm:"not null;size:32"`                                     // 行业缩写名称
	MarketImgPath    string    `gorm:"not null;size:32"`                                     // 市场图片地址
	OutsideImgPath   string    `gorm:"not null;size:32"`                                     // 外部图片地址
	CreateTime       time.Time `gorm:"default:current_timestamp"`                            // 创建时间，默认当前时间
	UpdateTime       time.Time `gorm:"default:current_timestamp;onupdate:current_timestamp"` // 更新时间，默认为当前时间，并在更新时自动设置
}

func (Industry) TableName() string {
	return "t_industry"
}

// AlgorithmInfo 算法信息表结构体
type AlgorithmInfo struct {
	AlgorithmID        int64     `gorm:"primaryKey;not null"`                                  // 算法编号，主键
	ModelKind          int64     `gorm:"not null"`                                             // 算法分类编号
	AlgorithmName      string    `gorm:"not null;size:32"`                                     // 算法名称
	AlgorithmNameEN    *string   `gorm:"size:32"`                                              // 算法英文名称
	AlgorithmDesc      string    `gorm:"not null;size:255"`                                    // 算法描述
	AlgorithmImgPath   string    `gorm:"not null"`                                             // 图片地址
	AlgorithmVideoPath string    `gorm:"not null"`                                             // 视频地址
	CreateTime         time.Time `gorm:"default:current_timestamp"`                            // 创建时间，默认当前时间
	UpdateTime         time.Time `gorm:"default:current_timestamp;onupdate:current_timestamp"` // 更新时间，默认为当前时间，并在更新时自动设置
}

func (AlgorithmInfo) TableName() string {
	return "t_algorithm_info"
}

type ModelKind struct {
	ModelKind  int       `gorm:"primary_key;column:model_kind;not null"`
	KindName   string    `gorm:"column:kind_name;size:32;not null"`
	KindNameEn *string   `gorm:"column:kind_name_en;size:32"`
	KindDesc   string    `gorm:"column:kind_desc;size:255;not null"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (ModelKind) TableName() string {
	return "t_algorithm_kind"
}

type AssetsManagement struct {
	global.CmbpModel
	UserID          string      `gorm:"not null" json:"user_id"`                                                                        // 用户ID
	AssetsID        string      `gorm:"not null;index;foreignKey:AssetsIDReferences;references:id on delete:SET NULL" json:"assets_id"` // 资产ID，外键关联t_assets_type表的id字段，删除时设置为NULL
	AssetsName      string      `gorm:"not null;size:32" json:"assets_name"`                                                            // 资产名称
	AssetsDesc      string      `gorm:"not null;size:32" json:"assets_desc"`                                                            // 资产描述
	AssetsMediaURL  *string     `gorm:"size:32" json:"assets_media_url"`                                                                // 资产图片/视频地址
	AssetsTypeID    string      `gorm:"not null;size:32" json:"assets_type_id"`                                                         // 资产类型：1-模型资产，2-数据资产，3-场景资产，4-业务模型，5-硬件类型（其余请参考t_assets_type表）
	CollectionCount *int        `json:"collection_count"`                                                                               // 收藏/点赞数量
	ViewCount       *int        `json:"view_count"`                                                                                     // 查看次数
	DownloadCount   *int        `json:"download_count"`                                                                                 // 下载次数
	IsDeleted       bool        `gorm:"default:false" json:"is_deleted"`                                                                // 是否已被作者删除，默认为false
	ExtraKwargs     *string     `gorm:"type:text" json:"extra_kwargs"`                                                                  // 其余参数，以JSON字符串形式存储
	AssetsType      *AssetsType `gorm:"constraint:OnDelete:SET NULL;"`                                                                  // 定义与t_assets_type表的关联关系
}

func (AssetsManagement) TableName() string {
	return "t_assets_management"
}

// AssetsType AI资产类型结构体
type AssetsType struct {
	global.CmbpModel
	AssetsNo   int    `gorm:"not null;unique;autoIncrement" json:"assets_no"` // 资产编号，唯一且自增
	AssetsName string `gorm:"not null;size:32" json:"assets_name"`            // 资产类型名称
	AssetsDesc string `gorm:"not null;size:32" json:"assets_desc"`            // 资产类型描述

}

func (AssetsType) TableName() string {
	return "t_assets_type"
}

// AssetsRecord AI资产申请记录/审批结果结构体
type AssetsRecord struct {
	global.CmbpModel
	AssetsID    string  `gorm:"not null;size:32" json:"assets_id"`      // 资产的ID
	AssetsName  string  `gorm:"not null;size:32" json:"assets_name"`    // 资产名称
	AssetsDesc  *string `gorm:"type:text" json:"assets_desc"`           // 资产描述
	ApplyType   int     `gorm:"not null;default:-1" json:"apply_type"`  // 申请类型，其值来源于t_assets_type表中的assets_type_no字段
	ApplyDetail string  `gorm:"not null;size:32" json:"apply_detail"`   // 申请记录中要申请的具体资产ID（模型、数据、场景等）
	UserID      string  `gorm:"not null;size:32" json:"user_id"`        // 用户ID
	ApplyStatus int     `gorm:"not null;default:0" json:"apply_status"` // 审批状态，0-待审批，1-审批通过，2-审批拒绝
}

func (AssetsRecord) TableName() string {
	return "t_assets_record"
}

// ModelMarketList 模型市场列表结构体
type ModelMarketList struct {
	global.CmbpModel
	IndustryCode           int    `gorm:"column:industry_code;nullable" json:"industry_code"`                              // 行业代码
	ModelName              string `gorm:"column:model_name;size:50;nullable" json:"model_name"`                            // 模型名称
	ModelChineseName       string `gorm:"column:model_chinese_name;size:255;nullable" json:"model_chinese_name"`           // 模型中文名
	ModelDescription       string `gorm:"column:model_description;size:500;nullable" json:"model_description"`             // 模型描述
	TechnicalDescription   string `gorm:"column:technical_description;size:255;nullable" json:"technical_description"`     // 技术描述
	PerformanceDescription string `gorm:"column:performance_description;size:255;nullable" json:"performance_description"` // 性能描述
	ModelKind              int    `gorm:"column:model_kind;nullable" json:"model_kind"`                                    // 模型种类
	AlgorithmID            int    `gorm:"column:algorithm_id;nullable" json:"algorithm_id"`                                // 算法ID
	BuildWay               int    `gorm:"column:build_way;nullable" json:"build_way"`                                      // 构建方式
	ModelType              string `gorm:"column:model_type;size:7;nullable" json:"model_type"`                             // 模型类型
	ModelField             string `gorm:"column:model_field;size:50;nullable" json:"model_field"`                          // 模型领域
	Accuracy               *int   `gorm:"column:accuracy;nullable" json:"accuracy"`                                        // 准确率
	TestDuration           *int   `gorm:"column:test_duration;nullable" json:"test_duration"`                              // 测试时长
	DownloadCount          int    `gorm:"column:down_load_count;nullable" json:"download_count"`                           // 下载次数
	ViewCount              int    `gorm:"column:view_count;nullable" json:"view_count"`                                    // 查看次数
	UserViewCount          int    `gorm:"column:user_view_count;nullable" json:"user_view_count"`                          // 用户查看次数
	VisitorViewCount       int    `gorm:"column:visitor_view_count;nullable" json:"visitor_view_count"`                    // 访客查看次数
	UsedCount              int    `gorm:"column:used_count;nullable" json:"used_count"`                                    // 使用次数
	CollectionCount        int    `gorm:"column:collection_count;nullable" json:"collection_count"`                        // 收藏次数
	TotalNum               int    `gorm:"column:total_num;nullable" json:"total_num"`                                      // 总数量
	ModelPurpose           int    `gorm:"column:model_purpose;nullable" json:"model_purpose"`                              // 模型用途
	Edition                int    `gorm:"column:edition;nullable" json:"edition"`                                          // 版本号
	ModelVersion           string `gorm:"column:model_version;size:20;nullable" json:"model_version"`                      // 模型版本

}

// TableName 定义数据库表名
func (ModelMarketList) TableName() string {
	return "v_model_market_list"
}

type Model struct {
	global.CmbpModel
	MineCode         string        `gorm:"column:mine_code;size:9"`
	ModelType        string        `gorm:"column:model_type;size:7;not null"`
	ModelName        string        `gorm:"column:model_name;size:50;not null"`
	ModelChineseName string        `gorm:"column:model_chinese_name;size:50"`
	ModelVersion     string        `gorm:"column:model_version;size:20;not null"`
	ModelDescription string        `gorm:"column:model_description;size:100;not null"`
	MD5              string        `gorm:"column:md5;size:50"`
	IsGPU            bool          `gorm:"column:is_GPU"`
	HardwareType     int           `gorm:"column:hardware_type"`
	IsImage          bool          `gorm:"column:is_image"`
	Cmd              string        `gorm:"column:cmd;size:500"`
	JSONURL          string        `gorm:"column:json_url;size:200"`
	ImgURL           string        `gorm:"column:img_url;size:200"`
	BusinessList     string        `gorm:"column:business_list;size:200"`
	BusinessParams   string        `gorm:"column:business_params;text"`
	OnBoot           bool          `gorm:"column:on_boot"`
	NeedGPU          bool          `gorm:"column:need_gpu"`
	SyncFlag         *int          `gorm:"column:sync_flag"`
	ModelAllID       string        `gorm:"column:model_all_id;size:32"`
	ModelConfig      []ModelConfig `gorm:"foreignKey:ModelID;references:ID;cascade:save"` // 在Gorm中处理一对多关系需要额外定义ModelConfig关联的model
	User             string        `gorm:"column:user;size:32;ForeignKey:TUserInfo.ID;references:ID;on_delete:CASCADE"`
	Channels         int           `gorm:"column:channels"`
	Deadline         string        `gorm:"column:deadline;size:30"`
	PayStatus        int           `gorm:"column:pay_status"`
	IsRealChannel    *string       `gorm:"column:is_real_channel;size:50"`
	Accuracy         int           `gorm:"column:accuracy"`
	TestDuration     float64       `gorm:"column:test_duration"`
	NewModelFlag     int           `gorm:"column:new_model_flag"`
	BusinessType     string        `gorm:"column:business_type;text"`
}

func (Model) TableName() string {
	return "t_model_info"
}

func (m *Model) ModelNameAndVersion() string {
	return m.ModelName + "V" + m.ModelVersion
}

func (m *Model) ModelZipFile() string {
	return filepath.Join(global.CMBP_CONFIG.CMBPBase.MineModelDir, m.MineCode, m.ModelChineseName+".zip")
}

type ModelConfig struct {
	ID                 string              `gorm:"primary_key;default:generate_uuid;size:32"`
	MineCode           string              `gorm:"column:mine_code;size:9";not null`
	EndID              string              `gorm:"column:end_id;size:32;ForeignKey:TEndInfo.ID;references:ID;on_delete:CASCADE"`
	ModelID            string              `gorm:"primaryKey;column:model_info_id;ForeignKey:Model.ID;references:ID;onDelete:CASCADE"`
	ModelInfo          *Model              `gorm:"foreignKey:ModelID;references:ID;cascade:delete"`
	ModelType          string              `gorm:"column:model_type;size:7";not null`
	ModelName          string              `gorm:"column:model_name;size:20";not null`
	ModelVersion       string              `gorm:"column:model_version;size:20";not null`
	ModelSubVersion    string              `gorm:"column:model_sub_version;size:10"`
	Path               string              `gorm:"column:path;size:50";not null;default:'/home/AIAgent/models/'"`
	ZipName            string              `gorm:"column:zip_name;size:50";not null`
	WeightsName        string              `gorm:"column:weights_name;text"`
	ModelDescription   string              `gorm:"column:model_description;size:100";not null`
	Flag               bool                `gorm:"column:flag"`                 // 是否下发, 未下发为0，已下发为1
	EnableFlag         int                 `gorm:"column:enable_flag;not null"` // 0 启用， 1 root禁用， 2 admin禁用
	CreateTime         time.Time           `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime         time.Time           `gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Model2VideoConfigs []Model2VideoConfig `gorm:"foreignkey:ModelConfigID;references:ID;cascade:save"` // 在Gorm中处理一对多关系需要额外定义关联模型
	ModelMonitor       *ModelMonitor       `gorm:"foreignkey:ModelConfigID;references:ID;cascade:save"`
	User               string              `gorm:"column:user;size:32;ForeignKey:TUserInfo.ID;references:ID;on_delete:CASCADE"`

	// 这里省略了自定义方法如model_name_and_version和model_zip_file等
}

func (ModelConfig) TableName() string {
	return "t_model_config"
}

// Model2VideoConfig 结构体对应t_model2video_config表
type Model2VideoConfig struct {
	ID             string    `gorm:"primary_key;default:generate_uuid;size:32"`
	MineCode       string    `gorm:"column:mine_code;size:9";not null`
	ModelConfigID  string    `gorm:"column:model_id;size:32;ForeignKey:ModelConfig.ID;references:ID;on_delete:CASCADE"`
	VideoIndex     int       `gorm:"column:video_index;not null"`
	URL            string    `gorm:"column:url;size:200";not null`
	VideoDescribe  string    `gorm:"column:video_describe;size:200"`
	JSONPort       *int      `gorm:"column:json_port"`
	ImgPort        *int      `gorm:"column:img_port"`
	BusinessID     string    `gorm:"column:business_id;size:32"`
	VideoServiceIP string    `gorm:"column:video_service_ip;size:50"`
	JSONPortAppend string    `gorm:"column:json_port_append;size:20"`
	ImgPortAppend  string    `gorm:"column:img_port_append;size:20"`
	BusinessParams []byte    `gorm:"column:business_params;text"`
	CreateTime     time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime     time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// 在Gorm中处理一对一关系需要额外定义关联模型
	Model2VideoMonitor Model2VideoMonitor `gorm:"foreignkey:Model2VideoID"`
}

func (Model2VideoConfig) TableName() string {
	return "t_model2video_config"
}

// ModelMonitor 结构体对应t_model_monitor表
type ModelMonitor struct {
	ID            string    `gorm:"primary_key;default:generate_uuid;size:32"`
	MineCode      string    `gorm:"column:mine_code;size:9";not null`
	EndID         string    `gorm:"column:end_id;size:32;ForeignKey:EndInfo.ID;references:ID"`
	ModelConfigID string    `gorm:"column:model_id;size:36;ForeignKey:ModelConfig.ID;references:ID;on_delete:CASCADE"`
	MonitorStatus int       `gorm:"column:monitor_status;not null"` // 0代表停止/1代表运行中
	IsAbnormal    int       `gorm:"column:is_abnormal;not_null"`    // 1 正常 -1异常
	CreateTime    time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime    time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	// 在Gorm中处理一对多关系需要额外定义关联模型
	//Model2VideoMonitors []Model2VideoMonitor `gorm:"foreignkey:ModelID"`
}

func (ModelMonitor) TableName() string {
	return "t_model_monitor"
}

// Model2VideoMonitor 结构体对应t_video_monitor表
type Model2VideoMonitor struct {
	ID            string             `gorm:"primary_key;default:generate_uuid;size:32"`
	MineCode      string             `gorm:"column:mine_code;size:9";not null`
	EndID         string             `gorm:"column:end_id;size:32;ForeignKey:TEndInfo.ID;references:ID"`
	ModelConfigID string             `gorm:"column:model_id;size:36;ForeignKey:ModelConfig.ID;references:ID"`
	Model2VideoID string             `gorm:"column:model2video_id;size:36;ForeignKey:Model2VideoConfig.ID;references:ID;on_delete:CASCADE"`
	MonitorStatus int                `gorm:"column:monitor_status;not null"` // 0代表停止/1代表运行中
	IsAbnormal    int                `gorm:"column:is_abnormal;not_null"`    // 1正常 -1异常
	StartTime     *time.Time         `gorm:"column:start_time"`
	CreateTime    time.Time          `gorm:"column:create_time;default:CURRENT_TIMESTAMP"`
	UpdateTime    time.Time          `gorm:"column:update_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Model2Video   *Model2VideoConfig `gorm:"foreignkey:Model2VideoID"`
}

func (Model2VideoMonitor) TableName() string {
	return "t_video_monitor"
}

type AutoUpdateInfo struct {
	global.CmbpModel
	UserID            string  `gorm:"column:user_id;size:32;not null;comment:'用户id'"`
	ModelID           string  `gorm:"column:model_id;size:32;not null;comment:'模型id'"`
	ModelUpdateStatus int     `gorm:"column:model_update_status;not null;comment:'模型更新状态'"`
	Log               *string `gorm:"column:log;size:255"`
}

func (AutoUpdateInfo) TableName() string {
	return "t_auto_update_info" //模型自动升级异常表
}

type RuntimeModels struct {
	global.CmbpModel
	RuntimeID  string     `gorm:"column:runtime_id;size:32;index;ForeignKey:RuntimeAll"`
	ModelAllID string     `gorm:"column:model_all_id;size:32;index;ForeignKey:ModelAll"`
	Runtime    RuntimeAll `gorm:"foreignkey:RuntimeID"`
	Models     ModelAll   `gorm:"foreignkey:ModelAllID"`
}

func (RuntimeModels) TableName() string {
	return "t_runtime_models"
}

type RuntimeAll struct {
	global.CmbpModel
	Name        string `gorm:"column:name;not null;size:50"`
	ChineseName string `gorm:"column:chinese_name;not_null;size:50"`
	Description string `gorm:"column:description;not_null;size:100"`
	Tag         string `gorm:"column:tag;not_null;size:20"`
	Developer   string `gorm:"column:developer;not_null;size:20"`
	Type        int    `gorm:"column:type"`      // 默认空为我们的CV镜像，1为数据分析模型镜像
	Status      int    `gorm:"column:status"`    // 运行镜像上传状态
	BuildWay    int    `gorm:"column:build_way"` // 运行镜像上传方式
}

func (r *RuntimeAll) TableName() string {
	return "t_runtime_all"
}

func (r *RuntimeAll) ModelZipFilePath() string {
	return filepath.Join(global.CMBP_CONFIG.CMBPBase.OssRuntimeLibrary, r.Name+r.Tag+".zip")
}

type AutoBuildTask struct {
	global.CmbpModel
	UserID        string  `gorm:"column:user_id;not null;size:32"`
	Type          int     `gorm:"column:type;default:1"`
	BuildTaskName string  `gorm:"column:build_task_name;not null;size:50"`
	ModelAllID    *string `gorm:"column:model_all_id;size:32"`
	Desc          *string `gorm:"column:desc;size:200"`
	LabelData     string  `gorm:"column:label_data;type:text"` // JSON数据通常以字符串存储
	ScenseID      string  `gorm:"column:scense_id;not null;size:32"`
}

func (AutoBuildTask) TableName() string {
	return "t_auto_build_task"
}

type AutoUpdateTask struct {
	global.CmbpModel
	UserID        string    `gorm:"column:user_id;not null;size:32;comment:'用户id'"`
	UpdateStatus  int       `gorm:"column:update_status;not null;comment:'更新状态：1为更新完成，2为更新中'"`
	TotalCount    int       `gorm:"column:total_count;not null;comment:'更新模型数量'"`
	SuccessCount  int       `gorm:"column:success_count;not null;comment:'更新成功数量'"`
	FailedCount   int       `gorm:"column:failed_count;not null;comment:'更新失败数量'"`
	TaskStartTime time.Time `gorm:"column:task_start_time"`
}

func (AutoUpdateTask) TableName() string {
	return "t_auto_update_task"
}

type HardwareArch struct {
	global.CmbpModel
	Code              int    `gorm:"column:code"`
	Name              string `gorm:"column:name;size:200"`
	Desc              string `gorm:"column:desc;size:500"`
	RealChannelNumber int    `gorm:"column:real_channel_number"`
}

func (HardwareArch) TableName() string {
	return "t_hardware_arch"
}

type AIModelAll struct {
	global.CmbpModel
	ModelName        string `gorm:"column:model_name;not null;size:50"`
	ModelDescription string `gorm:"column:model_description;not null;size:100"`
	HardwareType     int    `gorm:"column:hardware_type"`
	AlgorithmID      int    `gorm:"column:algorithm_id"`
	IsImage          bool   `gorm:"column:is_image"`
	NeedGPU          bool   `gorm:"column:need_gpu"`
	IsRealChannel    bool   `gorm:"column:is_real_channel"`
	Cmd              string `gorm:"column:cmd;size:500"`
	AIModelAPI       string `gorm:"column:ai_model_api;text"`
	AIModelPurpose   int    `gorm:"column:ai_model_purpose"`
	Env              string `gorm:"column:env;size:500"`  // 注意：此字段若为JSON格式，可能需要自定义JSON解析逻辑或转换为结构体
	Vol              string `gorm:"column:vol;size:500"`  // 同上，需考虑JSON格式处理
	Port             string `gorm:"column:port;size:500"` // 同上
	Extra            string `gorm:"column:extra;size:500"`
	StartCmd         string `gorm:"column:start_cmd;size:500"`
	JsonURL          string `gorm:"column:json_url;size:200"`
	ImgURL           string `gorm:"column:img_url;size:200"`
	Developer        string `gorm:"column:developer;size:20"`
	BuildFlag        int    `gorm:"column:build_flag"`
	RuntimeID        string `gorm:"column:runtime_id;size:32"`

	// 下面这个属性不是直接映射数据库字段，而是提供一个基于model_name的计算属性
	ModelZipFile func() string `gorm:"-"`
}

func (m *AIModelAll) TableName() string {
	return "t_ai_model_all"
}

func (m *AIModelAll) ModelZipFilePath() string {
	return "/home/models/ModelLibrary/" + m.ModelName
}

type BusinessModelAll struct {
	global.CmbpModel
	ModelName        string `gorm:"column:model_name;size:50"`
	ModelDescription string `gorm:"column:model_description;size:100"`
	NeedAI           bool   `gorm:"column:need_ai"`
	BusinessParams   string `gorm:"column:business_params;text"`
	JsonOutput       string `gorm:"column:json_output;text"`
	User             string `gorm:"column:user;size:32"`
	BusinessAPI      string `gorm:"column:business_api;text"`
	BusinessPurpose  int    `gorm:"column:business_purpose"`
	BusinessType     string `gorm:"column:business_type;text"`
}

func (BusinessModelAll) TableName() string {
	return "t_business_model_all"
}

// ModelAllBusinessModel 表示模型-业务模型关联表
type ModelAllBusinessModel struct {
	global.CmbpModel
	ModelAllID      string            `gorm:"type:char(32);not null" json:"model_all_id"`
	BusinessModelID string            `gorm:"type:char(32);not null" json:"business_model_id"`
	Model           *ModelAll         `gorm:"foreignKey:ModelAllID" json:"-"`
	BusinessModel   *BusinessModelAll `gorm:"foreignKey:BusinessModelID" json:"-"`
}

// TableName 返回表名
func (ModelAllBusinessModel) TableName() string {
	return "t_model_all_business_model"
}

type Notebook struct {
	global.CmbpModel
	UserID         string `gorm:"not null;size:32"`
	UUID           string `gorm:"not null;size:32"`
	Type           int    `gorm:"nullable"`
	Name           string `gorm:"not null;size:32"`
	Desc           string `gorm:"size:255;default:null"`
	URL            string `gorm:"not null"`
	Status         int    `gorm:"nullable"`
	ExpirationTime *int   `gorm:"nullable"`
}

func (Notebook) TableName() string {
	return "t_notebook"
}

type ApplicationRecord struct {
	global.CmbpModel
	ModelID           string `gorm:"column:model_id;size:32"`
	ApplicationStatus int    `gorm:"column:application_status"` // 99 免测申请中，100 免测申请拒绝
	Reason            string `gorm:"size:500"`
	User              string `gorm:"size:32"`
	ProcessType       int
}

func (ApplicationRecord) TableName() string {
	return "t_application_record"
}

type DataAnalyzeModelAll struct {
	global.CmbpModel
	ModelName        string `gorm:"not null;size:50" json:"model_name"`
	ModelChineseName string `gorm:"size:50" json:"model_chinese_name"`
	ModelDescription string `gorm:"not null;size:100" json:"model_description"`
	Env              string `gorm:"default:null;size:500" json:"env"`
	Vol              string `gorm:"default:null;size:500" json:"vol"`
	Port             string `gorm:"default:null;size:500" json:"port"`
	Extra            string `gorm:"default:null;size:500" json:"extra"`
	Cmd              string `gorm:"not null;size:500" json:"cmd"`
	AIModelAPI       string `gorm:"type:text;nullable" json:"ai_model_api"`
	AIModelPurpose   *int   `gorm:"column:ai_model_purpose;nullable" json:"ai_model_purpose"`
	Params           string `gorm:"default:null;type:text" json:"params"`
	User             string `gorm:"size:32" json:"user"`
	ConfigJSON       string `gorm:"default:null;type:text" json:"config_json"`
	// 自定义方法
	//ModelZipFile func() string `json:"-"`
}

func (d *DataAnalyzeModelAll) ModelZipFile() string {
	return filepath.Join("/home/models/DataAnalyzeModelLibrary", d.ModelName)
}

func (d *DataAnalyzeModelAll) TableName() string {
	return "t_data_analyze_model_all"
}

type WeightsManagement struct {
	global.CmbpModel
	WeightsName    string `gorm:"not null;size:255" json:"weights_name" doc:"权重名称"`
	WeightsDesc    string `gorm:"size:500" json:"weights_desc" doc:"权重描述"`
	WeightsCreator string `gorm:"not null;size:32" json:"weights_creator" doc:"权重创建者"`
	UserName       string `gorm:"size:32" json:"user_name" doc:"创建者名称"`
	WeightsFile    string `gorm:"not null;size:32" json:"weights_file" doc:"文件名"`
	WeightsType    int    `gorm:"not null" json:"weights_type" doc:"权重类型"`
}

func (WeightsManagement) TableName() string {
	return "t_weights_management"
}

type TestFreeModelRes struct {
	ModelID              string `json:"model_id"`
	ModelKind            int    `json:"-" gorm:"column:model_kind"`
	Edition              int    `json:"-" gorm:"column:edition"`
	ModelTypeDesc        string `json:"model_type_desc"`
	ModelFieldDesc       string `json:"model_field_desc"`
	ModelName            string `json:"model_name"`
	ModelChineseName     string `json:"model_chinese_name"`
	ModelVersion         string `json:"model_version"`
	ModelDescription     string `json:"model_description"`
	TechnicalDescription string `json:"technical_description"`
	PerformanceDesc      string `json:"performance_desc"`
	HardwareTypeName     string `json:"hardware_type_name"`
	IsImage              string `json:"is_image"`
	Cmd                  string `json:"cmd"`
	JsonURL              string `json:"json_url"`
	ImgURL               string `json:"img_url"`
	OnBoot               string `json:"on_boot"`
	NeedGPU              string `json:"need_gpu"`
	AuditState           string `json:"audit_state"`
	User                 string `json:"user"`
	Developer            string `json:"developer"`
	UploadTime           string `json:"upload_time"`
	ImgPath              string `json:"img_path" gorm:"-"`
	Img2Path             string `json:"img2_path" gorm:"-"`
	VideoPath            string `json:"video_path" gorm:"-"`
	Edit                 string `json:"edit"`
	TestStatus           int    `json:"test_status"`
	Reason               string `json:"reason"`
	Phone                string `json:"phone"`
	ApplicationTime      string `json:"application_time"`
	ProcessType          int    `json:"process_type"`
}

func (m *TestFreeModelRes) AfterFind(tx *gorm.DB) (err error) {
	// 图片路径生成
	if m.ModelKind != 1 {
		baseName := m.ModelName + "V" + m.ModelVersion
		m.ImgPath = filepath.Join(global.CMBP_CONFIG.CMBPBase.OssPath, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia, baseName+".jpg")
		m.VideoPath = filepath.Join(global.CMBP_CONFIG.CMBPBase.OssPath, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia, baseName+".mp4")
	} else {
		versionInt := m.ModelVersion
		formattedName := fmt.Sprintf("%sV%d.%s", m.ModelName, versionInt, m.Edition)
		m.ImgPath = filepath.Join(global.CMBP_CONFIG.CMBPBase.OssPath, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia, formattedName+".jpg")
		m.VideoPath = filepath.Join(global.CMBP_CONFIG.CMBPBase.OssPath, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia, formattedName+".mp4")
	}
	return nil
}

// DataModelConfig 对应 t_data_model_config 表
type DataModelConfig struct {
	MineCode      *string `gorm:"column:mine_code;type:varchar(9);default:null" json:"mine_code,omitempty"`
	ModelInfoID   *string `gorm:"column:model_info_id;type:varchar(32);default:null" json:"model_info_id,omitempty"`
	ContainerName *string `gorm:"column:container_name;type:varchar(50);default:null" json:"container_name,omitempty"`
	Desc          *string `gorm:"column:desc;type:varchar(32);default:null" json:"desc,omitempty"`        // 实例说明
	Env           *string `gorm:"column:env;type:varchar(500);default:null" json:"env,omitempty"`         // 环境变量
	Vol           *string `gorm:"column:vol;type:varchar(500);default:null" json:"vol,omitempty"`         // 挂载目录
	Port          *string `gorm:"column:port;type:varchar(500);default:null" json:"port,omitempty"`       // 端口
	Extra         *string `gorm:"column:extra;type:varchar(500);default:null" json:"extra,omitempty"`     // 其他
	Cmd           string  `gorm:"column:cmd;type:varchar(500);not null" json:"cmd"`                       // 命令
	Params        *string `gorm:"column:params;type:varchar(500);default:null" json:"params,omitempty"`   // 运行参数
	DeployStatus  int     `gorm:"column:deploy_status;type:int" json:"deploy_status"`                     // 部署状态
	ErrMsg        *string `gorm:"column:err_msg;type:varchar(500);default:null" json:"err_msg,omitempty"` // 错误信息
	User          *string `gorm:"column:user;type:varchar(32);default:null" json:"user,omitempty"`        // 部署模型用户
}

// TableName 设置表名
func (d *DataModelConfig) TableName() string {
	return "t_data_model_config"
}

// ModelFeedback 对应 t_model_feedback 表
type ModelFeedback struct {
	global.CmbpModel
	MineCode         *string    `gorm:"column:mine_code;type:varchar(9);default:null" json:"mine_code,omitempty"`
	EventID          *string    `gorm:"column:event_id;type:varchar(32);default:null;index" json:"event_id,omitempty"` // 外键关联 t_event.id
	ModelID          *string    `gorm:"column:model_id;type:varchar(32);default:null;index" json:"model_id,omitempty"` // 外键关联 t_model_all.id
	OriginalVideoURL *string    `gorm:"column:original_video_url;type:text;default:null" json:"original_video_url,omitempty"`
	ErrorVideoURL    *string    `gorm:"column:error_video_url;type:varchar(200);default:null" json:"error_video_url,omitempty"`
	VideoDescribe    *string    `gorm:"column:video_describe;type:varchar(200);default:null" json:"video_describe,omitempty"`
	Title            *string    `gorm:"column:title;type:varchar(200);default:null" json:"title,omitempty"`
	Flag             int        `gorm:"column:flag;type:int" json:"flag"` // 上传状态,0代表开始上传，1代表上传成功，-1代表上传失败
	StartTime        *time.Time `gorm:"column:start_time;type:datetime;default:null" json:"start_time,omitempty"`
	EndTime          *time.Time `gorm:"column:end_time;type:datetime;default:null" json:"end_time,omitempty"`
	Base64Image      *string    `gorm:"column:base64_image;type:text;default:null" json:"base64_image,omitempty"`
	Advice           *string    `gorm:"column:advice;type:varchar(500);default:null" json:"advice,omitempty"`
	Type             *int       `gorm:"column:type;type:int;default:null" json:"type,omitempty"`                   // 没有 代表1.3版本反馈，1代表探水系统，2代表客户端
	ConfigID         *string    `gorm:"column:config_id;type:varchar(32);default:null" json:"config_id,omitempty"` // 实例id
}

// TableName 设置表名
func (ModelFeedback) TableName() string {
	return "t_model_feedback"
}

// ModelDeploy 对应 t_model_deploy 表
type ModelDeploy struct {
	global.CmbpModel
	MineCode *string `gorm:"column:mine_code;type:varchar(9);default:null" json:"mine_code,omitempty"`
	EventID  *string `gorm:"column:event_id;type:varchar(32);default:null;index" json:"event_id,omitempty"` // 外键关联 t_event.id
	ModelID  *string `gorm:"column:model_id;type:varchar(32);default:null;index" json:"model_id,omitempty"` // 外键关联 t_model_all.id
	Title    *string `gorm:"column:title;type:varchar(200);default:null" json:"title,omitempty"`
	Type     int     `gorm:"column:type;type:int" json:"type"` // 类型 2:代表模型下发、3:代表模型更新
}

// TableName 设置表名
func (ModelDeploy) TableName() string {
	return "t_model_deploy"
}

// ModelEvaluation 对应 t_model_evaluation 表
type ModelEvaluation struct {
	global.CmbpModel
	MineCode         *string `gorm:"column:mine_code;type:varchar(9);default:null" json:"mine_code,omitempty"`
	EventID          *string `gorm:"column:event_id;type:varchar(32);default:null;index" json:"event_id,omitempty"` // 外键关联 t_event.id
	ModelID          *string `gorm:"column:model_id;type:varchar(32);default:null;index" json:"model_id,omitempty"` // 外键关联 t_model_all.id
	Title            string  `gorm:"column:title;type:varchar(200);default:'模型评价'" json:"title"`
	Entirety         int     `gorm:"column:entirety;type:int;default:0" json:"entirety"`                                 // 整体, default 为 0
	Accuracy         int     `gorm:"column:accuracy;type:int;default:0" json:"accuracy"`                                 // 准确性, default 为 0
	Stability        int     `gorm:"column:stability;type:int;default:0" json:"stability"`                               // 稳定性, default 为 0
	RecognitionModel *int    `gorm:"column:recognition_model;type:int;default:null" json:"recognition_model,omitempty"`  // 识别模型, 1 赞 -1 踩 0 不赞不踩
	BusinessList     *string `gorm:"column:business_list;type:varchar(500);default:null" json:"business_list,omitempty"` // 业务模型, 1 赞 -1 踩 0 不赞不踩 {"空载":1,停机:-1}
	Advice           *string `gorm:"column:advice;type:varchar(500);default:null" json:"advice,omitempty"`
}

// TableName 设置表名
func (ModelEvaluation) TableName() string {
	return "t_model_evaluation"
}

// Event 对应 t_event 表
type Event struct {
	global.CmbpModel
	MineCode        *string            `gorm:"column:mine_code;type:varchar(9);default:null" json:"mine_code,omitempty"`
	EventType       int                `gorm:"column:event_type;type:int" json:"event_type"`                              // 0:代表漏识别、1:代表误识别、2:代表模型下发、3:代表模型更新 4 模型评价 5 在线发布
	EventStatus     int                `gorm:"column:event_status;type:int" json:"event_status"`                          // 0:代表事件未读|未响应、1:代表事件已读|已响应、2:代表事件正在处理、3:代表事件处理解决 99 已关闭 -1事件撤销
	Creator         *string            `gorm:"column:creator;type:varchar(50);default:null" json:"creator,omitempty"`     // 代表由哪个用户创建该事件
	Responder       *string            `gorm:"column:responder;type:varchar(50);default:null" json:"responder,omitempty"` // 代表由哪个用户响应该事件
	Describe        *string            `gorm:"column:describe;type:varchar(500);default:null" json:"describe,omitempty"`  // 关于事件的一些描述
	ModelFeedback   []ModelFeedback    `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ModelEvaluation []ModelEvaluation  `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ModelDeploy     []ModelDeploy      `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	OnlinePublish   []OnlinePublishing `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// TableName 设置表名
func (Event) TableName() string {
	return "t_event"
}

// OnlinePublishing 对应 t_online_publishing 表
type OnlinePublishing struct {
	EventID             *string `gorm:"column:event_id;type:varchar(32);default:null;index" json:"event_id,omitempty"` // 外键关联 t_event.id
	TrainingName        *string `gorm:"column:training_name;type:varchar(100);default:null" json:"training_name,omitempty"`
	TrainingDescription *string `gorm:"column:training_description;type:varchar(200);default:null" json:"training_description,omitempty"`
	TrainingID          *string `gorm:"column:training_id;type:varchar(200);default:null" json:"training_id,omitempty"`
	Path                *string `gorm:"column:path;type:varchar(200);default:null" json:"path,omitempty"`                    // 路径
	TrainingStaff       *string `gorm:"column:training_staff;type:varchar(50);default:null" json:"training_staff,omitempty"` // 训练人员/共享人员
	Developer           *string `gorm:"column:developer;type:varchar(50);default:null" json:"developer,omitempty"`
	Title               *string `gorm:"column:title;type:varchar(200);default:null" json:"title,omitempty"`
}

// TableName 设置表名
func (OnlinePublishing) TableName() string {
	return "t_online_publishing"
}

// ClassRelatePangu 表示模型关联盘古大模型
type ClassRelatePangu struct {
	global.CmbpModel
	ModelID      string `gorm:"type:varchar(32);comment:'模型id'" json:"model_id"`
	PanguModelID string `gorm:"type:varchar(100);comment:'盘古大模型id'" json:"pangu_model_id"`
	ClassName    string `gorm:"type:text;comment:'分析类型'" json:"class_name"`
}

// TableName 返回表名
func (ClassRelatePangu) TableName() string {
	return "t_class_relate_pangu"
}

type ModelUpdateRecord struct {
	global.CmbpModel
	UserID           string `gorm:"column:user_id;type:varchar(32)" json:"user_id" doc:"用户ID"`
	ModelID          string `gorm:"column:model_id;type:varchar(32)" json:"model_id" doc:"模型ID"`
	ModelChineseName string `gorm:"column:model_chinese_name;type:varchar(255)" json:"model_chinese_name" doc:"模型名称"`
	Status           int    `gorm:"column:status type:int" json:"statu" doc:"状态 0 未更新 1 已更新"`
}

func (ModelUpdateRecord) TableName() string {
	return "t_model_update_record"
}
