package tests

import (
	"fmt"
	"testing"
	"github.com/wintrysky/orm"
	"github.com/wintrysky/orm/tests/model"
)

func TestQuery(t *testing.T) {

	f := orm.NewSession()

	t.Run("GetItemWhere", func(t *testing.T) {
		var items []model.UnitTestModel
		f.GetItemWhere(&items, "operate_type = ?", "init data")
		if f.Error != nil {
			t.Error(f.Error)
		}
		if f.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		fmt.Println("GetItemWhere Result:", len(items))
	})

	t.Run("GetItemWhereFirst", func(t *testing.T) {
		var item model.UnitTestModel
		f.GetItemWhereFirst(&item, "operate_type = ?", "init data")
		if f.Error != nil {
			t.Error(f.Error)
		}
		if f.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		fmt.Println("GetItemWhereFirst Result:", item.UUID)
	})

	t.Run("GetItemByCondition", func(t *testing.T) {
		var items []model.UnitTestModel

		var conns []orm.Condition
		var con1 orm.Condition
		con1.MatchType = orm.Fuzzy
		con1.Field = "operate_type"
		con1.Value = "init"

		var con2 orm.Condition
		con2.MatchType = "" // accurate 精确匹配
		con2.Field = "ratio_uuid"
		con2.Value = "RatioUUID9"

		conns = append(conns,con1)
		conns = append(conns,con2)

		f.GetItemByCondition(&items,conns)
		if f.Error != nil {
			t.Error(f.Error)
		}
		if f.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		if len(items) == 0 {
			t.Error("没有发现数据")
		}

		fmt.Println("GetItemByCondition Items:", len(items))
		fmt.Println("GetItemByCondition Result:", items[0].UUID)
	})
}

func TestQueryRaw(t *testing.T) {
	f := orm.NewSession()

	t.Run("ExecuteTextQuery", func(t *testing.T) {
		var items []model.UnitTestModel
		f.ExecuteTextQuery(&items, "select * from unit_test_model where operate_type = ? and id > ?",
			"init data",30)
		if f.Error != nil {
			t.Error(f.Error)
		}
		if f.RowsAffected == 0 {
			t.Error("RowsAffected is zero")
		}
		fmt.Println("ExecuteTextQuery Items:", len(items))
		fmt.Printf("ExecuteTextQuery Result:%#+v", items[0])
	})

}