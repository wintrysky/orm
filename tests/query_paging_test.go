package tests

import (
	"fmt"
	"testing"
	"github.com/wintrysky/orm"
	"github.com/wintrysky/orm/tests/model"
)

func TestPagingQuery(t *testing.T) {
	f := orm.NewSession()

	t.Run("PagingQueryWhere", func(t *testing.T) {
		var items []model.UnitTestModel
		cnt,err := f.PagingQueryWhere(&items, 10,2,"id desc",
			"operate_type = ? and id > ?", "init data",10)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Println("PagingQueryWhere Items:", cnt)
		fmt.Printf("PagingQueryWhere Result:%#+v",items[0])
	})

	t.Run("PagingQueryByCondition", func(t *testing.T) {
		var items []model.UnitTestModel
		var conns []orm.Condition
		var con1 orm.Condition
		con1.MatchType = orm.Fuzzy
		con1.Field = "operate_type"
		con1.Value = "init"

		var con2 orm.Condition
		con2.MatchType = orm.Larger
		con2.Field = "operater_id"
		con2.Value = 50

		conns = append(conns,con1)
		conns = append(conns,con2)

		cnt,err := f.PagingQueryByCondition(&items, 10,2,"id desc", conns)
		if err != nil {
			t.Error(err.Error())
		}
		fmt.Println("PagingQueryByCondition Items:", cnt)
		fmt.Printf("PagingQueryByCondition Result:%#+v",items[0])
	})
}
