package orm

import (
	"fmt"
	"github.com/guregu/null"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"reflect"
	"strings"
	"time"
	"github.com/wintrysky/orm/internal"
)

func (x *GormDB) batchUpdateInternal(entities interface{}, cols []string, rowMap []map[string]interface{}) (string,
	[]interface{}) {

	// DryRun模式，不执行，取结构体
	stmt := x.db.Session(&gorm.Session{DryRun: true}).First(entities).Statement
	// TODO 更新updated_at
	// 判断是否需要更新updated_at
	//var hasUpdatedColumnTag bool
	//for _, field := range stmt.Statement.Statement.Schema.Fields {
	//	if field.DBName == internal.UpdatedAt {
	//		hasUpdatedColumnTag = true
	//	}
	//}
	//var hasUpdatedColumn bool
	//for _, field := range cols {
	//	if field == internal.UpdatedAt {
	//		hasUpdatedColumn = true
	//	}
	//}
	//if hasUpdatedColumnTag == true && hasUpdatedColumn == false {
	//	cols = append(cols,internal.UpdatedAt)
	//}
	sql,values := x.generateUpdateSQL(stmt, cols, rowMap)

	return sql,values
}

//UPDATE table_name
//SET
//type_id = CASE
//WHEN (id=1 AND UUID='v1') THEN 3
//WHEN (id=2 AND UUID='v2') THEN 4
//END,
//category_id = CASE
//WHEN (id=1 AND UUID='v1') THEN 6
//WHEN (id=2 AND UUID='v2') THEN 7
//END
//WHERE (id,UUID) IN ((1,'v1'),(2,'v2'))

func (x *GormDB) generateUpdateSQL(stmt *gorm.Statement, cols []string,
	rowMap []map[string]interface{}) (string,[]interface{}) {

	tableName := stmt.Table
	// 有序的主键列表
	var primaryKeys []string
	for _, field := range stmt.Statement.Statement.Schema.PrimaryFields {
		primaryKeys = append(primaryKeys, field.DBName)
	}

	// 获取数据库类型
	columnTypeMap := make(map[string]schema.DataType)
	for _, field := range stmt.Statement.Statement.Schema.Fields {
		columnTypeMap[field.DBName] = field.DataType
	}

	// 实体字段、数据库字段对应关系
	dbTagMap := make(map[string]string) // StructName-DbColumnName
	for _, field := range stmt.Statement.Statement.Schema.Fields {
		dbTagMap[field.DBName] = field.Name
	}

	// 组装SQL
	//UPDATE table_name
	//SET
	var sql string
	var valueList []interface{}
	sql = fmt.Sprintf("UPDATE `%s` SET ", tableName)

	//type_id = CASE
	//...
	//END,
	for _, colName := range cols {
		sql += fmt.Sprintf("`%s` = CASE ", colName)

		for _, row := range rowMap {
			// WHEN (id=1 AND UUID='xxx0') THEN 3
			ss, values := x.buildCaseLine(primaryKeys, colName, row, columnTypeMap, dbTagMap)
			sql += ss
			valueList = append(valueList,values...)
		}
		sql += "END,"
	}
	sql = strings.TrimSuffix(sql, ",")

	//WHERE (id,UUID) IN ((1,'v1'),(2,'v2'))
	ss,values :=x.buildWhereLine(primaryKeys, rowMap, columnTypeMap, dbTagMap)
	sql += ss
	valueList = append(valueList,values...)

	return sql,valueList
}

//WHERE (id,UUID) IN ((1,'v1'),(2,'v2'))
func (x *GormDB) buildWhereLine(primaryKeys []string, rowMap []map[string]interface{},
	columnTypeMap map[string]schema.DataType, dbTagMap map[string]string) (string,[]interface{}) {
	var values []interface{}
	sql := " WHERE ("
	for _, primaryKey := range primaryKeys {
		sql += fmt.Sprintf("`%s`,", primaryKey)
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += ") IN ("
	// ((1,'v1'),(2,'v2'))
	for _, item := range rowMap {
		sql += "("
		for _, primaryKey := range primaryKeys {
			name := dbTagMap[primaryKey]
			value := item[name]
			sql += "?,"
			values = append(values,value)
		}
		sql = strings.TrimSuffix(sql, ",")
		sql += "),"
	}
	sql = strings.TrimSuffix(sql, ",")
	sql += ")"
	return sql,values
}

//WHEN (id=1 AND UUID='v1') THEN 3
func (x *GormDB) buildCaseLine(primaryKeys []string, colName string, row map[string]interface{},
	columnTypeMap map[string]schema.DataType, dbTagMap map[string]string) (string, []interface{}) {
	var values []interface{}
	sql := "WHEN ("
	for _, primaryKey := range primaryKeys {
		name := dbTagMap[primaryKey]
		value := row[name]
		sql += fmt.Sprintf("`%s`=? AND ", primaryKey)
		values = append(values,value)
	}
	sql = strings.TrimSuffix(sql, " AND ")

	name := dbTagMap[colName]
	value := row[name]
	sql += ") THEN ? "
	values = append(values,value)

	return sql,values
}

// 将实体转换为map
func convertStructValue(batchSize int,rows interface{}) [][]map[string]interface{} {
	var itemList [][]map[string]interface{}
	var rowMap []map[string]interface{}
	raw := reflect.ValueOf(rows).Elem()
	total := raw.Len()

	for k := 0; k < total; k++ {
		valueMap := make(map[string]interface{})
		val := reflect.ValueOf(raw.Index(k).Interface())
		tp := reflect.Indirect(val).Type()
		for i := 0; i < val.NumField(); i++ {
			name := tp.Field(i).Name
			value := val.Field(i).Interface()
			valueMap[name] = getNullValue(value)
		}
		rowMap = append(rowMap, valueMap)
		if len(rowMap) == batchSize {
			itemList = append(itemList, rowMap)
			rowMap = nil
		}
	}

	if rowMap != nil && len(rowMap) > 0 {
		itemList = append(itemList, rowMap)
	}

	return itemList
}

func getNullValue(value interface{}) interface{}{
	vString,ok := value.(null.String)
	if ok {
		return vString.ValueOrZero()
	}
	vTime,ok := value.(null.Time)
	if ok {
		return vTime.ValueOrZero().Local().Format("2006-01-02 15:04:05.000000000")
	}
	vTime2,ok := value.(time.Time)
	if ok {
		return vTime2.Local().Local().Format("2006-01-02 15:04:05.000000000")
	}

	vFloat,ok := value.(null.Float)
	if ok {
		return vFloat.ValueOrZero()
	}
	vInt,ok := value.(null.Int)
	if ok {
		return vInt.ValueOrZero()
	}
	vBool,ok := value.(null.Bool)
	if ok {
		return vBool.ValueOrZero()
	}
	return value
}

// 根据数据库类型转换值
func convertDbValue_bk(columnTypeMap map[string]schema.DataType, colName string, value interface{}) interface{} {
	var sql interface{}

	if v, ok := columnTypeMap[colName]; ok {
		//if colName == internal.UpdatedAt {
		//	return fmt.Sprintf("'%s'", time.Now())
		//}

		switch v {
		case schema.String, schema.Time:
			if value == nil {
				sql = "''"
			} else {
				sql = fmt.Sprintf("%s", value)
			}
		case schema.Int, schema.Bool, schema.Uint:
			if value == nil {
				sql = 0
			} else {
				sql = fmt.Sprintf("%d", value)
			}
		case schema.Float:
			if value == nil {
				sql = 0
			} else {
				sql = fmt.Sprintf("%f", value)
			}
		default:
			internal.ThrowErrorMessage("%s不支持类型:%v"+colName, v)
		}
	} else {
		internal.ThrowErrorMessage("没有匹配上字段:" + colName)
	}

	return sql
}
