package request

import (
	"jykj-cmbp-dev-platform/server/model/common/request"
	"jykj-cmbp-dev-platform/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
