package response

import (
	"jykj-cmbp-dev-platform/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
