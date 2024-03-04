package response

type GetModelField struct {
	FieldName    string `json:"field_name"`
	FieldNameEn  string `json:"field_name_en"`
	FieldImgUrl  string `json:"field_img_url"`
	FieldCode    string `json:"field_code"`
	IndustryName string `json:"industry_name"`
	IndustryCode int64  `json:"industry_code"`
}

type ModelFieldRsp struct {
	Count          int             `json:"count"`
	ModelFieldList []GetModelField `json:"model_field_list"`
}

type GetModelType struct {
	SceneName    string `json:"scene_name"`
	SceneCode    string `json:"scene_code"`
	FieldName    string `json:"field_name"`
	FieldCode    string `json:"field_code"`
	IndustryName string `json:"industry_name"`
	IndustryCode int64  `json:"industry_code"`
}

type AlgorithmInfo struct {
	AlgorithmID     int64   `json:"algorithm_id"`
	AlgorithmName   string  `json:"algorithm_name"`
	AlgorithmNameEn *string `json:"algorithm_name_en"`
	AlgorithmDesc   string  `json:"algorithm_desc"`
	AlgorithmImg    string  `json:"algorithm_img"`
	AlgorithmVideo  string  `json:"algorithm_video"`
	IsDelete        int     `json:"is_delete"`
	CreateTime      string  `json:"create_time"`
}
