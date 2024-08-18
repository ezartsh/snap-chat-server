package database

import (
	"snap_chat_server/config"
	"snap_chat_server/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDbConnection() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: config.Env.Db.ConnectionString(),
	}), &gorm.Config{
		Logger: logger.NewGormPackageLogger(),
	})

	if err != nil {
		logger.AppLog.Fatalf(err, "Database connection failed.")
	}

	return db
}
