package tests

import (
	"testing"
	"github.com/wintrysky/orm"
)

// TestTxInsert 事务提交
func TestTxInsert(t *testing.T) {
	tx := orm.BeginTransaction()
	defer tx.EndTransaction()

	item := BuildRecord(1,"TxInsert")
	tx.Insert(&item)
	//panic("------")

	item2 := BuildRecord(1,"TxInsert")
	tx.Insert(&item2,[]string{"uuid","node_uuid","operate_type"})
	if tx.Error != nil {
		t.Error(tx.Error)
	}
	if tx.RowsAffected == 0 {
		t.Error("RowsAffected is zero")
	}

}