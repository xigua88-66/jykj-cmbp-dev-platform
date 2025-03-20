package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"io"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/system"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	systemRsp "jykj-cmbp-dev-platform/server/model/system/response"
	"jykj-cmbp-dev-platform/server/utils"
	"jykj-cmbp-dev-platform/server/utils/upload"
	"mime/multipart"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
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
			fieldRsp.FieldImgUrl = global.CMBP_CONFIG.CMBPBase.CmbpUrl + fieldObj.FieldImgPath
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
					err = global.CMBP_DB.Preload("ModelField").Model(&system.ModelType{}).
						Where("model_field_id = ?", modelField.ID).
						Where("model_type in ?", typeCode).
						Find(&modelType).Error
					if err != nil {
						return nil, err
					}
				} else {
					err = global.CMBP_DB.Preload("ModelField").Model(&system.ModelType{}).
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

			rspData := []interface{}{}
			rspData = append(rspData, typeRspList)
			return rspData, nil
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

	var modelList []system.ModelMarketList
	if params.AlgorithmID != nil {
		QUERY, _, hadModels, appliedModels, user := SetModelMarketListCondition(params, userId)
		QUERY = QUERY.Where("algorithm_id = ?", *params.AlgorithmID)
		QUERY.Limit(params.Limit).Offset(params.Limit * (params.Page - 1)).Find(&modelList)
		dataList := FormatModelList(modelList, *params.ModelPurpose, hadModels, appliedModels, user.MineCode, userId)
		rspList := map[string]interface{}{
			"model_list": dataList,
			"count":      len(dataList),
		}
		return rspList, nil
	} else if params.ModelKind != nil && *params.ModelKind != 0 {
		QUERY, _, hadModels, appliedModels, user := SetModelMarketListCondition(params, userId)

		QUERY = QUERY.Where("model_kind = ?", *params.ModelKind)
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
		var rspList []map[string]interface{}
		for _, kindInfo := range modelKind {
			var limit int
			if kindInfo.ModelKind == 2 {
				limit = 10
			} else {
				limit = params.Limit
			}
			QUERY, _, hadModels, appliedModels, user := SetModelMarketListCondition(params, userId)

			realQuery := QUERY.Where("model_kind = ?", kindInfo.ModelKind)
			err = realQuery.Limit(limit).Offset(limit * (params.Page - 1)).Find(&modelList).Error
			dataList := FormatModelList(modelList, *params.ModelPurpose, hadModels, appliedModels, user.MineCode, userId)
			rspData := map[string]interface{}{
				"model_kind":  kindInfo.ModelKind,
				"count":       len(dataList),
				"total_count": len(dataList),
				"model_list":  dataList,
			}
			rspList = append(rspList, rspData)
		}
		return rspList, nil
	}
}

func SetModelMarketListCondition(params systemReq.ModelListReq, userId string) (query *gorm.DB, err error, hadModels []string, appliedModels []string, user system.Users) {
	//var hadModels []string
	global.CMBP_DB.Model(&system.AssetsManagement{}).Where("user_id = ? AND assets_type_id = 1", userId).Pluck("assets_id", &hadModels)
	//var appliedModels []string
	global.CMBP_DB.Model(&system.AssetsRecord{}).Where("user_id = ? AND apply_status = 0", userId).Pluck("assets_id", &appliedModels)

	//var user system.UserRoles
	global.CMBP_DB.Where("id = ?", userId).First(&user)

	QUERY := global.CMBP_DB.Model(&system.ModelMarketList{})
	if params.IndustryCode != nil {
		QUERY = QUERY.Where("industry_code = ?", params.IndustryCode)
	}
	if params.ModelPurpose != nil {
		QUERY = QUERY.Where("model_purpose = ?", params.ModelPurpose)
	}
	if params.NameOrDesc != "" {
		QUERY = QUERY.Where("model_name LIKE ? OR model_chinese_name LIKE ? OR model_description LIKE ? OR technical_description LIKE ? OR performance_description LIKE ?", "%"+params.NameOrDesc+"%")
		// TODO 添加用户访问记录
	}
	if params.Type != "" {
		QUERY = QUERY.Where("model_type = ?", params.Type)
	} else if params.Code != "" {
		QUERY = QUERY.Where("model_field = ?", params.Code)
	}
	if params.ReqType == 1 {
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
	} else if params.ReqType == 0 {

	} else {
		return nil, errors.New("参数错误"), nil, nil, system.Users{}

	}
	if params.Flag == 1 {
		QUERY = QUERY.Order("down_load_count DESC")
	} else {
		QUERY = QUERY.Order("update_time DESC")
	}
	return QUERY, nil, hadModels, appliedModels, user
}

func FormatModelList(modelList []system.ModelMarketList, isCloud int, hadModel []string, appliedModels []string, mineCode string, userID string) []interface{} {

	rspData := []interface{}{}

	for _, m := range modelList {
		modelZhName := m.ModelChineseName
		baseUrl := global.CMBP_CONFIG.CMBPBase.OssPath
		if isCloud != 0 {
			nameSplit := strings.Split(modelZhName, "-")
			if len(nameSplit) >= 2 {
				modelZhName = "某企业" + strings.Join(nameSplit[1:], "-")
			}
			baseUrl, _ = url.JoinPath(baseUrl, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia)
		}
		baseUrl, _ = url.JoinPath(baseUrl, global.CMBP_CONFIG.CMBPBase.ModelMarketMedia)
		modelFullName := m.ModelName + "V" + m.ModelVersion

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
			if m.ModelKind == 1 {
				modelFullName = m.ModelName + "." + string(rune(m.Edition))
			}
		}
		imgPath, _ := url.JoinPath(baseUrl, modelFullName+".jpg")
		videoPath, _ := url.JoinPath(baseUrl, modelFullName+".mp4")
		if modelInfo.SyncFlag == nil || *modelInfo.SyncFlag == 1 {
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
			"img_path":                imgPath,
			"video_path":              videoPath,
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
			var allKindModel []interface{}
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

func ForEachStoreModel(params systemReq.ModelStoreRqe, user system.Users, autoUpdateModel []string) (rspData interface{}, err error) {
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
					if strconv.Itoa(*ml.TestStatus) == params.TestStatus &&
						(*ml.TestStatus == 3 ||
							*ml.TestStatus == 4 ||
							*ml.TestStatus == 5 ||
							*ml.TestStatus == 6) {
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
			//rspData = append(rspData, d)
			return d, nil

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
			//rspData = append(rspData, d)
			return d, nil
		}
	}
}

func FormatStoreModel(modelAll []system.ModelAll, userID string) (rspList []interface{}, err error) {
	rspDataList := []interface{}{}
	for _, m := range modelAll {
		var industryName string
		global.CMBP_DB.Model(&system.Industry{}).Where("industry_code = ?", m.IndustryCode).Pluck("industry_name", &industryName)
		var dCount int64
		global.CMBP_DB.Model(&system.Model{}).Joins("JOIN t_model_all ON t_model_all.model_name = t_model_all.model_name AND t_model_info.model_version = t_model_all.model_version").
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

		modelFullName := m.ModelName + "V" + m.ModelVersion
		baseUrl := global.CMBP_CONFIG.CMBPBase.OssPath
		baseUrl, _ = url.JoinPath(baseUrl, global.CMBP_CONFIG.CMBPBase.ModelWareHouseMedia)

		if m.ModelKind == 1 {
			baseUrl, _ = url.JoinPath(baseUrl, global.CMBP_CONFIG.CMBPBase.ModelMarketMedia)
			modelFullName = modelFullName + "." + string(rune(m.Edition))
		}

		imgPath, _ := url.JoinPath(baseUrl, modelFullName+".jpg")
		videoPath, _ := url.JoinPath(baseUrl, modelFullName+".mp4")

		rspData := map[string]interface{}{
			"model_id":                m.ID,
			"industry_code":           m.IndustryCode,
			"industry_name":           industryName,
			"model_type":              m.ModelType,
			"model_name":              m.ModelName,
			"model_chinese_name":      m.ModelChineseName,
			"model_version":           m.ModelVersion,
			"model_description":       fmt.Sprintf("功能描述：%s；\n技术描述：%s；\n性能描述：%s。", m.ModelDescription, m.TechnicalDescription, m.PerformanceDescription),
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
			"img_path":                imgPath,
			"video_path":              videoPath,
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
					return *m.TestStatus
				}
			}
			if (err1 == nil && status1 == 1) || (err2 == nil && status2 == 1) {
				return 3 // 文件压缩中
			} else if (err1 == nil && status1 == -1) || (err2 == nil && status2 == -1) {
				return 4 // 文件压缩失败
			} else if err1 != nil && err2 != nil && m.MD5 != "" {
				return *m.TestStatus //压缩成功，待测试
			} else {
				return 3
			}
		} else {
			if err1 == nil && (status1 == 0 || status1 == 1) {
				return 3
			} else if err1 == nil && status1 == -1 {
				return 4
			} else if err1 != nil && m.MD5 != "" {
				return *m.TestStatus
			} else {
				return 4
			}
		}
	} else {
		return *m.TestStatus
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

func (modelService *ModelService) UploadModel(params systemReq.UploadModelStoreReq, videoFile, imgFile *multipart.FileHeader, userId string) (map[string]string, error) {
	var busListDict []map[string]interface{}
	if params.BusinessDict != "" {
		err := json.Unmarshal([]byte(params.BusinessDict), &busListDict)
		if err != nil {
			return nil, errors.New("business_dict json解析失败 ：" + err.Error())
		}
		//_, ok := busListDict.([]interface{})
		//if !ok {
		//	return nil, errors.New("business_dict不符合规则")
		//}
	}
	var m system.ModelAll
	global.CMBP_DB.Where("model_name = ?", params.ModelName).Where("model_version = ?", params.ModelVersion).First(&m)

	if m.ID != "" {
		return nil, errors.New("模型重复")
	}

	modelName := params.ModelName + "V" + params.ModelVersion
	// OBS模型存放目录
	modelDir := fmt.Sprintf(global.CMBP_CONFIG.CMBPBase.OssModelPath, modelName)

	_, err := os.Stat(modelDir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(modelDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// 图片和视频保存目录 #/home/models/models
	mediaDir := global.CMBP_CONFIG.CMBPBase.OssModelMedia
	_, err = os.Stat(mediaDir)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(mediaDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	// 大文件切片保存 /home/models/fileSave/用户id/
	modelZipPath := "/home/models/fileSave/" + userId + "/" + modelName + ".zip"
	err = SaveSliceFile(params.TaskID, userId, modelZipPath)
	if err != nil {
		return nil, err
	}

	// 将保存的压缩包解压到"/home/tmp/userId/uuid"下面
	tmpPath := fmt.Sprintf("/home/tmp/%s/", userId)
	_, err = os.Stat(tmpPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(tmpPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	targetPath := tmpPath + params.UUID
	global.CMBP_REDIS.Set(context.Background(), params.UUID, 1, 3600*time.Second)
	_, err = utils.Unzip(modelZipPath, targetPath)
	if err != nil {
		return nil, err
	}

	err = os.RemoveAll(modelZipPath)
	if err != nil {
		return nil, err
	}

	// 拷贝End到模型临时目录下
	err = utils.CopyDir("/home/models/AIMonitorEnd", targetPath)
	if err != nil {
		return nil, err
	}
	global.CMBP_REDIS.Set(context.Background(), modelName, 1, 3600*time.Second)

	// 生成加密Python文件
	isDataModel := 0
	dirList, err := os.ReadDir(targetPath)
	if !(len(dirList) == 1 && dirList[0].Name() == "app") {
		isDataModel = 1
	}
	dirTree := utils.GetDirTree(targetPath, targetPath, isDataModel)
	jsonTree, err := utils.TreeToJson(dirTree)
	err = global.CMBP_REDIS.Set(context.Background(), fmt.Sprintf("%s_AIModel_tmp_dirs", params.UUID), jsonTree, 3600*time.Second).Err()
	if err != nil {
		return nil, err
	}

	// copy到model_warehouse
	modelWareHouse := fmt.Sprintf("/home/model_warehouse/%s", modelName)

	_, err = os.Stat(modelWareHouse)
	if os.IsExist(err) {
		err = os.Remove(modelWareHouse)
		if err != nil {
			return nil, err
		}
	}

	// 更新模型可编辑代码到model_warehouse
	err = utils.CopyDir(targetPath, modelWareHouse)
	if err != nil {
		return nil, errors.New("模型仓库model_warehouse拷贝失败: " + err.Error())
	}

	processor := "x86"
	armHardWare := []int{1, 3, 5, 7, 14, 30, 91, 100}
	for _, h := range armHardWare {
		if params.HardwareType == h {
			processor = "arm"
			break
		}
	}
	// 生成加密文件并转so
	randStr := utils.UniqueRandomStr()
	encrypted, err := utils.AddEncryptFile(targetPath, randStr, processor)
	if err != nil {
		return nil, errors.New("model添加加密文件失败" + err.Error())
	}

	// 上传plugins到OBS
	pluginsPath := path.Join(targetPath, "plugins")
	_, err = os.Stat(pluginsPath)
	if os.IsExist(err) {
		os.RemoveAll(pluginsPath)
	}
	err = os.MkdirAll(pluginsPath, 0755)
	if err != nil {
		return nil, err
	}
	global.CMBP_REDIS.Set(context.Background(), params.UUID+"plugins", 1, 3600*time.Second)
	err = utils.CopyEnd(pluginsPath, randStr, processor)
	if err != nil {
		return nil, err
	}
	py2so, err := utils.Plugins2So(targetPath, processor)
	if !py2so && err != nil {
		global.CMBP_LOG.Error(fmt.Sprintf("业务模型转换so失败: %s", err.Error()))
		return nil, errors.New("业务模型转换so失败")
	}
	exists, _ := utils.PathExists(targetPath + "plugins/plugins")
	if exists {
		err = os.RemoveAll(targetPath + "plugins/plugins")
		if err != nil {
			return nil, err
		}
	}

	err = utils.CopyDir(pluginsPath, targetPath+"plugins/plugins")
	if err != nil {
		return nil, err
	}
	// 压缩plugin目录
	err = utils.CompressZip(pluginsPath, modelDir, modelName, false)
	if err != nil {
		return nil, errors.New("压缩plugins目录失败：" + err.Error())
	}
	global.CMBP_REDIS.Del(context.Background(), params.UUID+"plugins")

	// 上传AIModel到OBS
	AiModelPath := fmt.Sprintf("/home/tmp/%s/%s/", userId, params.UUID+"AIModel")
	global.CMBP_LOG.Info(fmt.Sprintf("识别模型临时目录：%s", AiModelPath))
	os.Mkdir(AiModelPath, 0755)
	global.CMBP_REDIS.Set(context.Background(), params.UUID+"AIModel", 1, 3600*time.Second)
	err = utils.CopyEnd(AiModelPath, randStr, processor)
	if err != nil {
		return nil, err
	}

	zipCmd := fmt.Sprintf("cd /home/tmp/%s/%s/AIModel && zip -r -1 -P %s AIModel.zip * && mv AIModel.zip ..", userId, params.UUID, randStr)

	// 创建命令
	command := exec.Command("/bin/bash", "-c", zipCmd)
	global.CMBP_LOG.Info(fmt.Sprintf("识别模型压缩命令: %s", command.String()))
	// 执行命令
	output, err := command.CombinedOutput()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			if waitStatus.ExitStatus() == 1 {
				return nil, errors.New("识别模型包压缩失败：" + err.Error() + string(output))
			}
		}
		global.CMBP_LOG.Error(fmt.Sprintf("执行命令时出现错误: %v, 输出: %s\n", err, output))
	}
	// 如果命令执行成功，这里可以继续处理其他逻辑
	global.CMBP_LOG.Info("识别模型包压缩成功")

	err = utils.FileCopy(fmt.Sprintf("/home/tmp/%s/%s/AIModel.zip", userId, params.UUID), fmt.Sprintf("/home/tmp/%s/%sAIModel/AIModel.zip", userId, params.UUID))
	if err != nil {
		return nil, errors.New("识别模型拷贝失败: " + err.Error())
	}
	utils.CompressZip(AiModelPath, modelDir, modelName+"AIModel", false)
	global.CMBP_REDIS.Del(context.Background(), params.UUID+"AIModel")

	// 压缩完整的模型包
	utils.CompressZip(targetPath, modelDir, modelName, false)
	global.CMBP_REDIS.Del(context.Background(), modelName+".zip")

	// docker镜像压缩
	runTimePath := fmt.Sprintf("/home/tmp/%s/%s/", userId, params.UUID+"Runtime")
	if params.IsImage != nil && *params.IsImage == 1 {
		os.MkdirAll(runTimePath, 0755)
		global.CMBP_REDIS.Set(context.Background(), params.UUID+"Runtime", 1, 3600*time.Second)
		utils.CopyEnd(runTimePath, "", "")
		os.ReadDir(targetPath)

		for _, f := range dirList {
			if len(f.Name()) > 4 && f.Name()[len(f.Name())-4:] == ".tar" {
				utils.FileCopy(filepath.Join(targetPath, f.Name()), runTimePath)
				break
			}
		}
		global.CMBP_REDIS.Set(context.Background(), modelName+"Runtime.zip", 1, 3600*time.Second)
		runTimeZipCmd := fmt.Sprintf("cd /home/tmp/%s/%s && zip -r -1 ../%s.zip *", userId, params.UUID+"Runtime", modelName+"Runtime")

		command = exec.Command("/bin/bash", "-c", runTimeZipCmd)

		// 执行命令
		output, err = command.CombinedOutput()
		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				if waitStatus.ExitStatus() == 1 {
					return nil, errors.New("镜像包压缩失败")
				}
			}
			fmt.Printf("执行命令时出现错误: %v, 输出: %s\n", err, output)
		}
		// 上传到minio TODO 增加兼容华为云OBS的上传功能
		runTimeZip := fmt.Sprintf("/home/tmp/%s/%s", userId, modelName+"Runtime.zip")
		var minio utils.MinIO
		err = minio.StreamUpload("obs-isf", fmt.Sprintf("models/models/%s/%s/", modelName, modelName+"Runtime.zip"), runTimeZip)
		if err != nil {
			return nil, errors.New("镜像压缩包上传对象存储失败" + err.Error())
		}
		_, err = os.Stat(runTimeZip)
		if os.IsExist(err) {
			os.RemoveAll(runTimeZip)
		}
		global.CMBP_REDIS.Del(context.Background(), modelName+"Runtime.zip")
		global.CMBP_REDIS.Del(context.Background(), params.UUID+"Runtime")

	}
	// 生成MD5
	md5sum, err := utils.FileMD5(filepath.Join(modelDir, modelName+".zip"))
	if err != nil {
		return nil, err
	}

	// 删除剩余临时目录
	err = os.RemoveAll(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	err = os.RemoveAll(runTimePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// 保存图片和视频到对象存储
	err = utils.SaveFile(imgFile, filepath.Join("/home/OBS/models/models_media/", modelName+".jpg"))
	if err != nil {
		return nil, err
	}
	err = utils.SaveFile(videoFile, filepath.Join("/home/OBS/models/models_media/", modelName+".mp4"))
	if err != nil {
		return nil, err
	}

	// 解析业务模型参数

	busList := ""
	busDict := map[string]interface{}{}
	busApi := map[string]interface{}{}
	busType := map[string]interface{}{}

	if busListDict != nil {
		for _, bus := range busListDict {
			busName := bus["business_name"]
			if busName != "" {
				busList = busList + busName.(string) + "|"
			}
			busParams := bus["business_params"]
			apis := bus["business_apis"]
			busType[busName.(string)] = bus["business_type"]
			busDict[busName.(string)] = busParams
			busApi[busName.(string)] = apis
		}
		busList = busList[:len(busList)-1]
	}

	var modelKind int
	global.CMBP_DB.Model(&system.AlgorithmInfo{}).Where("algorithm_id = ?", params.AlgorithmID).Pluck("model_kind", &modelKind)

	var modelFiled system.ModelField
	err = global.CMBP_DB.Model(&system.ModelField{}).Joins("JOIN t_model_type ON t_model_type.model_field_id = t_model_field.id").First(&modelFiled).Error
	if err != nil {
		return nil, err
	}
	industryCode := modelFiled.IndustryCode
	fieldCode := modelFiled.Code

	t := time.Now().Format("20060102150405")
	isDocker := false
	if params.IsImage != nil && *params.IsImage == 1 {
		isDocker = true
	}
	busParams, err := json.Marshal(busDict)
	if err != nil {
		return nil, err
	}
	busTypeStr, err := json.Marshal(busType)
	if err != nil {
		return nil, err
	}
	busApiStr, err := json.Marshal(busApi)
	if err != nil {
		return nil, err
	}

	needGpu := false
	if params.NeedGPU != nil && *params.NeedGPU == 1 {
		needGpu = true
	}
	onBoot := false
	if params.OnBoot == 1 {
		onBoot = true
	}

	testStatus := 0
	modelAll := system.ModelAll{
		ModelType:              params.ModelType,
		FieldCode:              fieldCode,
		ModelName:              params.ModelName,
		ModelChineseName:       params.ModelChineseName + "-Mark-" + t,
		ModelVersion:           params.ModelVersion,
		ModelDescription:       params.ModelDescription,
		TechnicalDescription:   params.TechnicalDescription,
		PerformanceDescription: params.PerformanceDescription,
		Advantage:              params.Advantage,
		MD5:                    md5sum,
		EncryptPassword:        encrypted,
		HardwareType:           params.HardwareType,
		IsImage:                isDocker,
		CMD:                    params.Cmd,
		JsonURL:                strings.TrimSpace(params.JsonUrl),
		ImgURL:                 strings.TrimSpace(params.ImgUrl),
		BusinessList:           busList,
		BusinessParams:         string(busParams),
		ModelPurpose:           params.ModelPurpose,
		BusinessAPI:            string(busApiStr),
		AiModelAPI:             params.AIModelAPI,
		BusinessType:           string(busTypeStr),
		NeedGPU:                needGpu,
		OnBoot:                 onBoot,
		TestDuration:           params.TestDuration,
		Accuracy:               params.Accuracy,
		IsRealChannel:          params.IsRealChannel,
		AuditState:             1,
		User:                   userId,
		TestStatus:             &testStatus,
		NewModelFlag:           true,
		IndustryCode:           industryCode,
		ModelKind:              modelKind,
		AlgorithmID:            params.AlgorithmID,
		UUID:                   params.UUID,
		BuildWay:               1,
		IsProcess:              1,
	}

	err = global.CMBP_DB.Save(&modelAll).Error
	if err != nil {
		return nil, err
	}
	// TODO 1、Redis 删除缓存 2、图片视频加水印
	rspData := map[string]string{
		"model_id": modelAll.ID,
	}
	return rspData, nil
}

func SaveSliceFile(taskId, userId, modelZipPath string) error {
	zf, err := os.Create(modelZipPath)
	if err != nil {
		return err
	}
	defer func() {
		// 在函数退出前尝试关闭zip文件
		if err := zf.Close(); err != nil {
			fmt.Println("Error closing zip file:", err)
		}
	}()

	chunk := 0
	for {
		sourceFilePath := "/home/models/fileSave/" + userId + "/" + taskId + strconv.Itoa(chunk)
		if _, err := os.Stat(sourceFilePath); os.IsNotExist(err) {
			if chunk == 0 {
				return errors.New("分片文件不存在" + sourceFilePath)
			} else {
				break
			}
		}
		fileToAppend, err := os.Open(sourceFilePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zf, fileToAppend)
		if err != nil {
			return err
		}
		err = fileToAppend.Close()
		if err != nil {
			return err
		}

		err = os.Remove(sourceFilePath)
		if err != nil {
			fmt.Println("分配片删除失败：", err.Error())
		}
		chunk++
	}
	return nil
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
	global.CMBP_REDIS.Set(context.Background(), shortUUID, fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600*time.Second)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"AIModel", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600*time.Second)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"plugins", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600*time.Second)
	global.CMBP_REDIS.Set(context.Background(), shortUUID+"Runtime", fmt.Sprintf("/home/tmp/%s/%s", userID, shortUUID), 3600*time.Second)
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
				"market_img_url":     path.Join(global.CMBP_CONFIG.CMBPBase.CmbpUrl, industry.MarketImgPath),
				"outside_img_url":    path.Join(global.CMBP_CONFIG.CMBPBase.CmbpUrl, industry.OutsideImgPath),
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
			"market_img_url":     path.Join(global.CMBP_CONFIG.CMBPBase.CmbpUrl, industry.MarketImgPath),
			"outside_img_url":    path.Join(global.CMBP_CONFIG.CMBPBase.CmbpUrl, industry.OutsideImgPath),
			"create_time":        industry.CreateTime.Format("2006-01-01 02:03:05"),
		}
		return d, nil
	}
}

func (modelService *ModelService) UnPublishModel(uuid string, token string) error {
	tempPath := fmt.Sprintf("/home/tmp/%s", uuid)
	_, err := os.Stat(tempPath)
	if !os.IsNotExist(err) { // 存在
		err := os.RemoveAll(tempPath)
		if err != nil {
			return err
		}
	}
	global.CMBP_REDIS.Del(context.Background(), uuid)
	global.CMBP_REDIS.Del(context.Background(), "upload_"+token)
	return nil
}

func (modelService *ModelService) CancelUpload(uuid, userID string) (err error) {
	busModelDir := fmt.Sprintf("/home/models/BusinessModelsLibrary/%s/%s", userID, uuid)
	busModelDirBak := fmt.Sprintf("/home/models/BusinessModelsLibrary/%s/%s_bak", userID, uuid)
	totalCount := int64(0)
	global.CMBP_DB.Model(&system.BusinessModelAll{}).Where("id = ?", uuid).Count(&totalCount)
	_, bakErr := os.Stat(busModelDirBak)
	_, moderr := os.Stat(busModelDir)
	if !os.IsNotExist(bakErr) {
		err = os.RemoveAll(busModelDirBak)
	} else if os.IsNotExist(moderr) && totalCount > 0 {
		err = os.RemoveAll(busModelDir)
	}
	if err != nil {
		return errors.New(fmt.Sprintf("文件删除失败， %s", err.Error()))
	}
	redis := global.CMBP_REDIS.Get(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s", userID, uuid))
	pid := redis.Val()
	if pid != "" {
		err = utils.KillProcess(pid)
		if err != nil {
			return err
		}
		global.CMBP_REDIS.Del(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s", userID, uuid))
		global.CMBP_REDIS.Del(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s_port", userID, uuid))
	}
	aiPid := global.CMBP_REDIS.Get(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s_ai", userID, uuid)).Val()
	if aiPid != "" {
		err = utils.KillProcess(aiPid)
		if err != nil {
			return err
		}
		global.CMBP_REDIS.Del(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s", userID, uuid))
		global.CMBP_REDIS.Del(context.Background(), fmt.Sprintf("jupyter_lab_%s_%s_port", userID, uuid))
	}
	err = global.CMBP_DB.Model(&system.Notebook{}).Where("user_id = ?", userID).Where("uuid = ?", uuid).Where("status = 1").Update("status", 0).Error
	if err != nil {
		return err
	}
	return nil
}

func (modelService *ModelService) GetTestFreeApplication(params systemReq.GetTestFreeApply) (rspData interface{}, err error) {
	selectStatus := []int{99, 100, 101}

	modelId := []string{}
	QUERY := global.CMBP_DB.Model(&system.ApplicationRecord{}).
		Joins("JOIN t_model_all ON t_application_record.model_id = t_model_all.id").
		Where("application_status IN ?", selectStatus)
	if params.NameOrDesc != "" {
		var user system.Users
		var model system.ModelAll
		global.CMBP_DB.Where("username = ?", params.NameOrDesc).First(&user)
		global.CMBP_DB.Where("model_name LIKE ? OR model_chinese_name LIKE ?", "%s"+params.NameOrDesc+"%s").First(&model)
		QUERY.Where("reason LIKE ? OR user = ? OR model_id = ?", "%s"+params.NameOrDesc+"%s", user.ID, model.ID).
			Order("t_application_record.create_time DESC").Pluck("model_id", &modelId)
	} else {
		QUERY.Order("t_application_record.create_time DESC").Pluck("model_id", &modelId)
	}
	//count := len(modelId)
	modelId = modelId[(*params.Page-1)**params.Limit : *params.Page**params.Limit]

	var testFreeModelRes []system.TestFreeModelRes

	global.CMBP_DB.Table("t_model_all AS mal").
		Joins("JOIN t_model_type AS mty ON mal.model_type = mty.model_type").
		Joins("JOIN t_model_field AS mf ON mty.model_field_id = mf.id").
		Joins("JOIN t_hardware_arch AS hard ON hard.code = mal.hardware_type").
		Joins("JOIN t_user_info AS us ON mal.user = us.id").
		Joins("JOIN t_application_record AS record ON mal.id = record.model_id").
		Where("mal.id IN ?", modelId).
		Select("mal.*, mty.model_type_desc AS model_type_desc, mf.name AS model_field_desc, hard.name AS hardware_type_name, us.username AS developer").
		Scan(&testFreeModelRes)

	if err != nil {
		return nil, err
	}
	return testFreeModelRes, nil
	//for _, mal := range modelAllList {
	//	d := map[string]interface{}{
	//		"model_id": mal.ID,
	//		"model_type": mal.ModelType,
	//		"model_type_desc":GetModelType(attribute="model_type"),
	//		"model_field": GetModelFieldCode(attribute="model_type"),
	//		"model_field_desc": GetModelFieldName(attribute="model_type"),
	//		"model_name": fields.String,
	//		"model_chinese_name": fields.String,
	//		"model_version": fields.String,
	//		"model_description": fields.String,
	//		"technical_description": fields.String,
	//		"performance_description": fields.String,
	//		"download_count": DownloadCount(attribute="id"),
	//		"collection_count": CollectionCount(attribute="id"),
	//		"view_count": ViewCount(attribute="id"),
	//		"hardware_type": fields.String(attribute="hardware_type"),
	//		"hardware_type_name": GetHardwareTypeName(attribute="hardware_type"),
	//		"hardware_type_desc": GetHardwareType(attribute="hardware_type"),
	//		"is_image": GetStr(attribute="is_image"),
	//		"is_real_channel": GetIsRealChannel(attribute="id"),
	//		"cmd": fields.String,
	//		"json_url": fields.String,
	//		"img_url": fields.String,
	//		"on_boot": GetStr(attribute="on_boot"),
	//		"need_gpu": GetStr(attribute="need_gpu"),
	//		"audit_state": fields.String(attribute="audit_state"),
	//		"user": fields.String,
	//		"developer": GetUser(default="root", attribute="user"),  # 模型开发上传人员
	//		"upload_time": fields.String(attribute="update_time"),
	//		"business_dict": JsonToDict(attribute="business_params"),
	//		"img_path": GetModelImgPath(attribute="id"),
	//		"img2_path": GetModelImg2Path(attribute="id"),
	//		"video_path": GetModelVideoPath(attribute="id"),
	//		"edit": GetEdit(attribute="user"),
	//		"model2video_config_list": GetModel2VideoList(attribute='id'),
	//		"test_status":GetTestStatus(attribute='id'),
	//		"reason":GetReason(attribute='id'),          #模型免测原因
	//		"phone":GetApplicant(attribute='id'),
	//		"application_time":GetApplicantTime(attribute='id'),  #模型免测申请时间
	//		"metadata_update_flag":GetUpdateFlag1(attribute='id'),
	//		"ai_model_update_flag":GetUpdateFlag2(attribute='id'),
	//		"business_model_update_flag":GetUpdateFlag3(attribute='id'),
	//		"runtime_update_flag":GetUpdateFlag4(attribute='id'),
	//		"model_kind": GetModelKind(attribute='id'),
	//		"model_status":GetModelStatus(attribute='id'),
	//		"process_type": GetProcessType(attribute='id'),
	//		"power_list": QueryPower(attribute='id'),
	//	}
	//}
}

func (modelService *ModelService) TestFreeApplyCation(params systemReq.TestFreeApplyCation, userID string) (rspData interface{}, err error) {
	var modelAll system.ModelAll
	err = global.CMBP_DB.Model(&system.ModelAll{}).Where("id=?", params.ModelID).First(&modelAll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("模型不存在，请传递正确模型id")
		}
		return nil, err
	}

	isProcess := modelAll.IsProcess
	var applyCationOld system.ApplicationRecord
	err = global.CMBP_DB.Model(&system.ApplicationRecord{}).Where("model_id=?", modelAll.ID).First(&applyCationOld).Error
	if err != nil {
		global.CMBP_LOG.Warn(fmt.Sprintf("记录查询报错是:%v", err.Error()))
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	modelUpdateStatus := 0
	if isProcess == 1 {
		*modelAll.TestStatus = 99
		modelAll.IsProcess = 0
		err = global.CMBP_DB.Save(&system.ApplicationRecord{
			ModelID:           params.ModelID,
			ApplicationStatus: 99,
			Reason:            params.Reason,
			User:              userID,
			ProcessType:       1,
		}).Error
		if err != nil {
			global.CMBP_LOG.Error(fmt.Sprintf("数据库报错报错了：:%v", err.Error()))
			return nil, err
		}
	} else {
		modelUpdateStatus = 1
		if modelAll.ModelKind != 1 {
			*modelAll.TestStatus = 201

			// TODO 更新边缘平台下载的模型
		} else {
			modelAll.TestStatus = nil
		}
		err = global.CMBP_DB.Save(&modelAll).Error
		if err != nil {
			global.CMBP_LOG.Error(fmt.Sprintf("数据库报错报错了：:%v", err.Error()))
			return nil, err
		}

		var developer string
		err = global.CMBP_DB.Model(&system.Users{}).Where("id=?", modelAll.User).Pluck("username", &developer).Error
		if err != nil {
			return nil, err
		}
		if applyCationOld.ID != "" {
			// TODO 判断用了这个模型的用户，给他们的钉钉发消息
		}
		global.CMBP_LOG.Info(fmt.Sprintf("开发者:%v", developer))
	}

	if applyCationOld.ApplicationStatus == 101 {
		err = global.CMBP_DB.Save(system.ModelUpdateRecord{
			UserID:           userID,
			ModelID:          params.ModelID,
			ModelChineseName: modelAll.ModelChineseName,
			Status:           modelUpdateStatus,
		}).Error
		if err != nil {
			global.CMBP_LOG.Error(fmt.Sprintf("数据库报错报错了：:%v", err.Error()))
			return nil, err
		}
	}
	if applyCationOld.ID != "" {
		err = global.CMBP_DB.Delete(&system.ApplicationRecord{}, "model_id=?", params.ModelID).Error
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func asyncPublishModel(modelID, userName string) {
	global.CMBP_LOG.Info(fmt.Sprintf("更新矿编码下模型：:%v", modelID))
	var modelAll system.ModelAll
	err := global.CMBP_DB.Model(&system.ModelAll{}).Where("id=?", modelID).First(&modelAll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.CMBP_LOG.Error(fmt.Sprintf("模型不存在:%v", modelID))
			return
		}
	}
	var mineModel []system.Model
	err = global.CMBP_DB.Model(&system.Model{}).Where("model_all_id=?", modelID).Find(&mineModel).Error

}

func (modelService *ModelService) CheckName(params systemReq.CheckName, userId string) (interface{}, error) {
	exists := false
	obsPath := ""
	var count int64
	zipName := params.ModelName + "V" + params.ModelVersion

	switch params.Flag {
	case 1:
		if params.ModelVersion != "" {
			obsPath = fmt.Sprintf("/OBS/models/models/%s/%s.zip", zipName, zipName)
		} else {
			return nil, errors.New("model_version参数不正确")
		}
		global.CMBP_DB.Model(&system.ModelAll{}).Where("model_name = ?", params.ModelName).Where("model_version = ?", params.ModelVersion).Count(&count)
		if count > 0 {
			exists = true
		}
	case 2:
		if params.UUID != "" {
			obsPath = fmt.Sprintf("/OBS/BusinessModelLibray/%s.zip", params.UUID)
		}
		global.CMBP_DB.Model(&system.BusinessModelAll{}).Where("model_name = ?", params.ModelName).Where("model_version = ?", userId).Count(&count)
		if count > 0 {
			exists = true
		}
	case 3:
		obsPath = fmt.Sprintf("/OBS/ModelLibrary/%s.zip", params.ModelName)
		global.CMBP_DB.Model(&system.AIModelAll{}).Where("model_name = ?", params.ModelName).Count(&count)
		if count > 0 {
			exists = true
		}
	case 4:
		obsPath = fmt.Sprintf("/OBS/DataAnalyzeModelLibrary/%s.zip", params.ModelName)
		global.CMBP_DB.Model(&system.DataAnalyzeModelAll{}).Where("model_name = ?", params.ModelName).Count(&count)
		if count > 0 {
			exists = true
		}
	case 5:
		if strings.ToLower(params.ModelName) != params.ModelName {
			return nil, errors.New("镜像名称不能包含大写字母")
		}
		obsPath = fmt.Sprintf("/OBS/RuntimeLibrary/%s.zip", zipName)
		global.CMBP_DB.Model(&system.RuntimeAll{}).Where("name = ?", params.ModelName).Where("tag = ?", params.ModelVersion).Count(&count)
		if count > 0 {
			exists = true
		}
		ctx := context.Background()

		// 创建Docker客户端实例
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		defer cli.Close()

		// 设置您的镜像名称 todo 常量参数抽取
		imageName := fmt.Sprintf(global.CMBP_CONFIG.CMBPBase.DockerRegistry, params.ModelName, params.ModelVersion)

		// 检查镜像是否存在
		_, _, err = cli.ImageInspectWithRaw(ctx, imageName)
		if err != nil {
			fmt.Printf("镜像不存在: %v\n", err)
		} else {
			exists = true
		}
	default:
		return nil, errors.New("flag参数不正确")
	}
	if obsPath != "" {
		_, err := os.Stat(obsPath)
		if os.IsExist(err) {
			exists = true
		}
	}

	resData := map[string]bool{
		"exist": exists,
	}
	return resData, nil
}

func (modelService *ModelService) DeleteModel(modelID string) (resData interface{}, err error) {
	var model system.ModelAll
	global.CMBP_DB.Model(system.ModelAll{}).Where("id = ?", modelID).First(&model)
	if model.ID == "" {
		return nil, errors.New("模型不存在或已删除")
	}
	global.CMBP_LOG.Info(fmt.Sprintf("待删除的模型是：%s", model.ModelChineseName))
	if model.ModelKind == 1 {
		// 删除已部署的大数据模型记录，删除model_info 下载信息
		var modelInfo []system.Model
		global.CMBP_DB.Model(&system.Model{}).Where("model_all_id = ?", model.ID).Find(&modelInfo)
		for _, v := range modelInfo {
			err = global.CMBP_DB.Delete(&system.DataModelConfig{}, "model_info_id = ?", v.ID).Error
			if err != nil {
				global.CMBP_LOG.Error("数据模型信息删除异常" + err.Error())
				return nil, err
			}
		}
		// 删除服务器上镜像
		_, err = utils.DeleteImages(fmt.Sprintf("%s/%s:%s", global.CMBP_CONFIG.CMBPBase.OssPath, model.ModelName, fmt.Sprintf("%s.%d", model.ModelVersion, model.Edition)))
		if err != nil {
			return nil, errors.New("镜像删除失败" + err.Error())
		}
	}

	// 删除服务器上模型文件和模型图片
	global.CMBP_LOG.Info("删除models模型目录")

	if model.NewModelFlag {
		modelFile := fmt.Sprintf(global.CMBP_CONFIG.CMBPBase.OssModelPath, model.ModelNameAndVersion())
		exists, _ := utils.PathExists(modelFile)
		if exists {
			utils.DeLFile(modelFile)
		}
		modelWareHouse := filepath.Join(global.CMBP_CONFIG.CMBPBase.ModelWareHouse, model.ModelNameAndVersion())
		existsWareHouse, _ := utils.PathExists(modelWareHouse)
		if existsWareHouse {
			utils.DeLFile(modelWareHouse)
		}
	} else {
		modelZipFile := model.ModelZipFile()
		zipExists, _ := utils.PathExists(modelZipFile)
		if zipExists {
			utils.DeLFile(modelZipFile)
		}
	}
	// 模型仓库媒体
	imgFile := model.ModelImgFile()
	videoFile := model.ModelVideoFile()
	// 模型市场媒体
	maskModelImg := filepath.Join(global.CMBP_CONFIG.CMBPBase.OssMarketModelPath, model.ModelNameAndVersion(), ".jpg")
	maskModelVideo := filepath.Join(global.CMBP_CONFIG.CMBPBase.OssMarketModelPath, model.ModelNameAndVersion(), ".mp4")

	err = utils.DeLFile(imgFile)
	if err != nil {
		return nil, err
	}
	global.CMBP_LOG.Info("删除模型图片：" + imgFile)

	err = utils.DeLFile(videoFile)
	if err != nil {
		return nil, err
	}
	global.CMBP_LOG.Info("删除模型视频: " + videoFile)

	err = utils.DeLFile(maskModelImg)
	if err != nil {
		return nil, err
	}
	global.CMBP_LOG.Info("删除模型市场图片: " + maskModelImg)
	err = utils.DeLFile(maskModelVideo)
	if err != nil {
		return nil, err
	}
	global.CMBP_LOG.Info("删除模型市场视频: " + maskModelVideo)

	var modelFeedback []string
	var modelDeploy []string
	var modelEvaluation []string
	global.CMBP_DB.Model(&system.ModelFeedback{}).Where("model_id = ?", model.ID).Pluck("model_id", &modelFeedback)
	global.CMBP_DB.Model(&system.ModelDeploy{}).Where("model_id = ?", model.ID).Pluck("model_id", &modelDeploy)
	global.CMBP_DB.Model(&system.ModelEvaluation{}).Where("model_id = ?", model.ID).Pluck("model_id", &modelEvaluation)

	allModelID := append(modelFeedback, modelDeploy...)
	allModelID = append(allModelID, modelEvaluation...)
	for _, v := range allModelID {
		err = global.CMBP_DB.Delete(&system.Event{}, "id = ?", v).Error
		if err != nil {
			return nil, err
		}
	}
	global.CMBP_DB.Delete(&model)
	return allModelID, nil
}

func (modelService *ModelService) JupyterNoteBook(userID string) (respData interface{}, err error) {
	var noteBook []system.Notebook
	global.CMBP_DB.Model(&system.Notebook{}).Where("user_id = ? AND status = 1", userID).Order("status DESC").Order("create_time DESC").Find(&noteBook)
	expireTime := 0
	jupyterList := []interface{}{}
	for _, v := range noteBook {
		expireDateTime := v.CreateTime.Add(time.Duration(*v.ExpirationTime) * time.Minute)
		now := time.Now()
		// 检查是否过期
		if !(expireDateTime.Before(now) || expireDateTime.Equal(now)) {
			expireTime = int(expireDateTime.Sub(now).Minutes())
		}
		jupyterList = append(jupyterList, map[string]interface{}{
			"notebook_name":   v.Name,
			"notebook_desc":   v.Desc,
			"create_time":     v.CreateTime.Format("2006-01-02 15:04:05"),
			"expire_time":     expireTime,
			"url":             v.URL,
			"notebook_status": v.Status,
			"uuid":            v.UUID,
			"notebook_id":     v.ID,
		})
	}
	respData = map[string]interface{}{
		"count":        len(noteBook),
		"jupyter_list": jupyterList,
	}
	return respData, nil
}

func (modelService *ModelService) Runtime(params systemReq.RunTime, role, userName string) (res interface{}, err error) {
	if params.RuntimeID != "" {
		var runTime system.RuntimeAll
		global.CMBP_DB.Model(&system.RuntimeAll{}).Where("id = ?", params.RuntimeID).First(&runTime)
		if runTime.ID == "" {
			return nil, errors.New("运行镜像不存在")
		}
		exists, err := utils.PathExists(filepath.Join(global.CMBP_CONFIG.CMBPBase.OssRuntimeLibrary, fmt.Sprintf("%s/%s.zip", runTime.Name, runTime.Tag)))
		if err != nil {
			return nil, errors.New("校验镜像OBS文件失败" + err.Error())
		}
		if !exists {
			return nil, errors.New("运行镜像丢失")
		}
		roleList := strings.Split(global.CMBP_CONFIG.CMBPBase.RuntimeDownRole, ",")
		hasPermission := false
		for _, v := range roleList {
			if v == role {
				hasPermission = true
				break
			}
		}
		downUrl := "No Permission"

		if hasPermission {
			runTimeLib := strings.Split(global.CMBP_CONFIG.CMBPBase.OssRuntimeLibrary, "/")[len(strings.Split(global.CMBP_CONFIG.CMBPBase.OssRuntimeLibrary, "/"))-1]
			if strings.ToLower(global.CMBP_CONFIG.CMBPBase.OssMode) == "minio" {
				downUrl, err = url.JoinPath(global.CMBP_CONFIG.CMBPBase.OssPath, fmt.Sprintf("/%s/%s/%s.zip", runTimeLib, runTime.Name, runTime.Tag))
				if err != nil {
					return nil, errors.New("生成镜像下载链接失败：" + err.Error())
				}
			} else if strings.ToLower(global.CMBP_CONFIG.CMBPBase.OssMode) == "obs" {
				downUrl, err = upload.HuaWeiObs.SignDownUrl(fmt.Sprintf("%s/%s/%s", runTimeLib, runTime.Name, runTime.Tag), 1800)
			}
		}

		tp := 0
		if &runTime.Type != nil {
			tp = runTime.Type
		}

		res = map[string]interface{}{
			"id":           runTime.ID,
			"name":         runTime.Name,
			"chinese_name": runTime.ChineseName,
			"description":  runTime.Description,
			"tag":          runTime.Tag,
			"file_size":    fmt.Sprintf("%0.3f", float64(utils.GetFileSize(runTime.ModelZipFilePath())/1024.0/1024.0)),
			"status":       runTime.Status,
			"upload_time":  runTime.UpdateTime,
			"down_url":     downUrl,
			"build_way":    runTime.BuildWay,
			"type":         tp,
		}
		return res, nil
	}
	relationRuntimeId := ""
	bindRuntime := []interface{}{}
	if params.ModelID != "" {
		var modelRuntime system.RuntimeModels
		global.CMBP_DB.Where("model_all_id = ?", params.ModelID).First(&modelRuntime)
		if modelRuntime.ID != "" {
			relationRuntimeId = modelRuntime.ID
			var relationRuntime system.RuntimeAll
			global.CMBP_DB.Where("id = ?", relationRuntimeId).First(&relationRuntime)
			fileSize := fmt.Sprintf("%0.3f", float64(utils.GetFileSize(relationRuntime.ModelZipFilePath())/1024.0/1024.0))
			tp := 0
			if &relationRuntime.Type != nil {
				tp = relationRuntime.Type
			}
			opFlag := 0
			if relationRuntime.Developer == userName {
				opFlag = 1
			}

			res = map[string]interface{}{
				"runtime_id":     relationRuntime.ID,
				"name":           relationRuntime.Name,
				"chinese_name":   relationRuntime.ChineseName,
				"description":    relationRuntime.Description,
				"tag":            relationRuntime.Tag,
				"file_size":      fileSize,
				"developer":      relationRuntime.Developer,
				"build_way":      relationRuntime.BuildWay,
				"status":         relationRuntime.Status,
				"type":           tp,
				"operation_flag": opFlag,
				"upload_time":    relationRuntime.UpdateTime,
				"is_relation":    1,
			}
			bindRuntime = append(bindRuntime, res)
		}
	}

	filterBy := make(map[string]interface{})

	if params.IsUsable != 0 {
		filterBy["status"] = 1
	}

	if params.Type == "null" {
		filterBy["type"] = nil
	} else {
		filterBy["type"] = strings.Trim(params.Type, " ")
	}

	filtersOr := []*gorm.DB{}

	if params.NameOrDesc != "" {
		searchTerm := "%" + params.NameOrDesc + "%"
		filtersOr = append(filtersOr, global.CMBP_DB.Where("description LIKE ?", searchTerm))
		filtersOr = append(filtersOr, global.CMBP_DB.Where("chinese_name LIKE ?", searchTerm))
		filtersOr = append(filtersOr, global.CMBP_DB.Where("name LIKE ?", searchTerm))
	}

	var query = global.CMBP_DB.Model(&system.RuntimeAll{}).Where(filterBy)

	if len(filtersOr) > 0 {
		for _, filter := range filtersOr {
			query = query.Or(filter)
		}
	}

	if relationRuntimeId != "" {
		query = query.Not("id", relationRuntimeId)
	}

	query = query.Order("update_time DESC")
	var count int64

	err = query.Count(&count).Error
	if err != nil {
		return nil, err
	}

	if params.Page != 0 && params.Limit != 0 {
		offset := (params.Page - 1) * params.Limit
		if len(bindRuntime) > 0 {
			if params.Page == 1 {
				query = query.Limit(params.Limit - 1).Offset(offset)
			} else {
				query = query.Limit(params.Limit).Offset(offset - 1)
			}
		} else {
			query = query.Limit(params.Limit).Offset(offset)
		}
	}

	var runtimes []system.RuntimeAll

	err = query.Find(&runtimes).Error
	if err != nil {
		return nil, err
	}
	for _, runtime := range runtimes {
		tp := 0
		if &runtime.Type != nil {
			tp = runtime.Type
		}
		opFlag := 0
		if runtime.Developer == userName {
			opFlag = 1
		}
		bindRuntime = append(bindRuntime, map[string]interface{}{
			"runtime_id":     runtime.ID,
			"name":           runtime.Name,
			"chinese_name":   runtime.ChineseName,
			"description":    runtime.Description,
			"tag":            runtime.Tag,
			"file_size":      fmt.Sprintf("%0.3f", float64(utils.GetFileSize(runtime.ModelZipFilePath())/1024.0/1024.0)),
			"developer":      runtime.Developer,
			"build_way":      runtime.BuildWay,
			"status":         runtime.Status,
			"type":           tp,
			"operation_flag": opFlag,
			"upload_time":    runtime.UpdateTime,
			"is_relation":    1,
		})
	}

	resData := map[string]interface{}{
		"count":      count,
		"model_list": bindRuntime,
	}
	return resData, nil
}

func (modelService *ModelService) PublishModel(params systemReq.PublishModel, userID string) (res interface{}, err error) {
	var modelAll system.ModelAll
	global.CMBP_DB.Where("id = ?", params.ModelID).Find(&modelAll)
	if modelAll.ID == "" {
		return nil, errors.New("请传递正确的模型id")
	}
	global.CMBP_LOG.Info("模型查询成功")
	aiModelFlag := false
	aiModelName := ""
	aiModelDesc := ""
	var dataModel system.DataAnalyzeModelAll
	var aiModel system.AIModelAll
	onLineFlag := 1 // 线上制作
	if modelAll.AiModelName != "" {
		if modelAll.ModelKind == 1 {
			global.CMBP_DB.Where("model_name = ?", modelAll.AiModelName).Find(&dataModel)
			if dataModel.ID != "" {
				aiModelDesc = dataModel.ModelDescription
				aiModelName = dataModel.ModelName
				if modelAll.UpdateTime.Before(*dataModel.UpdateTime) {
					aiModelFlag = true
				}
			}
		} else {
			global.CMBP_DB.Where("model_name = ?", modelAll.AiModelName).Where("id = ?", modelAll.ID).Find(&dataModel)
			if dataModel.ID != "" {
				return nil, errors.New("未找到该算法")
			}
			aiModelDesc = aiModel.ModelDescription
			aiModelName = aiModel.ModelName
			if modelAll.UpdateTime.Before(*aiModel.UpdateTime) {
				aiModelFlag = true
			}
		}
	} else {
		aiModelDesc = modelAll.ModelDescription
		onLineFlag = 0
	}
	targetPath := fmt.Sprintf("/home/tmp/%s/%s", userID, params.UUID)
	global.CMBP_LOG.Info("判断模型类型")
	modelPath := filepath.Join(global.CMBP_CONFIG.CMBPBase.ModelWareHouse, modelAll.ModelNameAndVersion())

	if modelAll.ModelKind != 1 {
		global.CMBP_LOG.Info("模型的路径：" + modelPath)
		exists, _ := utils.PathExists(modelPath)
		if exists {
			businessModelIDList := []string{}
			for _, bus := range modelAll.BusinessModelAllProxy {
				businessModelIDList = append(businessModelIDList, bus.BusinessModel.ID)
			}
			if len(businessModelIDList) > 0 {
				fileList, err := os.ReadDir(filepath.Join(modelPath, "plugins"))
				if err != nil {
					return nil, err
				}
				for _, file := range fileList {
					for _, b := range businessModelIDList {
						if strings.HasPrefix(file.Name(), "_"+b) || strings.HasPrefix(file.Name(), b) {
							// 以业务模型ID开头的文件夹，均需要删除
							err := utils.DeLFile(filepath.Join(filepath.Join(modelPath, "plugins"), file.Name()))
							if err != nil {
								return nil, err
							}
							break
						}
					}
				}
			}
			for _, soFile := range []string{"file_operation_86.so", "file_operation_arm.so"} {
				soFilePath := filepath.Join(modelPath, soFile)
				global.CMBP_LOG.Info("删除AIMonitorEnd的so文件" + soFilePath)
				err = utils.DeLFile(soFilePath)
				if err != nil {
					global.CMBP_LOG.Error("删除AIMonitorEnd下的so文件失败" + soFilePath)
				}
			}
		} else {
			utils.Unzip(modelAll.ModelZipFile(), targetPath)
		}
	}
	// 大数据和图像分析模型都统一一个复制方法
	err = utils.CopyDir(modelPath, targetPath)
	if err != nil {
		return nil, errors.New("拷贝model_warehouse到tmp目录失败" + err.Error())
	}
	err = utils.CopyEnd(targetPath, "", "")
	if err != nil {
		return nil, errors.New("拷贝End到tmp目录失败" + err.Error())
	}
	global.CMBP_LOG.Info("复制完成model_warehouse，AIMonitorEnd到tmp目录")

	var modelType system.ModelType
	global.CMBP_DB.Where("model_type = ?", modelAll.ModelType).Find(&modelType)

	//if modelAll.User == "" {
	//
	//}
	var user system.Users
	global.CMBP_DB.Where("id = ?", modelAll.User).Find(&user)

	businessList := []interface{}{}
	businessNameList := []interface{}{}
	var busParams map[string]interface{}
	for _, b := range modelAll.BusinessModelAllProxy {
		err = json.Unmarshal([]byte(b.BusinessModel.BusinessParams), &busParams)
		if err != nil {
			return nil, errors.New("业务模型参数business_params提取失败" + err.Error())
		}
		if len(busParams) > 0 {
			for _, k := range busParams {
				businessNameList = append(businessNameList, k)
			}
		}
	}
	var businessType map[string]interface{}
	if modelAll.BusinessType != "" {
		err = json.Unmarshal([]byte(modelAll.BusinessType), &businessType)
		if err != nil {
			return nil, errors.New("业务模型参数BusinessType提取失败" + err.Error())
		}
	} else {
		businessType = make(map[string]interface{})
	}
	var businessApi map[string]interface{}
	err = json.Unmarshal([]byte(modelAll.BusinessAPI), &businessApi)
	if err != nil {
		return nil, errors.New("业务模型参数BusinessAPI提取失败" + err.Error())
	}
	businessClass := map[string]interface{}{}
	if modelAll.BusinessClass != "" {
		err = json.Unmarshal([]byte(modelAll.BusinessClass), &businessClass)
		if err != nil {
			return nil, errors.New("业务模型参数BusinessClass提取失败" + err.Error())
		}
	}
	if modelAll.BusinessParams != "" {
		var modelAllBusParams map[string]interface{}
		_ = json.Unmarshal([]byte(modelAll.BusinessParams), &modelAllBusParams)
		for k, v := range modelAllBusParams {
			business := map[string]interface{}{}
			for _, bName := range businessNameList {
				if bName == k {
					// 如果业务模型是从业务模型库关联的，则返回给前端的时候，去掉关联的业务模型名称和参数信息
					continue
				}
				business["business_name"] = k
				business["business_params"] = v
			}
			for _, i := range v.([]map[string]interface{}) {
				for vk, _ := range i {
					if !strings.Contains("paras_text", vk) {
						i["paras_text"] = i["paras_desc"]
					}
				}
			}
			business["business_type"] = []interface{}{}
			for _, btype := range businessType {
				if btype == k {
					business["business_type"] = businessType[k]
				}
			}
			bAPI := map[string]interface{}{}
			if businessApi[k] != nil {
				bAPI = businessApi[k].(map[string]interface{})
			}
			bClass := []interface{}{}
			if businessClass[k] != nil {
				bClass = businessClass[k].([]interface{})
			}
			business["business_apis"] = bAPI
			business["business_classes"] = bClass
			businessList = append(businessList, business)
		}
	} else if businessType != nil {
		for k, v := range businessType {
			business := map[string]interface{}{}
			business["business_name"] = k
			business["business_type"] = v
			business["business_class"] = []interface{}{}
			business["business_apis"] = map[string]interface{}{}
			business["business_params"] = []interface{}{}
		}
	}
	global.CMBP_LOG.Info("提取业务模型参数")
	var hardArch system.HardwareArch
	var runtimeModel system.RuntimeModels
	global.CMBP_DB.Where("code = ?", params.ModelID).Find(&hardArch)
	global.CMBP_DB.Where("model_all_id = ?", params.ModelID).Find(&runtimeModel)
	runtimeID := ""
	if runtimeModel.ID != "" {
		runtimeID = runtimeModel.ID
	}
	var industryName, modelKindName, algorithmName string
	global.CMBP_DB.Model(&system.Industry{}).Where("industry_code = ?", modelAll.IndustryCode).Pluck("industry_name", &industryName)
	global.CMBP_DB.Model(&system.ModelKind{}).Where("model_kind = ?", modelAll.ModelKind).Pluck("kind_name", &modelKindName)
	global.CMBP_DB.Model(&system.AlgorithmInfo{}).Where("algorithm_id = ?", modelAll.AlgorithmID).Pluck("algorithm_name", &algorithmName)

	modelChinessName := modelAll.ModelChineseName
	chinessName := strings.Split(modelAll.ModelChineseName, "-")
	if len(chinessName) > 2 {
		modelChinessName = strings.Join(chinessName[:len(chinessName)-2], "-")
	}

	resData := map[string]interface{}{
		"model_id":        modelAll.ID,
		"industry_code":   modelAll.IndustryCode,
		"industry_name":   industryName,
		"model_kind":      modelAll.ModelKind,
		"model_kind_name": modelKindName,
		"algorithm_id":    modelAll.AlgorithmID,
		"algorithm_name":  algorithmName,
		"build_way":       modelAll.BuildWay,
		"model_field":     modelType.ModelField.Name,
		"model_big_type":  modelType.ModelField.Code,

		"model_type_desc":         modelType.ModelTypeDesc,
		"model_type":              modelType.ModelType,
		"model_name":              modelAll.ModelName,
		"model_chinese_name":      modelChinessName,
		"model_version":           modelAll.ModelVersion,
		"ai_model_description":    aiModelDesc,
		"model_description":       modelAll.ModelDescription,
		"technical_description":   modelAll.TechnicalDescription,
		"performance_description": modelAll.PerformanceDescription,
		"hardware_type":           hardArch.Name,
		"hardware_type_code":      hardArch.Code,
		"is_image": func() string {
			if modelAll.IsImage {
				return "1"
			} else {
				return "0"
			}
		},
		"need_gpu": func() string {
			if modelAll.NeedGPU {
				return "1"
			} else {
				return "0"
			}
		},
		"is_real_channel": func() string {
			if modelAll.IsRealChannel == "" {
				return "1"
			} else {
				return "0"
			}
		},
		"ai_model_name": func() string {
			if aiModelName != "" {
				return aiModelName
			} else {
				return modelAll.ModelName
			}
		},
		"cmd":      modelAll.CMD,
		"json_url": modelAll.JsonURL,
		"img_url":  modelAll.ImgURL,
		"on_boot": func() string {
			if modelAll.OnBoot {
				return "1"
			} else {
				return "0"
			}
		},
		"user":          modelAll.User,
		"developer":     user.Username, // 模型开发上传人员
		"upload_time":   modelAll.UpdateTime,
		"business_dict": businessList,
		"test_duration": modelAll.TestDuration,
		"accuracy":      modelAll.Accuracy,
		"runtime_id": func() *string {
			if modelAll.IsImage {
				return &runtimeID
			} else {
				return nil
			}
		},
		"business_model_list": func() []map[string]interface{} {
			var busRes []map[string]interface{}
			for _, b := range modelAll.BusinessModelAllProxy {
				busRes = append(busRes, map[string]interface{}{
					"id":   b.BusinessModel.ID,
					"name": b.BusinessModel.ModelName,
				})
			}
			return busRes
		}, //所选的业务模型id
		"img_path": filepath.Join("/models/models/", modelAll.ModelNameAndVersion(), ".jpg"),
		"video_path": func() string {
			exists, _ := utils.PathExists(modelAll.ModelVideoFile())
			if exists {
				return filepath.Join("/models/models", modelAll.ModelNameAndVersion(), ".mp4")
			} else {
				return ""
			}
		},

		"ai_model_flag":   aiModelFlag,
		"on_line_flag":    onLineFlag,
		"advantage":       modelAll.Advantage,
		"model_purpose":   modelAll.ModelPurpose,
		"is_relate_pangu": 0,
		"pangu_model_id":  "",
	}
	var panguBigModel system.ClassRelatePangu
	global.CMBP_DB.Where("model_id = ?", params.ModelID).Find(&panguBigModel)
	//if panguBigModel.ID
	return resData, nil
}
