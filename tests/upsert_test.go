package tests

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/guregu/null"
	"testing"
	"time"
	"vv/orm"
	"vv/orm/tests/model"
)

// TestBatchUpsert 测试批量新增或更新
func TestBatchUpsert(t *testing.T) {
	qb := orm.NewSession()

	t.Run("BatchUpsert", func(t *testing.T) {
		var items []model.UnitTestModel
		uid := "xxxxx01"
		for i := 80; i < 120; i++ {
			var item model.UnitTestModel
			item.ID = cast.ToInt64(i)
			if i > 100 {
				item.ID = 0
			}
			item.CategoryID = 1
			item.UUID = uuid.New().String()
			if i == 80 || i == 119 {
				item.UUID = uid
			}
			item.NodeUUID = "default node"
			item.OperaterName = null.StringFrom("admin")
			item.TypeID = 5
			item.OperaterTime = null.TimeFrom(time.Now())
			item.WorkflowUUID = "workflow uuid"
			item.DefaultNode = "Y"
			item.OperaterID = cast.ToInt64(i)
			item.OperateType = "batch_upsert"
			item.RatioType = 111
			item.RatioUUID = "ratio uuid"
			items = append(items, item)
		}
		fmt.Println("=====BatchUpsert===========================================")
		qb.BatchUpsert(&items, 10,[]string{"operate_type","ratio_type","uuid"},[]string{"uuid"})
		fmt.Println("=====BatchUpsert===========================================")
		if qb.Error != nil {
			t.Error(qb.Error)
		}
		if qb.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}

		// 校验
		var checkItems []model.UnitTestModel
		qb.GetItemWhere(&checkItems,"operate_type = ?","batch_upsert")
		if len(checkItems) == 0 {
			t.Errorf("BatchUpsert失败")
		}
	})
}
