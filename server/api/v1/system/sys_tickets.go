package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/common/response"
	"jykj-cmbp-dev-platform/server/model/system"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	systemRsp "jykj-cmbp-dev-platform/server/model/system/response"
	"jykj-cmbp-dev-platform/server/utils"
	"math"
	"time"
)

type TicketsApi struct {
}

var FEEDBACK_SCHEDULE = map[int]int{
	78: 0,
	79: 10,
	80: 20,
	82: 30,
	83: 40,
	84: 50,
	85: 60,
	86: 70,
	87: 80,
	88: 90,
	89: 100,
}

var DEVELOP_SCHEDULE = map[int]int{
	69: 0,
	70: 15,
	71: 30,
	72: 45,
	73: 60,
	74: 75,
	77: 90,
	76: 100,
}

func (t *TicketsApi) GetTickets(c *gin.Context) {
	var req systemReq.GetUserTickets
	//err := c.ShouldBindJSON(&req)
	err := c.BindQuery(&req)
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	err = utils.Verify(req, utils.GetUserTicketsVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userName := utils.GetUserName(c)
	userId := utils.GetUserID(c)
	var roleName string
	err = global.CMBP_DB.Model(system.Roles{}).Joins("JOIN t_user_roles on t_roles_info.id=t_user_roles.role_id").Where("t_user_roles.user_id = ?", userId).Pluck("t_roles_info.name", &roleName).Error
	if err != nil {
		return
	}
	ticketsByte, err := ticketService.QueryTickets(userName, roleName, req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var ticketList systemRsp.UserTicketList
	json.Unmarshal(ticketsByte, &ticketList)
	fmt.Println(ticketList)
	for i := range ticketList.TicketList {
		ticketInfo := &ticketList.TicketList[i] // 获取当前ticket的指针，以便修改

		var diff time.Duration
		if ticketInfo.StateID == 45 { // 假设config是加载的配置
			gmtModified, _ := time.Parse("2006-01-02 15:04:05", ticketInfo.GmtModified)
			gmtCreated, _ := time.Parse("2006-01-02 15:04:05", ticketInfo.GmtCreated)
			diff = gmtModified.Sub(gmtCreated)
		} else {
			gmtCreated, _ := time.Parse("2006-01-02 15:04:05", ticketInfo.GmtCreated)
			diff = time.Now().Sub(gmtCreated)
		}

		dDays := int(diff.Hours() / 24)
		dHours := math.Ceil(diff.Hours()) - float64(dDays*24)

		duration := ""
		if dDays > 0 {
			duration += fmt.Sprintf("%d天", dDays)
		}
		if dHours > 0 {
			duration += fmt.Sprintf("%d小时", int(dHours))
		}
		ticketInfo.Duration = duration
		var ticketRecord system.TicketRecord
		result := global.CMBP_DB.Model(&system.TicketRecord{}).Where("ticket_id = ?", ticketInfo.ID).First(&ticketRecord)

		if ticketInfo.WorkflowID == 18 {
			if (result.Error != nil) && errors.Is(result.Error, gorm.ErrRecordNotFound) {
				ticketInfo.DatasetID = ""
				ticketInfo.InpaintTaskID = ""
				ticketInfo.QualityInspectID = ""
				ticketInfo.ModelName = ""
			} else {
				schedule, exists := FEEDBACK_SCHEDULE[ticketInfo.StateID]
				if exists {
					ticketInfo.Schedule = fmt.Sprintf("%d%%", schedule)
				}
				ticketInfo.DatasetID = ticketRecord.DatasetID
				ticketInfo.InpaintTaskID = ticketRecord.InpaintTaskID
				ticketInfo.QualityInspectID = ticketRecord.QualityInspectID
				var modelChineseName = ""
				global.CMBP_DB.Model(&system.ModelAll{}).Where("id =?", ticketRecord.ModelID).Pluck("model_chinese_name", &modelChineseName)
				ticketInfo.ModelName = modelChineseName

			}
		} else if ticketInfo.WorkflowID == 17 {
			schedule, exists := DEVELOP_SCHEDULE[ticketInfo.StateID]
			if exists {
				ticketInfo.Schedule = fmt.Sprintf("%d%%", schedule)
			}
			if ticketRecord.Ident != "" {
				ticketInfo.IdentInfo = ticketRecord.Ident
			} else {
				ticketInfo.IdentInfo = ""
			}
		}
	}
	if ticketList.TicketList == nil {
		ticketList.TicketList = []systemRsp.UserTicketResponse{}
	}
	response.OkWithDetailed(ticketList, "工单获取成功", c)
}
