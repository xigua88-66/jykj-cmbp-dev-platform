package system

import "time"

type TicketRecord struct {
	ID               string     `gorm:"primaryKey;type:string;size:32;default:generate_uuid"`
	TicketCreator    string     `gorm:"type:string;size:32;not null"`
	FeedbackType     *int       `gorm:"type:int;"`
	WorkflowID       int        `gorm:"type:int;not null"`
	TicketID         int        `gorm:"type:int;not null"`
	VideoID          string     `gorm:"type:string;size:64;not null"`
	EventID          string     `gorm:"type:string;size:64;not null"`
	ProductLine      string     `gorm:"type:text;not null"`
	IsUpgraded       int        `gorm:"type:int;not null;default:0"`
	Ident            string     `gorm:"type:text;not null"`
	ContractType     int        `gorm:"type:int;not null"`
	InpaintTaskID    string     `gorm:"type:string;size:64;not null"`
	DatasetID        string     `gorm:"type:string;size:64;not null"`
	QualityInspectID string     `gorm:"type:string;size:64;not null"`
	ModelID          *string    `gorm:"type:string;size:32;"`
	WeightsFile      *string    `gorm:"type:text;"`
	ModelCreator     string     `gorm:"type:string;size:32;not null"`
	MineName         *string    `gorm:"type:string;size:200;"`
	Status           int        `gorm:"type:int;not null"`
	Extra            *string    `gorm:"type:text;"`
	StartTime        time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	EndTime          *time.Time `gorm:"type:datetime;"`
	Duration         *int       `gorm:"type:int;"`
	CreateTime       time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdateTime       time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP;autoUpdateTime"`
}

func (TicketRecord) TableName() string {
	return "t_ticket_record"
}
