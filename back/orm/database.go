package orm

import (
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	GormDB *gorm.DB
)

func InitDatabase() {
	var err error
	var gormLogger logger.Interface
	if os.Getenv("MODE") == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}
	GormDB, err = gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{
		PrepareStmt: true,
		Logger:      gormLogger,
		NowFunc:     func() time.Time { return time.Now().UTC() }, // fuck timezones
	})
	if err != nil {
		panic("failed to connect database " + err.Error())
	}

	err = GormDB.AutoMigrate(&Challenge{})

	if err != nil {
		panic("failed to migrate database " + err.Error())
	}
}
