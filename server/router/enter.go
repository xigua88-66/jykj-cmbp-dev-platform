package router

import (
	"jykj-cmbp-dev-platform/server/router/example"
	"jykj-cmbp-dev-platform/server/router/system"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
