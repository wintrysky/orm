package orm

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/wintrysky/orm/internal"
)

// GetItemWhere 获取数据
// entities 返回的数据列表
// query 查询条件
// args 查询条件参数
// GetItemWhere(&items,"workflow_uuid = ? and address like ?","xxx","%xxx%")
func (x *GormDB) GetItemWhere(entitiesPtr interface{}, query string, args ...interface{}) {
	defer x.sessionHandler()

	x.db = x.db.Session(&gorm.Session{NewDB: true})
	x.db = x.db.Where(query, args...).Find(entitiesPtr)
	if x.db.Error != nil {
		if x.db.Error == context.DeadlineExceeded {
			x.db.Error = errors.New("SQL超时")
		}
	}
}

// entity 返回的实体数据
// query 查询条件
// args 查询条件参数
// GetItemWhereFirst(&item,"workflow_uuid = ? and address like ?","xxx","%xxx%")
func (x *GormDB) GetItemWhereFirst(entityPtr interface{}, query string, args ...interface{}) {
	defer x.sessionHandler()

	x.db = x.db.Session(&gorm.Session{NewDB: true})
	x.db = x.db.Where(query, args...).First(entityPtr)
}

// GetItemByCondition 根据条件查询
func (x *GormDB) GetItemByCondition(entityPtr interface{}, conditions []Condition) {
	defer x.sessionHandler()

	query, args, err := BuildCondition(conditions)
	internal.ThrowError(err)
	x.db = x.db.Session(&gorm.Session{NewDB: true})
	x.db = x.db.Where(query, args...).Find(entityPtr)
	if x.db.Error != nil {
		if x.db.Error == context.DeadlineExceeded {
			x.db.Error = errors.New("SQL执行超时")
		}
	}
}

// QueryRawSQL 根据自定义SQL查询
// entity 待返回的数据实体列表
// rawSql 完整的SQL语句
// args 参数
// QueryRawSQL(&items,"select * from table_xxx where uuid = ?","xxx")
func (x *GormDB) ExecuteTextQuery(entityPtr interface{}, rawSQL string, args ...interface{}) {
	defer x.sessionHandler()

	x.db = x.db.Session(&gorm.Session{NewDB: true})
	clone := x.db.Raw(rawSQL, args...)
	internal.ThrowError(clone.Error)

	x.db = clone.Scan(entityPtr)
}

// Insert 新增一行数据
// entity 实体或实体数组
// cols 【可选】定义新增字段,如果为空，更新表中所有字段
// 更新单行: Insert(&item)
// 更新单行且只更新特定字段: Insert(&item,[]string{"uuid","node"})
func (x *GormDB) Insert(entityPtr interface{}, cols ...[]string) {
	defer x.sessionHandler()

	// 判断cols数组长度
	if len(cols) > 1 {
		internal.ThrowErrorMessage("cols数组长度不能大于1")
	}
	x.db = x.db.Session(&gorm.Session{NewDB: true})

	if len(cols) == 0 {
		x.db = x.db.Create(entityPtr)
	} else {
		msg := x.checkColumns(entityPtr, cols[0])
		if msg != "" {
			internal.ThrowErrorMessage("新增失败，不存在字段%s", msg)
		}
		x.db = x.db.Select(cols[0]).Create(entityPtr)
	}
}

// BatchInsert 批量新增或插入
func (x *GormDB) BatchInsert(entitiesPtr interface{}, batchSize int, cols ...[]string) {
	defer x.sessionHandler()

	if len(cols) > 1 {
		internal.ThrowErrorMessage("cols参数不能大于1")
	}
	if batchSize < 1 {
		internal.ThrowErrorMessage("batchSize不能小于1")
	}
	x.db = x.db.Session(&gorm.Session{
		DisableNestedTransaction: true,
		NewDB:                    true,
	})

	if len(cols) == 1 {
		msg := x.checkColumns(entitiesPtr, cols[0])
		if msg != "" {
			internal.ThrowErrorMessage("BatchInsert失败，不存在字段%s", msg)
		}

		x.db = x.db.Select(cols[0])
	}
	x.db = x.db.CreateInBatches(entitiesPtr, batchSize)
}

// Update 更新单行数据
// entity 待更新实体
// cols 更新字段 []string{"name","address"}
// 备注: entity主键不能为空,当为复合主键时,支持其中任意一个主键不为空
func (x *GormDB) Update(entityPtr interface{}, cols []string) {
	defer x.sessionHandler()

	if len(cols) == 0 {
		internal.ThrowErrorMessage("cols不能为空")
	}
	x.db = x.db.Session(&gorm.Session{NewDB: true})

	msg := x.checkColumns(entityPtr, cols)
	if msg != "" {
		internal.ThrowErrorMessage("Update失败，不存在字段%s", msg)
	}
	x.db = x.db.Select(cols).Model(entityPtr).Updates(entityPtr)
}

// BatchUpdate 批量更新
func (x *GormDB) BatchUpdate(entitiesPtr interface{}, batchSize int, cols []string) {
	defer x.sessionHandler()

	if len(cols) == 0 {
		internal.ThrowErrorMessage("cols不能为空")
	}

	// 分批处理
	// 将实体转换为map
	var rowList [][]map[string]interface{}
	rowList = convertStructValue(batchSize,entitiesPtr)

	// 减少事务等待时间
	var sqlList []string
	var valueList [][]interface{}
	for _, rows := range rowList {
		sql,values := x.batchUpdateInternal(entitiesPtr, cols, rows)
		sqlList = append(sqlList, sql)
		valueList = append(valueList,values)
	}
	x.db = x.db.Session(&gorm.Session{NewDB: true})

	if x.isTx == false {
		x.db = x.db.Begin()
		defer x.EndTransaction()
	}

	var cnt int64
	for idx, sql := range sqlList {
		//start := time.Now()
		value := valueList[idx]
		clone := x.db.Exec(sql,value...)
		if clone.Error != nil {
			internal.ThrowErrorMessage("Error %v", clone.Error)
		}
		cnt += clone.RowsAffected
		//el := time.Since(start)
	}
	x.RowsAffected = cnt
}

// BatchUpsert 批量更新或新增
// entities 待新增的实体列表
// batchSize 每一条批量更新SQL更新的行数，超过这个数目，将自动拆分，并在一个事务中执行
// cols 需要更新的字段
// uniqueCols 唯一索引列表或唯一索引组合键
// BatchInsert(&items,500,[]string{"uuid","node"},[]string{"uuid"})
// INSERT INTO `table_name` (`node`,`uuid`) VALUES ('1','xxx1'),('2','xxx2')
func (x *GormDB) BatchUpsert(entitiesPtr interface{}, batchSize int, cols []string, uniqueCols []string)  {
	defer x.sessionHandler()

	if batchSize < 1 {
		internal.ThrowErrorMessage("batchSize不能小于1")
	}
	x.db = x.db.Session(&gorm.Session{
		DisableNestedTransaction: true,
		CreateBatchSize:          batchSize,
		PrepareStmt:              false,
		NewDB:                    true,
	})

	if len(uniqueCols) == 0 {
		internal.ThrowErrorMessage("upsert 必须包含唯一索引")
	}

	var ucols []clause.Column
	for _, item := range uniqueCols {
		ucols = append(ucols, clause.Column{Name: item})
	}

	x.db = x.db.Clauses(clause.OnConflict{
		Columns:   ucols,
		DoUpdates: clause.AssignmentColumns(cols),
	})
	x.db = x.db.Create(entitiesPtr)
}

// Execute 自定义执行DDL/DML语句
// Execute("update user set uuid = ?,type_id=? where id = ?","xxx",3,121)
func (x *GormDB) Execute(rawSQL string, args ...interface{}) {
	defer x.sessionHandler()

	x.db = x.db.Session(&gorm.Session{NewDB: true})
	x.db = x.db.Exec(rawSQL, args...)
}

// PagingQueryByCondition 分页查询
func (x *GormDB) PagingQueryByCondition(entityPtr interface{}, pageSize int,
	currentPage int, orderBy string, conditions []Condition) (int64, error) {
	// build where condition
	whereStr, param, err := BuildCondition(conditions)
	internal.ThrowError(err)

	return x.PagingQueryWhere(entityPtr, pageSize, currentPage, orderBy, whereStr, param...)
}

// PagingQueryWhere 分页查询
func (x *GormDB) PagingQueryWhere(entityPtr interface{}, pageSize int,
	currentPage int, orderBy string, where string, args ...interface{}) (count int64, err error) {
	// 获取总行数
	x.db = x.db.Session(&gorm.Session{NewDB: true})
	clone := x.db.Model(entityPtr).Where(where, args...).Count(&count)
	internal.ThrowError(clone.Error)

	if orderBy != "" {
		x.db = x.db.Order(orderBy)
	}

	if currentPage == 0 { // 初始页为1
		internal.ThrowErrorMessage("分页查询首页为1")
	}

	var offset = (currentPage - 1) * pageSize
	x.db = x.db.Where(where, args...).Limit(pageSize).Offset(offset).Find(entityPtr)

	return count, x.db.Error
}

//  AutoMigrate 会创建表、缺失的外键、约束、列和索引。 如果大小、精度、是否为空可以更改，
// 则 AutoMigrate 会改变列的类型。 出于保护您数据的目的，它 不会 删除未使用的列
func (x *GormDB) AutoMigrate(entityPtr interface{}) (err error) {
	defer x.sessionHandler()

	return x.db.AutoMigrate(entityPtr)
}

func (x *GormDB) sessionHandler() {
	if x.RowsAffected == 0 {
		x.RowsAffected = x.db.RowsAffected
	}

	if x.db.Error != nil {
		x.Error = x.db.Error
	}

	if r := recover(); r != nil {
		msg := fmt.Sprint(r)
		x.Error = errors.New(msg)
	}
}

func (x *GormDB) checkColumns(entityPtr interface{}, cols []string) string {
	stmt := x.db.Session(&gorm.Session{DryRun: true}).First(entityPtr).Statement
	for _, col := range cols {
		if _, ok := stmt.Schema.FieldsByDBName[col]; !ok {
			return col
		}
	}

	return ""
}