package models

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/Quons/go-gin-example/pkg/logging"
	"github.com/Quons/go-gin-example/pkg/setting"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

var wdb *gorm.DB
var readDbSlice []*gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	var err error
	wdb, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=5s",
		setting.DatabaseSetting.WUser,
		setting.DatabaseSetting.WPassword,
		setting.DatabaseSetting.WHost,
		setting.DatabaseSetting.Name))

	if err != nil {
		logrus.Error(err)
	}

	wdb.LogMode(true)
	wdb.SetLogger(log.New(logging.GetGinLogWriter(), "[GORM] ", log.Ldate))
	// 全局禁用表名复数,默认表名为结构体复数
	wdb.SingularTable(true)
	//wdb.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//wdb.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	//wdb.Callback().Delete().Replace("gorm:delete", deleteCallback)
	wdb.DB().SetMaxIdleConns(10)
	wdb.DB().SetMaxOpenConns(100)

	rHosts := strings.Split(setting.DatabaseSetting.RHost, "|")
	for _, rHost := range rHosts {
		rdb, err := gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=5s",
			setting.DatabaseSetting.RUser,
			setting.DatabaseSetting.RPassword,
			rHost,
			setting.DatabaseSetting.Name))

		if err != nil {
			logrus.Error(err)
		}
		rdb.LogMode(true)
		rdb.SetLogger(log.New(logging.GetGinLogWriter(), "[GORM] ", log.Ldate))
		// 全局禁用表名复数，默认表名为结构体复数
		rdb.SingularTable(true)
		rdb.DB().SetMaxIdleConns(10)
		rdb.DB().SetMaxOpenConns(100)
		readDbSlice = append(readDbSlice, rdb)
	}
}

//获取读库
func readDB() *gorm.DB {
	return readDbSlice[rand.Intn(len(readDbSlice))]
}

//获取写库
func WriteDB() *gorm.DB {
	return wdb
}

func CloseDB() {
	defer wdb.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func formatTimeStampForQueryCallback(scope *gorm.Scope) {

	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")

		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
