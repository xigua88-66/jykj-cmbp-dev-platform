package service

import (
	"jykj-cmbp-dev-platform/server/service/example"
	"jykj-cmbp-dev-platform/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
