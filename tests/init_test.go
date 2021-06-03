package tests

import (
	"fmt"
	"github.com/hyahm/goconfig"
	"testing"
	"vv/orm"
	"vv/orm/internal"
	"vv/orm/tests/model"
)

type baseConfig struct {
	DBSettings orm.DBSettings
}

func readBaseConfig(bconfig *baseConfig, confFile string) {
	goconfig.InitConf(confFile, goconfig.INI)

	bconfig.DBSettings.DBName = goconfig.ReadString("DBSettings.DBName", "")
	bconfig.DBSettings.Port = goconfig.ReadInt("DBSettings.Port")
	bconfig.DBSettings.User = goconfig.ReadString("DBSettings.User", "")
	bconfig.DBSettings.Host = goconfig.ReadString("DBSettings.Host", "")
	bconfig.DBSettings.Password = goconfig.ReadString("DBSettings.Password", "")
}

func TestMain(m *testing.M) {
	baseConfig := baseConfig{}
	readBaseConfig(&baseConfig, "./config.ini")

	err := orm.InitDB(baseConfig.DBSettings)
	internal.ThrowError(err)

	createTable()
	cleanData()
	initData()

	m.Run()

	//cleanData()
}

// 创建数据库测试表
func createTable() {
	fmt.Println("初始化表：UnitTestModel")
	var newTable model.UnitTestModel
	f := orm.NewSession()
	err := f.AutoMigrate(&newTable)
	internal.ThrowError(err)
}

// 初始化数据
func initData() {
	var items []model.UnitTestModel

	for i := 1; i < 101; i++ {
		item := BuildRecord(i,"init data")
		items = append(items, item)
	}

	qb := orm.NewSession()
	qb.BatchInsert(&items, 500)
}

// 清除数据
func cleanData() {
	qa := orm.NewSession()

	qa.Execute("truncate table unit_test_model")
}