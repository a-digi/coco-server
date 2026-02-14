package di_test

import (
	"testing"
	"github.com/a-digi/coco-server/server/di"
)

type dummyContext struct{}

func (d dummyContext) GetDatabaseManager() interface{} { return nil }
func (d dummyContext) GetLogger() interface{} { return nil }

func TestContextInterface(t *testing.T) {
	var _ di.Context = dummyContext{}
}
