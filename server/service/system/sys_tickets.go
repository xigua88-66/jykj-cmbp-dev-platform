package system

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"jykj-cmbp-dev-platform/server/global"
	systemReq "jykj-cmbp-dev-platform/server/model/system/request"
	"jykj-cmbp-dev-platform/server/utils"
	"strconv"
	"time"
)

type TicketService struct {
}

func (ticketService *TicketService) QueryTickets(userName string, roleName string, req systemReq.GetUserTickets) (resp []byte, err error) {
	reqUrl := global.CMBP_CONFIG.CMBPBase.WorkFlowUrl
	header := ticketService.GetHeader()
	header["username"] = userName
	fmt.Println(userName)
	if req.Category == "view" || req.Category == "relation" {
		req.Category = "all"
	}
	if req.CoalMineName != "" {
		req.QueryField = "enterprise_info"
		req.QueryValue = req.CoalMineName
	} else if req.StateId != nil {
		req.StateIds = req.StateId
	}
	reqJson, _ := json.Marshal(req)
	result, err := utils.HttpService(reqUrl, "GET", reqJson, header)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil

}

func (ticketService *TicketService) GetHeader() map[string]string {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	oriStr := timeStamp + global.CMBP_CONFIG.CMBPBase.WorkFlowAppSk
	singnture := fmt.Sprintf("%x", md5.Sum([]byte(oriStr)))
	headers := map[string]string{
		"signature": singnture,
		"timestamp": timeStamp,
		"appname":   global.CMBP_CONFIG.CMBPBase.WorkFlowAppName,
	}
	return headers
}
