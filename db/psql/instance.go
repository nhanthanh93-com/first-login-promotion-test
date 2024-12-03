package psql

import (
	"gorm.io/gorm"
)

type Instance[T any] struct {
	DBName string
	DBAddr string
	DB     *gorm.DB
}

func NewInstance[T any](dbManager *DBManager, dbName string) *Instance[T] {
	db, exists := dbManager.GetDB(dbName)
	if !exists {
		panic("Database connection not found: " + dbName)
	}

	if err := db.AutoMigrate(new(T)); err != nil {
		panic("Failed to auto-migrate: " + err.Error())
	}

	return &Instance[T]{DBName: dbName, DB: db}
}
