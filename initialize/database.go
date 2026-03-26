package initialize

import (
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/irorikon/cloudgram-go/config"
	"github.com/irorikon/cloudgram-go/logger"
	"github.com/irorikon/cloudgram-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CountTable 用于统计指定表中的记录数量
// tableName: 要统计的表名
// 返回值: int64 类型，表示表中记录的总数
func CountTable(tableName string) (num int64) {
	// 使用数据库连接 config.DB 对指定表进行计数查询
	// 将查询结果赋值给 num 变量
	config.DB.Table(tableName).Count(&num)
	// 返回查询得到的记录数量
	return
}

// InitDB 函数用于初始化数据库连接
// 根据配置文件中指定的数据库类型，创建对应类型的数据库连接
// 支持多种数据库类型：SQLite、MSSQL、MySQL、PostgreSQL
// 如果未指定数据库类型，默认使用SQLite
// @return error 初始化数据库连接时可能返回的错误信息
func InitDB() (err error) {
	// 根据配置的数据库类型选择对应的数据库驱动
	switch strings.ToLower(config.DBTYPE) {
	case "sqlite":
		config.DB, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	case "mysql":
		// use MySQL
		config.DB, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	case "postgres":
		// use PostgreSQL
		config.DB, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	default:
		// use SQLite as default
		config.DB, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
			PrepareStmt: true, // precompile SQL
		})
	}
	return err
}

func CloseDB() error {
	sqlDB, err := config.DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}

func RegisterTables() {
	db := config.DB
	err := db.AutoMigrate(
		model.Channel{},
		model.File{},
		model.FileChunk{},
		model.TempChunk{},
	)
	if err != nil {
		logger.SysErrorF("register table failed: %v", err)
		os.Exit(0)
	}
	logger.SysLog("register table success")
}
