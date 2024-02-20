package v1

import (
	"jykj-cmbp-dev-platform/server/api/v1/example"
	"jykj-cmbp-dev-platform/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
