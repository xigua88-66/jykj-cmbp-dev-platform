package initialize

import (
	_ "jykj-cmbp-dev-platform/server/source/example"
	_ "jykj-cmbp-dev-platform/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
