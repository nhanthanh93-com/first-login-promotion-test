package psql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

type DBConn struct {
	DB *gorm.DB
}

var (
	dbInstance *DBConn
	once       sync.Once
	loc        *time.Location
)

type DBManager struct {
	databases map[string]*gorm.DB
}

func NewDBManager() *DBManager {
	return &DBManager{
		databases: make(map[string]*gorm.DB),
	}
}

func (manager *DBManager) Connect(name, dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	manager.databases[name] = db
	return nil
}

func (manager *DBManager) GetDB(name string) (*gorm.DB, bool) {
	db, exists := manager.databases[name]
	return db, exists
}

func (manager *DBManager) Disconnect(name string) error {
	db, exists := manager.databases[name]
	if !exists {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
