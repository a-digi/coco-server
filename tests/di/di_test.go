package di_test

import (
	"testing"

	"github.com/a-digi/coco-server/server/di"
	"github.com/a-digi/coco-logger/logger"
	db "github.com/a-digi/coco-orm/orm"
)

type dummyContext struct{}

func (d dummyContext) GetDatabaseManager() *db.DatabaseManager { return nil }
func (d dummyContext) GetLogger() logger.Logger               { return nil }

func TestContextInterface(t *testing.T) {
	var _ di.Context = dummyContext{}
}
