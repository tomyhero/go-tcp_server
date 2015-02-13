package server

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/util"
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
	return handler.GetAuthorizer().Auth(c)
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
		loginActions[login_field_name] = reflect.ValueOf(handler.GetAuthorizer()).MethodByName("Login")
	}
	return &Dispatcher{Handlers: handlers, mapHandlers: mapHandlers, Actions: actions, LoginActions: loginActions}
}

func (d *Dispatcher) BeforeExecute(c *context.Context, cmd string) {
	prefix := strings.Split(cmd, "_")[0]
	handler := d.mapHandlers[prefix]
	handler.HookBeforeExecute(c)
}

func (d *Dispatcher) AfterExecute(c *context.Context, cmd string) {
	prefix := strings.Split(cmd, "_")[0]
	handler := d.mapHandlers[prefix]
	handler.HookAfterExecute(c)
}

func (d *Dispatcher) HookInitialize(gstore map[string]interface{}) {
	for _, handler := range d.Handlers {
		myStore := map[string]interface{}{}
		gstore[handler.Prefix()] = myStore
		handler.HookInitialize(gstore, myStore)
	}
}
func (d *Dispatcher) HookDestroy(gstore map[string]interface{}) {
	for _, handler := range d.Handlers {

		myStore := gstore[handler.Prefix()].(map[string]interface{})
		handler.HookDestroy(gstore, myStore)
	}
}
