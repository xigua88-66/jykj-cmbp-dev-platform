package request

import (
	"jykj-cmbp-dev-platform/server/model/common/request"
	"jykj-cmbp-dev-platform/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}

type AddFrontOpsLog struct {
	EnterTime string `json:"enter_time" binding:"required"`
	LeaveTime string `json:"leave_time" binding:"required"`
	PageName  string `json:"page_name" binding:"required"`
}
