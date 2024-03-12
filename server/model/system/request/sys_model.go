package request

type ModelFiled struct {
	Page         *int   `form:"page"`
	Limit        *int   `form:"limit"`
	IndustryCode *int   `form:"industry_code"`
	Flag         string `form:"flag"`
	Code         string `form:"code"`
	Keywords     string `form:"key_code"`
	ModelPurpose int    `form:"model_purpose"`
}

type AlgorithmRqe struct {
	Page        *int   `form:"page"`
	Limit       *int   `form:"limit"`
	ModelKind   *int   `form:"model_kind"`
	AlgorithmID *int   `form:"algorithm_id"`
	KeyWords    string `form:"key_words"`
}

type ModelListReq struct {
	Limit        int    `form:"limit"`
	Page         int    `form:"page"`
	IndustryCode *int   `form:"industry_code"`
	NameOrDesc   string `form:"name_or_desc"`
	Code         string `form:"code"`
	Type         string `form:"type"`
	Flag         int    `form:"flag"`
	AlgorithmID  *int   `form:"algorithm_id"`
	IsCloud      int    `form:"is_cloud"` // 注意：根据实际业务判断是否应为bool类型
	ReqType      int    `form:"req_type"`
	ModelPurpose *int   `form:"model_purpose"`
	ModelKind    *int   `form:"model_kind"`
}

type ModelStoreRqe struct {
	Page         int    `form:"page" binding:"required"`
	Limit        int    `form:"limit" binding:"required"`
	IndustryCode string `form:"industry_code"`
	NameOrDesc   string `form:"name_or_desc"`
	Flag         int    `form:"flag"`
	Mine         *int   `form:"mine"`
	TestStatus   string `form:"test_status"`
	ModelKind    *int   `form:"model_kind"`
	AlgorithmID  *int   `form:"algorithm_id"`
	UpdateStatus int    `form:"update_status"`
}

type UploadModelStoreReq struct {
	ModelName              string  `form:"model_name"`
	UUID                   string  `form:"uuid"`
	ModelChineseName       string  `form:"model_chinese_name"`
	ModelVersion           string  `form:"model_version" `
	ModelDescription       string  `form:"model_description"`
	TechnicalDescription   string  `form:"technical_description"`
	PerformanceDescription string  `form:"performance_description"`
	ModelType              string  `form:"model_type" `
	HardwareType           string  // 在Golang中硬件类型一般是枚举，但这里假设它是字符串
	IsImage                *int    `form:"is_image" `
	Cmd                    string  `form:"cmd" `
	JsonUrl                string  `form:"json_url"`
	ImgUrl                 string  `form:"img_url" `
	BusinessDict           string  `form:"business_dict"`
	OnBoot                 *int    `form:"on_boot"`
	NeedGPU                *int    `form:"need_gpu" `
	TestDuration           float64 `form:"test_duration"`
	Accuracy               int     `form:"accuracy"`
	IsRealChannel          int     `form:"is_real_channel"`
	TaskID                 string  `form:"task_id" `
	AlgorithmID            int     `form:"algorithm_id" `
	ModelPurpose           int     `form:"model_purpose;default=1"`
	AIModelAPI             string  `form:"ai_model_api;default=[]"`
	Advantage              string  `form:"advantage;default=''"`

	// 文件字段在Golang中通常不会直接包含在结构体中，而是通过请求处理函数中的FormFile方法获取
	//VideoFile multipart.File `form:"video_file"`
	//ModelImg  multipart.File `form:"model_img"`
}

type GetModelKind struct {
	Page         *int   `form:"page"`
	Limit        *int   `form:"limit"`
	ModelKind    int    `from:"model_kind"`
	Keywords     string `form:"keywords"`
	ModelPurpose int    `form:"model_purpose"`
}

type GetHardWare struct {
	Code string `form:"code"`
	Flag *int   `form:"flag"`
}

type GetIndustry struct {
	Page         *int   `form:"page"`
	Limit        *int   `form:"limit"`
	IndustryCode *int   `form:"industry_code"`
	Keywords     string `form:"keywords"`
	ModelPurpose int    `form:"model_purpose"`
}

type GetTestFreeApply struct {
	Page       *int   `form:"page"`
	Limit      *int   `form:"limit"`
	NameOrDesc string `form:"name_or_desc"`
}
