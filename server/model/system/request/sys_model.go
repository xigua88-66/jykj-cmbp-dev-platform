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
