package orm

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

type DBSettings struct {
	DBName   string
	Host     string
	User     string
	Password string
	Port     int
}
type GormDB struct {
	db             *gorm.DB // 局部
	isTx bool
	Error          error
	RowsAffected   int64
}


func InitDB(conf DBSettings) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:              time.Second * 10,   // Slow SQL threshold
			LogLevel:                   logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,          // Disable color
		},
	)

	var source string
	source = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)
	db, err = gorm.Open(mysql.Open(source), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: false,
		Logger: newLogger,
	})

	return err
}

func BeginTransaction() GormDB {
	var da GormDB

	da.db = db.Begin().Session(&gorm.Session{NewDB: true})
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(180))
	da.db = da.db.WithContext(ctx)
	da.isTx = true
	return da
}

func (x *GormDB)EndTransaction() {
	if r := recover(); r != nil {
		x.db.Rollback()
		panic(r)
	}
	if x.Error != nil {
		x.db.Logger.Error(context.Background(), "rollback",x.Error.Error())
		x.db.Rollback()
	} else {
		x.db.Logger.Info(context.Background(), "committed")
		x.db.Commit()
	}
}

// NewDB 新连接会话
func NewSession() GormDB {
	var da GormDB
	da.db = db.Session(&gorm.Session{NewDB: true})

	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(180))
	da.db = da.db.WithContext(ctx)

	return da
}