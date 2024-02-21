package request

type SysAuthorityBtnReq struct {
	MenuID      string   `json:"menuID"`
	AuthorityId string   `json:"authorityId"`
	Selected    []string `json:"selected"`
}
