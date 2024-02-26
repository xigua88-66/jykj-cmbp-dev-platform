package response

type UserTicketList struct {
	TicketList []UserTicketResponse `json:"value"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"per_page"`
	Total      int                  `json:"total"`
}

type UserTicketResponse struct {
	ActStateID          int             `json:"act_state_id"`
	AddNodeMan          string          `json:"add_node_man"`
	Creator             string          `json:"creator"`
	CreatorInfo         CreatorInfo     `json:"creator_info"`
	DatasetID           string          `json:"dataset_id"`
	Duration            string          `json:"duration"`
	GmtCreated          string          `json:"gmt_created"`
	GmtModified         string          `json:"gmt_modified"`
	ID                  int             `json:"id"`
	InAddNode           bool            `json:"in_add_node"`
	InpaintTaskID       string          `json:"inpaint_task_id"`
	IsDeleted           bool            `json:"is_deleted"`
	ModelName           string          `json:"model_name"`
	MultiAllPerson      string          `json:"multi_all_person"`
	ParentTicketID      int             `json:"parent_ticket_id"`
	ParentTicketStateID int             `json:"parent_ticket_state_id"`
	Participant         string          `json:"participant"`
	ParticipantInfo     ParticipantInfo `json:"participant_info"`
	ParticipantTypeID   int             `json:"participant_type_id"`
	QualityInspectID    string          `json:"quality_inspect_id"`
	Relation            string          `json:"relation"`
	Schedule            string          `json:"schedule"`
	ScriptRunLastResult bool            `json:"script_run_last_result"`
	Sn                  string          `json:"sn"`
	State               State           `json:"state"`
	StateID             int             `json:"state_id"`
	Title               string          `json:"title"`
	WorkflowID          int             `json:"workflow_id"`
	WorkflowInfo        WorkflowInfo    `json:"workflow_info"`
	IdentInfo           string          `json:"ident_info"`
}

type CreatorInfo struct {
	Alias    string      `json:"alias"`
	DeptInfo interface{} `json:"dept_info"` // 使用interface{}因为dept_info为空对象{}
	Email    string      `json:"email"`
	IsActive bool        `json:"is_active"`
	Phone    string      `json:"phone"`
	Username string      `json:"username"`
}

type ParticipantInfo struct {
	Participant         string `json:"participant"`
	ParticipantAlias    string `json:"participant_alias"`
	ParticipantName     string `json:"participant_name"`
	ParticipantTypeID   int    `json:"participant_type_id"`
	ParticipantTypeName string `json:"participant_type_name"`
}

type State struct {
	StateID    int         `json:"state_id"`
	StateLabel interface{} `json:"state_label"` // 使用interface{}因为state_label为空对象{}
	StateName  string      `json:"state_name"`
}

type WorkflowInfo struct {
	WorkflowID   int    `json:"workflow_id"`
	WorkflowName string `json:"workflow_name"`
}
