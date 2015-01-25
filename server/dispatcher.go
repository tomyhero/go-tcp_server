package server

import (
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/util"
	"reflect"
)

type Dispatcher struct {
	Handlers []context.IHandler
	Actions  map[string]reflect.Value
}

func NewDispatcher(handlers []context.IHandler) *Dispatcher {

	actions := map[string]reflect.Value{}
	for _, handler := range handlers {
		util.GetMethods(actions, handler)
	}

	return &Dispatcher{Handlers: handlers, Actions: actions}
}
