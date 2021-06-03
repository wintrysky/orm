package model

import (
	"github.com/guregu/null"
	"time"
)

// UnitTestModel 测试表
type UnitTestModel struct {
	CategoryID   int64     `gorm:"column:category_id" json:"category_id"`
	DefaultNode  string    `gorm:"column:default_node;size:100" json:"default_node"`
	ID           int64     `gorm:"primary_key;column:id" json:"id"`
	NodeUUID     string    `gorm:"column:node_uuid;size:100" json:"node_uuid"`
	OperateType  string    `gorm:"column:operate_type;size:100" json:"operate_type"`
	OperaterID   int64     `gorm:"column:operater_id" json:"operater_id"`
	OperaterName null.String    `gorm:"column:operater_name;size:100" json:"operater_name"`
	OperaterTime null.Time `gorm:"column:operater_time" json:"operater_time"`
	RatioType    int       `gorm:"column:ratio_type" json:"ratio_type"`
	RatioUUID    string    `gorm:"column:ratio_uuid;size:100" json:"ratio_uuid"`
	TypeID       int64     `gorm:"primary_key;column:type_id" json:"type_id"`
	UUID         string    `gorm:"column:uuid;size:100;unique;not null" json:"uuid"`
	WorkflowUUID string    `gorm:"column:workflow_uuid;size:100" json:"workflow_uuid"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (s *UnitTestModel) TableName() string {
	return "unit_test_model"
}