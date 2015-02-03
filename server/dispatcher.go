package server

import (
	"fmt"
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/util"
	"reflect"
	"strings"
)

type Dispatcher struct {
	Handlers     []context.IHandler
	mapHandlers  map[string]context.IHandler
	Actions      map[string]reflect.Value
	LoginActions map[string]reflect.Value
}

func (d *Dispatcher) ExecAuth(c *context.Context, cmd string) bool {
	prefix := strings.Split(cmd, "_")[0]
	handler := d.mapHandlers[prefix]
	return handler.AuthorizerHandler().Auth(c)
}

func (d *Dispatcher) GetHandler(prefix string) context.IHandler {
	return d.mapHandlers[prefix]
}

func NewDispatcher(handlers []context.IHandler) *Dispatcher {
	actions := map[string]reflect.Value{}
	loginActions := map[string]reflect.Value{}
	mapHandlers := map[string]context.IHandler{}

	for _, handler := range handlers {
		mapHandlers[handler.Prefix()] = handler
		util.GetMethods(actions, handler)
		login_field_name := fmt.Sprintf("%s_login", handler.Prefix())
		loginActions[login_field_name] = reflect.ValueOf(handler.AuthorizerHandler()).MethodByName("Login")
	}
	return &Dispatcher{Handlers: handlers, mapHandlers: mapHandlers, Actions: actions, LoginActions: loginActions}
}
