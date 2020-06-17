package db_mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"time"
)

var db *gorm.DB

type NewMysqlConfig struct {
	Host            string
	Debug           bool
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}

func NewMysql(cfg NewMysqlConfig) error {
	var err error
	db, err = gorm.Open("mysql", cfg.Host)
	if err != nil {
		return err
	}
	db.SingularTable(true)
	db.DB().SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	db.LogMode(cfg.Debug)
	return nil
}

func Close() {
	defer db.Close()
}

//一致事务
func Translation(fun func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%+v", r))
			tx.Rollback()
		}
	}()
	if err := fun(tx); err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return
}

func GetDb() *gorm.DB {
	return db
}
