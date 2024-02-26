package request

type GetUserTickets struct {
	SN                  string `form:"sn"`
	Title               string `form:"title"`
	Creator             string `form:"creator"`
	CreateStart         string `form:"create_start"`
	CreatorEnd          string `form:"creator_end"`
	WorkflowIds         string `form:"workflow_ids"`
	StateIds            *int   `form:"state_ids"`
	TicketIds           string `form:"ticket_ids"`
	Reverse             string `form:"reverse"`
	Page                int    `form:"page"`
	PerPage             string `form:"per_page"`
	ActStateId          int    `form:"act_state_id"`
	ParentTicketId      int    `form:"parent_ticket_id"`
	ParentTicketStateId int    `form:"parent_ticket_state_id"`
	Category            string `form:"category"`
	CoalMineName        string `form:"coal_mine_name"`
	Ident               string `form:"ident"`
	StateId             *int   `form:"state_id"`
	OpType              int    `form:"op_type"`
	QueryField          string `form:"query_field"`
	QueryValue          string `form:"query_value"`
	OrderBy             string `form:"order_by"`
}
