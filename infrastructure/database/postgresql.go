package database

import (
	"database/sql"
	"fmt"
	"time"
	"workshop-management/pkg/logger"
	"workshop-management/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnDb() (db *gorm.DB, sqlDB *sql.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		utils.GetEnv("DB_HOST", "").(string),
		utils.GetEnv("DB_PORT", "").(string),
		utils.GetEnv("DB_USERNAME", "").(string),
		utils.GetEnv("DB_PASS", "").(string),
		utils.GetEnv("DB_NAME", "").(string))
	logger.WriteLog(logger.LogLevelDebug, fmt.Sprintf("ConnDb; Initialize db connection..."))

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("ConnDb; %s Error: %s", dsn, err.Error()))
		return
	}

	maxIdle := 10
	maxIdleTime := 5 * time.Minute
	maxConn := 100
	maxLifeTime := time.Hour

	sqlDB, err = db.DB()
	if err != nil {
		logger.WriteLog(logger.LogLevelError, fmt.Sprintf("ConnDb.sqlDB; %s Error: %s", dsn, err.Error()))
		return
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(maxConn)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifeTime)

	db.Debug()

	return
}
