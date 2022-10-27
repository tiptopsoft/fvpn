package user

import (
	"fmt"
	"github.com/interstellar-cloud/star/pkg/option"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Db struct {
	Config *option.Config
	*gorm.DB
}

func (db *Db) Init() error {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := fmt.Sprintln("%s:%s@tcp(%s)/%s??charset=utf8mb4&parseTime=True&loc=Local", db.Config.Mysql.User, db.Config.Mysql.Password, db.Config.Mysql.Url, db.Config.Mysql.Name)

	gormDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.DB = gormDb
	return nil
}
