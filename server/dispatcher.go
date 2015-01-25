package server

import (
	"github.com/tomyhero/ore_server/util"
	"reflect"
)

type Dispatcher struct {
	Handlers []interface{}
	Actions  map[string]reflect.Value
}

func NewDispatcher(handlers []interface{}) *Dispatcher {

	actions := map[string]reflect.Value{}
	for _, handler := range handlers {
		util.GetMethods(actions, handler)
	}

	return &Dispatcher{Handlers: handlers, Actions: actions}
}
