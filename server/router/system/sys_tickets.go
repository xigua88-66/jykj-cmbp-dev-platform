package system

import (
	"github.com/gin-gonic/gin"
	v1 "jykj-cmbp-dev-platform/server/api/v1"
)

type TicketRouter struct {
}

func (s *TicketRouter) InitTicketRouter(Router *gin.RouterGroup) {

	//ticketsRouter := Router.Group("tickets")
	ticketsRouter := Router
	ticketsApi := v1.ApiGroupApp.SystemApiGroup.TicketsApi
	{
		ticketsRouter.GET("/v1.9/tickets", ticketsApi.GetTickets) // 管理员注册账号
	}
}
