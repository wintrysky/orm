package tests

import (
	"fmt"
	"testing"
	"vv/orm"
	"vv/orm/tests/model"
)

// TestInsert 测试新增单条数据
func TestInsert(t *testing.T) {
	q := orm.NewSession()

	t.Run("InsertOne", func(t *testing.T) {
		item := BuildRecord(1,"InsertOne")
		q.Insert(&item)
		if q.Error != nil {
			t.Error(q.Error)
		}
		if q.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		if item.ID < 1 {
			t.Errorf("插入数据失败")
		}
	})

	t.Run("InsertOneWithField", func(t *testing.T) {
		item := BuildRecord(1,"InsertOneWithField")
		q.Insert(&item,[]string{"uuid","operate_type","operater_time"})
		if q.Error != nil {
			t.Error(q.Error)
		}
		if q.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		if item.ID < 1 {
			t.Errorf("插入数据失败")
		}
		fmt.Println("插入数据，ID：",item.ID)
	})
}

// TestBatchInsert 测试批量新增
func TestBatchInsert(t *testing.T) {
	qb := orm.NewSession()

	t.Run("BatchInsert", func(t *testing.T) {
		var items []model.UnitTestModel
		for i := 102; i < 201; i++ {
			item := BuildRecord(i,"BatchInsert")
			//item.UUID="xxx"
			items = append(items, item)
		}

		qb.BatchInsert(&items, 500)
		if qb.Error != nil {
			t.Error(qb.Error)
		}
		if qb.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}

		if items[0].ID < 1 {
			t.Errorf("插入数据失败")
		}
		fmt.Println("BatchInsert,ID：",items[0].ID)
	})

	t.Run("BatchInsertWithSlice", func(t *testing.T) {
		var items []model.UnitTestModel
		for i := 300; i < 310; i++ {
			item := BuildRecord(i,"BatchInsertWithSlice")
			items = append(items, item)
		}

		qb.BatchInsert(&items, 3,[]string{"uuid","operate_type","node_uuid"})
		if qb.Error != nil {
			t.Error(qb.Error)
		}
		if qb.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		if items[0].ID < 1 {
			t.Errorf("插入数据失败")
		}

		// 校验
		var checkItems []model.UnitTestModel

		qb.GetItemWhere(&checkItems,"operate_type = ?","BatchInsertWithSlice")
		if len(checkItems) != 10 {
			t.Errorf("BatchInsertWithSlice插入数据失败")
		}

		fmt.Println("BatchInsertCheck,Items：",len(checkItems))
	})
}