package tests

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/spf13/cast"
	"time"
	"github.com/wintrysky/orm/tests/model"
)

func BuildRecord(idx int,operateType string) model.UnitTestModel{
	var item model.UnitTestModel
	item.CategoryID = 1
	item.UUID = uuid.New().String()
	item.NodeUUID = "NodeUUID" + cast.ToString(idx)
	item.OperaterName = null.StringFrom("admin")
	item.TypeID = 5
	item.OperaterTime = null.TimeFrom(time.Now())
	item.WorkflowUUID = "WorkflowUUID" + cast.ToString(idx)
	item.DefaultNode = "Y"
	item.OperaterID = 1001
	item.OperateType = operateType
	item.RatioType = 1
	item.RatioUUID = "RatioUUID" + cast.ToString(idx)

	return item
}
