package tests

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/guregu/null"
	"testing"
	"time"
	"vv/orm"
	"vv/orm/tests/model"
)

// TestUpdate 测试更新单条数据
func TestUpdateItem(t *testing.T) {
	q := orm.NewSession()

	t.Run("UpdateOne", func(t *testing.T) {
		var item model.UnitTestModel
		q.GetItemWhereFirst(&item,"node_uuid = ?","NodeUUID10")
		if item.NodeUUID != "NodeUUID10" {
			t.Error("没有发现记录")
		}

		item.WorkflowUUID = "wf_uuid"
		item.OperaterName = null.StringFrom("opt_name")
		q.Update(&item,[]string{"workflow_uuid","operater_name"})
		if q.Error != nil {
			t.Error(q.Error)
		}
		if q.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}

		var item2 model.UnitTestModel
		q.GetItemWhereFirst(&item2,"node_uuid = ?","NodeUUID10")
		if item2.WorkflowUUID != "wf_uuid" {
			t.Error("更新失败")
		}
	})
}

// TestBatchUpdate 测试批量更新
func TestBatchUpdate(t *testing.T) {
	q := orm.NewSession()

	var items []model.UnitTestModel
	q.GetItemWhere(&items,"id < 3")

	t.Run("BatchUpdate", func(t *testing.T) {
		for idx := range items {
			//items[idx].UUID = "xxxx"
			items[idx].DefaultNode = "DN" + cast.ToString(idx)
			items[idx].UpdatedAt = time.Now()
			items[idx].OperaterTime = null.Time{}
			items[idx].OperaterName = null.StringFrom("BatchUpdate"+cast.ToString(idx))
		}

		q.BatchUpdate(&items,30,[]string{"default_node","operater_name","updated_at","operater_time"})
		if q.Error != nil {
			t.Error(q.Error)
		}
		if q.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}

		var checkItem model.UnitTestModel
		qc := orm.NewSession()
		qc.GetItemWhereFirst(&checkItem,"default_node = ?","DN0")
		if checkItem.DefaultNode != "DN0" {
			t.Errorf("更新失败")
		}
		fmt.Printf("BatchUpdate Result:%#+v",checkItem)
	})
}