package response

import "jykj-cmbp-dev-platform/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
