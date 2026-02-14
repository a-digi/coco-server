package di

import (
	db "github.com/a-digi/coco-orm/orm"
	"github.com/a-digi/coco-logger/logger"
)

// Context defines the interface for dependency injection.
type Context interface {
	GetDatabaseManager() *db.DatabaseManager
	GetLogger() logger.Logger
}
