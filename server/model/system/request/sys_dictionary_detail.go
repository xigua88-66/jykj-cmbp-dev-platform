package request

import (
	"jykj-cmbp-dev-platform/server/model/common/request"
	"jykj-cmbp-dev-platform/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
