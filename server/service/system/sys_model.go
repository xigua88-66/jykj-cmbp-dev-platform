package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/system"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	systemRsp "jykj-cmbp-dev-platform/server/model/system/response"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ModelService struct {
}

func (modelService *ModelService) GetModelField(params systemReq.ModelFiled) (data interface{}, err error) {
	var fieldCode []string
	var typeCode []string

	if params.ModelPurpose == 2 {
		err = global.CMBP_DB.Model(&system.ModelField{}).
			Select("DISTINCT model_field").
			Where("model_purpose = ?", 2).
			Where("test_status IS NULL").
			Group("model_field").Pluck("model_field", &fieldCode).Error

		err = global.CMBP_DB.Model(&system.ModelField{}).
			Select("DISTINCT model_type").
			Where("model_purpose = ?", 2).
			Where("test_status IS NULL").
			Group("model_type").Pluck("model_type", &typeCode).Error
		if err != nil {
			return nil, err
		}
	}

	var QUERY = global.CMBP_DB.Model(&system.ModelField{})

	if params.Flag == "0" {
		var allFieldCode []system.ModelField
		if len(fieldCode) > 0 {
			QUERY.Where("code in ?", fieldCode)
		}
		if params.IndustryCode != nil {
			QUERY = QUERY.Where("industry_code = ?", *params.IndustryCode)
		}
		if params.Keywords != "" {
			QUERY.Where("name LIKE ? or field_name_en LIKE ?", "%"+params.Keywords+"%")
		}
		if params.Page == nil || params.Limit == nil {
			err = QUERY.Find(&allFieldCode).Error
			if err != nil {
				return nil, err
			}
		} else {
			err = QUERY.Limit(*params.Limit).Offset(*params.Limit * (*params.Page - 1)).Find(&allFieldCode).Error
			if err != nil {
				return nil, err
			}
		}
		var fieldRspList []systemRsp.GetModelField

		for _, fieldObj := range allFieldCode {
			var industry system.Industry

			err = global.CMBP_DB.Model(&system.Industry{}).
				Where("industry_code = ?", fieldObj.IndustryCode).
				Find(&industry).Error
			if err != nil {
				return nil, err
			}

			var fieldRsp systemRsp.GetModelField
			fieldRsp.FieldName = fieldObj.Name
			fieldRsp.FieldNameEn = fieldObj.FieldNameEn
			// todo url前缀从配置文件中获取
			fieldRsp.FieldImgUrl = fieldObj.FieldImgPath
			fieldRsp.FieldCode = fieldObj.Code
			fieldRsp.IndustryName = industry.IndustryName
			fieldRsp.IndustryCode = industry.IndustryCode
			fieldRspList = append(fieldRspList, fieldRsp)
		}
		if params.Page != nil && params.Limit != nil {
			return systemRsp.ModelFieldRsp{Count: len(fieldRspList), ModelFieldList: fieldRspList}, nil
		} else {
			return fieldRspList, nil
		}
	} else {
		if params.Code != "" {
			var codeList = strings.Split(params.Code, ",")
			sort.Strings(codeList)
			var typeRspList []systemRsp.GetModelType

			for _, code := range codeList {
				var modelField system.ModelField
				err = global.CMBP_DB.Model(system.ModelField{}).Where("code = ?", code).First(&modelField).Error
				if err != nil {
					return nil, err
				}

				var modelType []system.ModelType
				if typeCode != nil {
					err = global.CMBP_DB.Model(&system.ModelType{}).
						Where("model_field_id = ?", modelField.ID).
						Where("model_type in ?", typeCode).
						Find(&modelType).Error
					if err != nil {
						return nil, err
					}
				} else {
					err = global.CMBP_DB.Model(&system.ModelType{}).
						Where("model_field_id = ?", modelField.ID).
						Find(&modelType).Error
					if err != nil {
						return nil, err
					}
				}
				for _, typeObj := range modelType {
					var industry system.Industry
					global.CMBP_DB.Model(&system.Industry{}).
						Where("industry_code = ?", modelField.IndustryCode).
						First(&industry)
					var typeRsp systemRsp.GetModelType
					typeRsp.SceneName = typeObj.ModelTypeDesc
					typeRsp.SceneCode = typeObj.ModelType
					typeRsp.FieldName = typeObj.ModelField.Name
					typeRsp.FieldCode = typeObj.ModelField.Code
					typeRsp.IndustryName = industry.IndustryName
					typeRsp.IndustryCode = industry.IndustryCode
					typeRspList = append(typeRspList, typeRsp)
				}
			}
			sort.Slice(typeRspList, func(i int, j int) bool {
				//return typeRspList[i].FieldCode+typeRspList[i].SceneCode < typeRspList[j].FieldCode+typeRspList[j].SceneCode
				fieldCodeLeft := typeRspList[i].FieldCode
				sceneCodeLeft := typeRspList[i].SceneCode
				combinedKeyI := fieldCodeLeft + sceneCodeLeft
				fieldCodeRight := typeRspList[j].FieldCode
				sceneCodeRight := typeRspList[j].SceneCode
				combinedKeyJ := fieldCodeRight + sceneCodeRight
				return combinedKeyI < combinedKeyJ
			})
			return typeRspList, nil
		} else {
			if params.IndustryCode != nil {
				QUERY = QUERY.Where("industry_code = ?", params.IndustryCode).Order("code")
			} else {
				QUERY = QUERY.Order("code DESC")
			}
			var fieldIdList []string
			err = QUERY.Pluck("id", &fieldIdList).Error
			if err != nil {
				return nil, err
			}

			var modelType []system.ModelType
			if typeCode != nil {
				err = global.CMBP_DB.Preload("ModelField").Model(&system.ModelType{}).
					Where("model_type IN ?", typeCode).
					Where("model_field_id IN ?", fieldIdList).Find(&modelType).Error
				if err != nil {
					return nil, err
				}
			} else {
				err = global.CMBP_DB.Preload("ModelField").Model(&system.ModelType{}).Where("model_field_id IN ?", fieldIdList).Find(&modelType).Error
				if err != nil {
					return nil, err
				}
			}
			var typeRspList []systemRsp.GetModelType
			for _, typeObj := range modelType {
				var industry system.Industry
				err = global.CMBP_DB.Model(&system.Industry{}).
					Where("industry_code = ?", typeObj.ModelField.IndustryCode).
					First(&industry).Error
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				var typeRsp systemRsp.GetModelType
				typeRsp.SceneName = typeObj.ModelTypeDesc
				typeRsp.SceneCode = typeObj.ModelType
				typeRsp.FieldName = typeObj.ModelField.Name
				typeRsp.FieldCode = typeObj.ModelField.Code
				typeRsp.IndustryName = industry.IndustryName
				typeRsp.IndustryCode = industry.IndustryCode
				typeRspList = append(typeRspList, typeRsp)
			}
			sort.Slice(typeRspList, func(i int, j int) bool {
				//return typeRspList[i].FieldCode+typeRspList[i].SceneCode < typeRspList[j].FieldCode+typeRspList[j].SceneCode
				fieldCodeLeft := typeRspList[i].FieldCode
				sceneCodeLeft := typeRspList[i].SceneCode
				combinedKeyI := fieldCodeLeft + sceneCodeLeft
				fieldCodeRight := typeRspList[j].FieldCode
				sceneCodeRight := typeRspList[j].SceneCode
				combinedKeyJ := fieldCodeRight + sceneCodeRight
				return combinedKeyI < combinedKeyJ
			})
			return typeRspList, nil
		}
	}
}

func (modelService *ModelService) GetAlgorithmLogic(params systemReq.AlgorithmRqe) (rspList []systemRsp.AlgorithmInfo, err error) {
	var AlgorithmInfo []system.AlgorithmInfo
	QUERY := global.CMBP_DB.Model(&system.AlgorithmInfo{})
	if params.ModelKind != nil {
		QUERY = QUERY.Where("model_kind = ?", params.ModelKind)
	}
	if params.KeyWords != "" {
		QUERY = QUERY.Where("algorithm_name LIKE ? OR algorithm_name_en LIKE ? OR algorithm_desc LIKE ?", "%"+params.KeyWords+"%")
	}
	err = QUERY.Find(&AlgorithmInfo).Error
	if err != nil {
		return nil, err
	}
	var algorithmList []systemRsp.AlgorithmInfo
	for _, algorithmInfo := range AlgorithmInfo {
		var modelAll system.ModelAll
		var isDeleted = 0
		err := global.CMBP_DB.Where("algorithm_id = ?", algorithmInfo.AlgorithmID).First(&modelAll).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			isDeleted = 1
		}
		var imgPath = ""
		var videoPath = ""
		if algorithmInfo.AlgorithmImgPath != "" {
			imgPath = algorithmInfo.AlgorithmImgPath
		}
		if algorithmInfo.AlgorithmVideoPath != "" {
			videoPath = algorithmInfo.AlgorithmVideoPath
		}
		var algorithmRsp systemRsp.AlgorithmInfo
		algorithmRsp.AlgorithmID = algorithmInfo.AlgorithmID
		algorithmRsp.AlgorithmName = algorithmInfo.AlgorithmName
		algorithmRsp.AlgorithmNameEn = algorithmInfo.AlgorithmNameEN
		algorithmRsp.AlgorithmDesc = algorithmInfo.AlgorithmDesc
		algorithmRsp.AlgorithmImg = imgPath
		algorithmRsp.AlgorithmVideo = videoPath
		algorithmRsp.IsDelete = isDeleted
		algorithmRsp.CreateTime = algorithmInfo.CreateTime.Format("2006-01-02 15:04:05")
		algorithmList = append(algorithmList, algorithmRsp)
	}
	return algorithmList, nil
}

func (modelService *ModelService) GetModelList(params systemReq.ModelListReq, userId string) (getModelList interface{}, err error) {
	var hadModels []string
	global.CMBP_DB.Model(&system.AssetsManagement{}).Where("user_id = ? AND assets_type_id = 1", userId).Pluck("assets_id", &hadModels)
	var appliedModels []string
	global.CMBP_DB.Model(&system.AssetsRecord{}).Where("user_id = ? AND apply_status = 0", userId).Pluck("assets_id", &appliedModels)

	var user system.Users
	global.CMBP_DB.Where("id = ?", userId).First(&user)

	QUERY := global.CMBP_DB.Model(&system.ModelMarketList{})
	if params.IndustryCode != nil {
		QUERY = QUERY.Where("industry_code = ?", params.IndustryCode)
	} else if params.ModelPurpose != nil {
		QUERY = QUERY.Where("model_purpose = ?", params.ModelPurpose)
	} else if params.NameOrDesc != "" {
		QUERY = QUERY.Where("model_name LIKE ? OR model_chinese_name LIKE ? OR model_description LIKE ? OR technical_description LIKE ? OR performance_description LIKE ?", "%"+params.NameOrDesc+"%")
		// TODO 添加用户访问记录
	} else if params.Type != "" {
		QUERY = QUERY.Where("model_type = ?", params.Type)
	} else if params.Code != "" {
		QUERY = QUERY.Where("model_field = ?", params.Code)
	} else if params.ReqType == 1 {
		hadModels = append(hadModels, appliedModels...)
		// TODO 对hadModels去重
		QUERY = QUERY.Where("id NOT IN ?", hadModels)
	} else if params.ReqType == 2 {
		QUERY = QUERY.Where("id IN ?", hadModels)
	} else if params.ReqType == 3 {
		var downloadModels []string
		if user.MineCode != "" {
			global.CMBP_DB.Model(&system.Model{}).Where("mine_code = ?", user.MineCode).Pluck("model_all_id", &downloadModels)
		}
		QUERY = QUERY.Where("id IN ?", hadModels).Where("model_kind != 1 AND id NOT IN ?", downloadModels)
	} else if params.ReqType == 4 {
		QUERY = QUERY.Where("id NOT IN ?", hadModels).Where("id IN ?", appliedModels)
	} else {
		return nil, errors.New("参数错误")
	}
	if params.Flag == 1 {
		QUERY = QUERY.Order("down_load_count DESC")
	} else {
		QUERY = QUERY.Order("update_time DESC")
	}
	var modelList []system.ModelMarketList
	if params.AlgorithmID != nil {
		QUERY = QUERY.Where("algorithm_id = ?", params.AlgorithmID)
		QUERY.Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&modelList)
		dataList := FormatModelList(modelList, *params.ModelPurpose, hadModels, appliedModels, user.MineCode, userId)
		rspList := map[string]interface{}{
			"model_list": dataList,
			"count":      len(dataList),
		}
		return rspList, nil
	} else if params.ModelKind != nil {
		QUERY = QUERY.Where("model_kind = ?", params.ModelKind)
		err = QUERY.Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&modelList).Error
		if err != nil {
			return nil, err
		}
		dataList := FormatModelList(modelList, *params.ModelPurpose, hadModels, appliedModels, user.MineCode, userId)
		rspList := map[string]interface{}{
			"model_list": dataList,
			"count":      len(dataList),
		}
		return rspList, nil
	} else {
		var modelKind []system.ModelKind
		kindQuery := global.CMBP_DB.Model(&system.ModelKind{}).Order("model_kind")
		if params.NameOrDesc != "" {
			kindQuery = kindQuery.Where("kind_name LIKE ? OR kind_name_en LIKE ? OR kind_desc LIKE ?", params.NameOrDesc)
		}
		err = kindQuery.Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&modelKind).Error
		if err != nil {
			return nil, err
		}
		for _, kindInfo := range modelKind {
			var limit int
			if kindInfo.ModelKind == 2 {
				limit = 10
			} else {
				limit = params.Limit
			}
			QUERY = QUERY.Where("model_kind = ?", kindInfo.ModelKind)
			err = QUERY.Limit(limit).Offset(limit * (params.Page - 1)).Find(&modelList).Error
			dataList := FormatModelList(modelList, *params.ModelPurpose, hadModels, appliedModels, user.MineCode, userId)
			rspList := map[string]interface{}{
				"model_kind":  kindInfo.ModelKind,
				"count":       len(dataList),
				"total_count": len(dataList),
				"model_list":  dataList,
			}
			return rspList, nil
		}
	}

	return nil, err
}

func FormatModelList(modelList []system.ModelMarketList, isCloud int, hadModel []string, appliedModels []string, mineCode string, userID string) []interface{} {

	var rspData []interface{}

	for _, m := range modelList {
		modelZhName := m.ModelChineseName
		if isCloud != 0 {
			nameSplit := strings.Split(modelZhName, "-")
			if len(nameSplit) >= 2 {
				modelZhName = "某企业" + strings.Join(nameSplit[1:], "")
			}
		}
		// TODO 通过配置文件返回生成地址
		var isHad int
		if strings.Contains(m.ID, strings.Join(hadModel, "")) {
			isHad = 1
		} else {
			if strings.Contains(m.ID, strings.Join(appliedModels, "")) {
				isHad = 2
			} else {
				isHad = 0
			}
		}
		var isDownload int
		var adminUser system.Users
		global.CMBP_DB.Where("mine_code = ?", mineCode).Order("create_time DESC").First(&adminUser)

		var user *string
		if userID == adminUser.ID {
			user = nil
		} else {
			user = &userID
		}
		var modelInfo system.Model
		if m.ModelKind != 1 {
			global.CMBP_DB.
				Where("mine_code = ?", mineCode).
				Where("model_all_id = ?", m.ID).
				Where("model_version = ?", m.ModelVersion).
				Where("user = ? OR user = ?", user, userID).First(&modelInfo)
		} else {
			global.CMBP_DB.
				Where("mine_code = ?", mineCode).
				Where("model_all_id = ?", m.ID).
				Where("model_version = ?", fmt.Sprintf("%s.%d", m.ModelVersion, m.Edition)).
				Where("user = ? OR user = ?", user, userID).First(&modelInfo)
		}
		if *modelInfo.SyncFlag == 1 || modelInfo.SyncFlag == nil {
			isDownload = 1 // 已下载
		} else if *modelInfo.SyncFlag == 0 || *modelInfo.SyncFlag == -1 {
			isDownload = 2 // 下载中
		} else {
			isDownload = 0 // 可下载
		}
		d := map[string]interface{}{
			"model_id":                m.ID,
			"industry_code":           m.IndustryCode,
			"model_name":              m.ModelName,
			"model_chinese_name":      modelZhName,
			"model_description":       fmt.Sprintf("功能描述：%s；\\n技术描述：%s；\\n性能描述：%s。", m.ModelDescription, m.TechnicalDescription, m.PerformanceDescription),
			"technical_description":   m.TechnicalDescription,
			"performance_description": m.PerformanceDescription,
			"img_path":                "", // TODO
			"video_path":              "", // TODO
			"download_count":          m.DownloadCount,
			"collection_count":        m.CollectionCount,
			"view_count":              m.ViewCount,
			"model_kind":              m.ModelKind,
			"algorithm_id":            m.AlgorithmID,
			"build_way":               m.BuildWay,
			"upload_time":             m.UpdateTime.Format("2006-01-02 15:04:05"),
			"accuracy": func() interface{} {
				if m.Accuracy != nil {
					return strconv.Itoa(*m.Accuracy)
				}
				return 90 // 默认值为90
			}(),
			"test_duration": func() interface{} {
				if m.TestDuration != nil {
					return m.TestDuration
				}
				return 12 // 默认值为12分钟
			}(),
			"model_purpose": m.ModelPurpose,
			"is_have":       isHad,
			"is_download":   isDownload,
		}
		rspData = append(rspData, d)
	}
	return rspData
}

func (modelService *ModelService) GetModelStore(params systemReq.ModelStoreRqe, user system.Users, role string) (rspData interface{}, err error) {
	var autoUpdateModel []string

	if params.UpdateStatus == 1 {
		err := global.CMBP_DB.Model(&system.AutoUpdateInfo{}).Where("user_id = ? AND model_update_status = 1", user.ID).Pluck("model_id", &autoUpdateModel).Error
		if err != nil {
			return rspData, err
		}
		if len(autoUpdateModel) == 0 {
			return rspData, nil
		}
	} else if params.UpdateStatus == 2 {
		err := global.CMBP_DB.Model(&system.AutoUpdateInfo{}).Where("user_id = ? AND model_update_status = 0", user.ID).Pluck("model_id", &autoUpdateModel).Error
		if err != nil {
			return rspData, err
		}
		if len(autoUpdateModel) == 0 {
			return rspData, nil
		}
	}

	roleList := []string{"ROOT", "MODEL", "AIMODEL", "MODEL_OUT", "ADMIN"}
	found := false
	for _, r := range roleList {
		if role == r {
			found = true
			break
		}
	}
	//if strings.Contains(role, strings.Join(roleList, ",")) {
	if found {
		if params.ModelKind != nil {
			rspData, err := ForEachStoreModel(params, user, autoUpdateModel)
			if err != nil {
				return nil, err
			} else {
				return rspData, nil
			}
		} else {
			var modelKind []system.ModelKind
			global.CMBP_DB.Model(&system.ModelKind{}).Order("model_kind").Find(&modelKind)
			var allKindModel [][]interface{}
			for _, mk := range modelKind {
				if mk.ModelKind == 2 {
					params.Limit = 10
				}
				rspData, _ := ForEachStoreModel(params, user, autoUpdateModel)
				allKindModel = append(allKindModel, rspData)
			}
			return allKindModel, nil
		}
	}
	return rspData, nil
}

func ForEachStoreModel(params systemReq.ModelStoreRqe, user system.Users, autoUpdateModel []string) (rspData []interface{}, err error) {
	var modelALl []system.ModelAll
	QUERY := global.CMBP_DB.Model(&system.ModelAll{})
	if len(autoUpdateModel) > 0 {
		QUERY = QUERY.Where("id NOT IN ?", autoUpdateModel)
	}
	if params.IndustryCode != "" {
		QUERY = QUERY.Where("is_image IS NOT NULL").Where("industry_code = ?", params.IndustryCode).
			Where("model_kind = ?", params.ModelKind)
	} else {
		QUERY = QUERY.Where("is_image IS NOT NULL").Where("model_kind = ?", params.ModelKind)
	}

	totalCount := int64(0)
	QUERY.Count(&totalCount)

	if params.AlgorithmID != nil {
		QUERY = QUERY.Where("algorithm_id = ?", params.AlgorithmID)
	}
	if params.Mine != nil {
		QUERY = QUERY.Where("user = ?", user.ID)
	}
	if params.TestStatus != "" {
		testStatus := []string{"3", "4", "5", "6"}

		var idList []string
		if params.TestStatus == "null" {
			QUERY = QUERY.Where("audit_state = 1 AND test_status = NULL")
		} else if strings.Contains(params.TestStatus, strings.Join(testStatus, ",")) {
			err = QUERY.Find(&modelALl).Error
			if err != nil {
				return rspData, err
			}
			for _, ml := range modelALl {
				if ml.UUID != "" {
					status1, err1 := global.CMBP_REDIS.Get(context.Background(), ml.UUID+"status1").Int()
					status2, err2 := global.CMBP_REDIS.Get(context.Background(), ml.UUID+"status2").Int()
					if ml.IsImage {
						if ml.ModelKind == 1 && err1 == nil && status1 == 0 && params.TestStatus == "5" {
							idList = append(idList, ml.ID)
						} else if ml.ModelKind == 1 && err1 == nil && status1 == -1 && params.TestStatus == "6" {
							idList = append(idList, ml.ID)
						} else if (err1 == nil && status1 == 1) || (err2 == nil && status2 == 1) && params.TestStatus == "3" {
							idList = append(idList, ml.ID)
						} else if (err1 == nil && status1 == -1) || (err2 == nil && status2 == -1) && params.TestStatus == "4" {
							idList = append(idList, ml.ID)
						}
					} else {
						if err1 == nil && (status1 == 1 || status1 == 0) && params.TestStatus == "3" {
							idList = append(idList, ml.ID)
						} else if err1 == nil && status1 == -1 && params.TestStatus == "4" {
							idList = append(idList, ml.ID)
						}
					}
				} else {
					if strconv.Itoa(ml.TestStatus) == params.TestStatus &&
						(ml.TestStatus == 3 ||
							ml.TestStatus == 4 ||
							ml.TestStatus == 5 ||
							ml.TestStatus == 6) {
						// TODO 去重
						idList = append(idList, ml.ID)
					}
				}
			}
			QUERY = QUERY.Where("id IN ?", idList)
		} else {
			tst, _ := strconv.Atoi(params.TestStatus)
			QUERY = QUERY.Where("audit_state = 1 AND test_status = ?", tst)
		}
	}
	if params.NameOrDesc == "" {
		err = QUERY.Order("test_status DESC").Order("update_time DESC").Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&modelALl).Error
		if err != nil {
			return nil, err
		}

		rspDataList, err := FormatStoreModel(modelALl, user.ID)
		if err != nil {
			return nil, err
		} else {
			d := map[string]interface{}{
				"model_list":  rspDataList,
				"count":       len(rspDataList),
				"total_count": totalCount,
			}
			rspData = append(rspData, d)
			return rspData, nil

		}
	} else {
		err = QUERY.Limit(params.Limit).
			Offset(params.Limit * (params.Page - 1)).
			Find(&modelALl).Error
		rspDataList, err := FormatStoreModel(modelALl, user.ID)
		if err != nil {
			return nil, err
		} else {
			d := map[string]interface{}{
				"model_list":  rspDataList,
				"count":       len(rspDataList),
				"total_count": totalCount,
			}
			rspData = append(rspData, d)
			return rspData, nil
		}
	}
}

func FormatStoreModel(modelAll []system.ModelAll, userID string) (rspList []interface{}, err error) {
	rspDataList := []interface{}{}
	for _, m := range modelAll {
		var industryName string
		global.CMBP_DB.Model(&system.Industry{}).Where("industry_code = ?", m.IndustryCode).Pluck("industry_name", &industryName)
		var dCount int64
		global.CMBP_DB.Model(&system.Model{}).Joins("JOIN t_model_info ON t_model_info.model_name = t_model_all.model_name AND t_model_info.model_version = t_model_all.model_version").
			Where("t_model_all.id = ?", m.ID).Count(&dCount)
		var runtime system.RuntimeModels
		runtimeId := ""
		runtimeId = runtime.RuntimeID
		global.CMBP_DB.Model(&system.RuntimeModels{}).Where("model_all_id = ?", m.ID).First(&runtime)
		var onBoot, needGpu string
		if m.OnBoot {
			onBoot = "1"
		} else {
			onBoot = "0"
		}
		if m.NeedGPU {
			needGpu = "1"
		} else {
			needGpu = "0"
		}
		userName := "root"
		global.CMBP_DB.Model(&system.Users{}).Where("id = ?", m.User).Pluck("username", &userName)

		model2video := []map[string]interface{}{}
		var modelConfig []system.ModelConfig
		global.CMBP_DB.Where("model_name = ?", m.ModelName).Where("model_version = ?", m.ModelVersion).Find(&modelConfig)
		for _, mc := range modelConfig {
			d := map[string]interface{}{
				"id":              mc.Model2VideoConfigs[0].ID,
				"video__describe": mc.Model2VideoConfigs[0].VideoDescribe,
			}
			model2video = append(model2video, d)
		}

		modelKindName := ""
		global.CMBP_DB.Model(&system.ModelKind{}).Where("model_kind = ?", m.ModelKind).Pluck("kind_name", &modelKindName)

		alName := ""
		global.CMBP_DB.Model(&system.AlgorithmInfo{}).Where("algorithm_id = ?", m.AlgorithmID).Pluck("algorithm_name", &alName)

		var buildTask string
		if m.ModelKind == 2 {
			var buildTaskObj system.AutoBuildTask
			global.CMBP_DB.Model(system.AutoBuildTask{}).Where("model_all_id = ?", m.ID).First(&buildTaskObj)
			buildTask = buildTaskObj.ID
		}

		rspData := map[string]interface{}{
			"model_id":                m.ID,
			"industry_code":           m.IndustryCode,
			"industry_name":           industryName,
			"model_type":              m.ModelType,
			"model_name":              m.ModelName,
			"model_chinese_name":      m.ModelChineseName,
			"model_version":           m.ModelVersion,
			"model_description":       fmt.Sprintf("功能描述：%s；\\n技术描述：%s；\\n性能描述：%s。", m.ModelDescription, m.TechnicalDescription, m.PerformanceDescription),
			"technical_description":   m.TechnicalDescription,
			"performance_description": m.PerformanceDescription,
			"download_count":          dCount,
			"hardware_type":           strconv.Itoa(m.HardwareType),
			"is_image":                strconv.FormatBool(m.IsImage),
			"runtime_id":              runtimeId,
			"cmd":                     m.CMD,
			"json_url":                m.JsonURL,
			"img_url":                 m.ImgURL,
			"on_boot":                 onBoot,
			"need_gpu":                needGpu,
			"audit_state":             m.AuditState,
			"user":                    m.User,
			"developer":               userName,
			"upload_time":             m.UpdateTime,
			"business_dict":           FormatBusParams(m),
			"img_path":                "", // TODO
			"video_path":              "", // TODO
			"edit":                    m.User == userID,
			"model2video_config_list": model2video,
			"test_status":             GetTestStatus(m),
			"model_kind":              m.ModelKind,
			"model_kind_name":         modelKindName,
			"algorithm_id":            m.AlgorithmID,
			"algorithm_name":          alName,
			"build_task_id":           buildTask,
			"build_way":               m.BuildWay,
		}
		rspDataList = append(rspDataList, rspData)
	}
	return rspDataList, nil
}

func GetTestStatus(m system.ModelAll) int {
	if m.UUID != "" {
		status1, err1 := global.CMBP_REDIS.Get(context.Background(), m.UUID+"status1").Int()

		if m.IsImage {
			status2, err2 := global.CMBP_REDIS.Get(context.Background(), m.UUID+"status2").Int()
			if m.ModelKind == 1 {
				if err1 == nil && status1 == 0 {
					return 5
				} else if err1 == nil && status1 == -1 {
					return 6
				} else {
					return m.TestStatus
				}
			}
			if (err1 == nil && status1 == 1) || (err2 == nil && status2 == 1) {
				return 3 // 文件压缩中
			} else if (err1 == nil && status1 == -1) || (err2 == nil && status2 == -1) {
				return 4 // 文件压缩失败
			} else if err1 != nil && err2 != nil && m.MD5 != "" {
				return m.TestStatus //压缩成功，待测试
			} else {
				return 3
			}
		} else {
			if err1 == nil && (status1 == 0 || status1 == 1) {
				return 3
			} else if err1 == nil && status1 == -1 {
				return 4
			} else if err1 != nil && m.MD5 != "" {
				return m.TestStatus
			} else {
				return 4
			}
		}
	} else {
		return m.TestStatus
	}
}

func FormatBusParams(model system.ModelAll) []map[string]interface{} {

	jsonStr := model.BusinessParams
	businessType := map[string]interface{}{}
	if err := json.Unmarshal([]byte(model.BusinessType), &businessType); err == nil {
	} else {
		businessType = map[string]interface{}{}
	}
	businessList := []map[string]interface{}{}
	if jsonStr != "" {
		paramsMap := map[string][]interface{}{}
		if err := json.Unmarshal([]byte(jsonStr), &paramsMap); err == nil {
			for k, v := range paramsMap {
				business := map[string]interface{}{}
				business["business_name"] = k
				business["business_params"] = v

				for i, _ := range v {
					if _, ok := v[i].(map[string]interface{})["paras_text"]; !ok {
						v[i].(map[string]interface{})["paras_text"] = v[i].(map[string]interface{})["paras_desc"]
					}
				}

				if bt, ok := businessType[k]; ok {
					business["business_type"] = bt
				} else {
					business["business_type"] = []interface{}{}
				}
				businessList = append(businessList, business)
			}
		}
	} else if len(businessType) > 0 {
		for k, v := range businessType {
			business := map[string]interface{}{}
			business["business_name"] = k
			business["business_type"] = v
			business["business_params"] = []interface{}{}
			businessList = append(businessList, business)
		}
	}

	return businessList
}

func (modelService *ModelService) GetAutoUpdateTask(userID string) (rspData map[string]interface{}, err error) {
	var updateEnd system.AutoUpdateTask
	err = global.CMBP_DB.Where("user_id = ?", userID).First(&updateEnd).Error
	d := map[string]interface{}{
		"update_status":     1,
		"update_start_time": time.Now().Format("2006-01-02 15:04:05"),
		"success_count":     0,
		"failed_count":      0,
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return d, nil
		} else {
			return nil, err
		}
	}

	d["update_status"] = updateEnd.UpdateStatus
	if updateEnd.UpdateTime != nil {
		d["update_start_time"] = updateEnd.UpdateTime.Format("2006-01-02 15:04:05")
	}
	d["success_count"] = updateEnd.SuccessCount
	d["failed_count"] = updateEnd.FailedCount
	return d, nil
}

func (modelService *ModelService) GetModelKind(params systemReq.GetModelKind) (rspData interface{}, err error) {
	allModelKind := []string{}

	var modelKindList []system.ModelKind
	QUERY := global.CMBP_DB.Model(&system.ModelKind{}).Order("model_kind")
	if params.ModelPurpose == 2 {
		global.CMBP_DB.Where("model_purpose = 2").Group("model_kind").Pluck("model_kind", &allModelKind)
	}
	if params.Keywords != "" {
		QUERY = QUERY.Where("kind_name LIKE ? OR kind_name_en LIKE ? OR kind_desc LIKE ?", "%"+params.Keywords+"%")
	}
	if len(allModelKind) > 0 {
		QUERY = QUERY.Where("model_kind IN ?", allModelKind)
	}
	var totalCount = int64(0)
	QUERY.Count(&totalCount)
	if params.Page != nil && params.Limit != nil {
		QUERY = QUERY.Limit(*params.Limit).Offset(*params.Limit * (*params.Page - 1)).Find(&modelKindList)
	} else {
		QUERY = QUERY.Find(&modelKindList)
	}
	rspDataList := []interface{}{}
	for _, mk := range modelKindList {
		isDelete := int64(0)
		global.CMBP_DB.Model(&system.AlgorithmInfo{}).Where("model_kind = ?", mk.ModelKind).Count(&isDelete)
		if isDelete > 0 {
			isDelete = 0
		} else {
			isDelete = 1
		}
		d := map[string]interface{}{
			"model_kind":         mk.ModelKind,
			"model_kind_name":    mk.KindName,
			"model_kind_name_en": mk.KindNameEn,
			"model_kind_desc":    mk.KindDesc,
			"create_time":        mk.CreateTime.Format("2006-03-05 02:03:03"),
			"is_delete":          isDelete,
		}
		rspDataList = append(rspDataList, d)
	}
	//rspData = map[string]interface{}{
	//	"count":           totalCount,
	//	"model_kind_list": rspDataList,
	//}
	return rspDataList, nil
}

func (modelService *ModelService) GetHardWare(params systemReq.GetHardWare) (rspData interface{}, err error) {
	if params.Code != "" {
		var hardWare system.HardwareArch
		global.CMBP_DB.Where("code = ?", params.Code).First(&hardWare)
		rspData = map[string]interface{}{
			"code":                hardWare.Code,
			"name":                hardWare.Name,
			"desc":                hardWare.Desc,
			"real_channel_number": hardWare.RealChannelNumber,
		}
		return hardWare, nil
	} else {
		var hardWareList []system.HardwareArch
		QUERY := global.CMBP_DB.Model(&system.HardwareArch{}).Order("code")
		if params.Flag != nil && *params.Flag == 1 {
			hardWareCode := []int{}
			global.CMBP_DB.Model(&system.ModelAll{}).Where("build_flag = 1").Pluck("hardware_type", &hardWareCode)
			QUERY.Where("code IN ?", hardWareCode).Find(&hardWareList)
		} else {
			QUERY.Find(&hardWareList)
		}
		var rspData []interface{}
		for _, hardWare := range hardWareList {
			d := map[string]interface{}{
				"code":                hardWare.Code,
				"name":                hardWare.Name,
				"desc":                hardWare.Desc,
				"real_channel_number": hardWare.RealChannelNumber,
			}
			rspData = append(rspData, d)
		}

		return rspData, nil
	}
}

func (modelService *ModelService) GetModelOpsUuid(userID string) (interface{}, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", errors.New("生成UUID失败")
	}
	shortUUID := strings.ToUpper(strings.Join(strings.Split(uid.String(), "-"), ""))
	global.CMBP_REDIS.Set(context.Background(), shortUUID, fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"AIModel", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"plugins", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"Runtime", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600)
	model_dir := fmt.Sprintf("/home/tmp/%s", userID)

	if _, err := os.Stat(model_dir); os.IsNotExist(err) {
		// 如果不存在，则尝试创建目录
		if err := os.MkdirAll(model_dir, 0755); err != nil { // 0755代表权限为rwxr-xr-x
			return "", fmt.Errorf("无法创建目录: %v", err)
		}
	}
	rspData := map[string]interface{}{
		"uuid": shortUUID,
	}
	return rspData, nil
}

func (modelService *ModelService) GetIndustry(params systemReq.GetIndustry) (rspData interface{}, err error) {
	if params.IndustryCode == nil {
		var industryList []system.Industry

		allCodeList := []int{}
		global.CMBP_DB.Model(&system.ModelDetails{}).Where("model_purpose = 2").Group("industry_code").Pluck("industry_code", &allCodeList)
		QUERY := global.CMBP_DB.Model(&system.Industry{})
		if params.Keywords != "" {
			QUERY = QUERY.Where("industry_name LIKE ? OR industry_desc LIKE ? OR industry_name_en LIKE ?", "%"+params.Keywords+"%")
		}
		if len(allCodeList) > 0 {
			QUERY = QUERY.Where("industry_code IN ?", allCodeList)
		}
		totalCount := int64(0)
		QUERY.Count(&totalCount)
		if params.Page != nil && params.Limit != nil {
			QUERY.Limit(*params.Limit).Offset(*params.Limit * (*params.Page - 1)).Find(&industryList)
		} else {
			QUERY.Find(&industryList)
		}
		rspDataList := []interface{}{}
		for _, industry := range industryList {
			d := map[string]interface{}{
				"industry_code":      industry.IndustryCode,
				"industry_name":      industry.IndustryName,
				"industry_desc":      industry.IndustryDesc,
				"industry_name_en":   industry.IndustryNameEN,
				"industry_abbr_name": industry.IndustryAbbrName,
				"market_img_url":     "http://172.24.3.26/" + industry.MarketImgPath,  // TODO
				"outside_img_url":    "http://172.24.3.26/" + industry.OutsideImgPath, // TODO
				"create_time":        industry.CreateTime.Format("2006-01-01 02:03:05"),
			}
			rspDataList = append(rspDataList, d)
		}
		if params.Page != nil && params.Limit != nil {
			rspData = map[string]interface{}{
				"count":         totalCount,
				"industry_list": rspDataList,
			}
			return rspData, nil
		} else {
			return rspDataList, nil
		}
	} else {
		var industry system.Industry
		global.CMBP_DB.Where("industry_code = ?", params.IndustryCode).First(&industry)
		d := map[string]interface{}{
			"industry_code":      industry.IndustryCode,
			"industry_name":      industry.IndustryName,
			"industry_desc":      industry.IndustryDesc,
			"industry_name_en":   industry.IndustryNameEN,
			"industry_abbr_name": industry.IndustryAbbrName,
			"market_img_url":     "http://172.24.3.26/" + industry.MarketImgPath,  // TODO
			"outside_img_url":    "http://172.24.3.26/" + industry.OutsideImgPath, // TODO
			"create_time":        industry.CreateTime.Format("2006-01-01 02:03:05"),
		}
		return d, nil
	}
}
