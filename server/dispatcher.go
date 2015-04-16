package server

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/util"
	"net"
	"reflect"
	"strings"
)

type Dispatcher struct {
	Handlers         []context.IHandler
	mapHandlers      map[string]context.IHandler
	Actions          map[string]reflect.Value
	LoginActions     map[string]reflect.Value
	ReconnectActions map[string]reflect.Value
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
	reconnectActions := map[string]reflect.Value{}
	mapHandlers := map[string]context.IHandler{}

	for _, handler := range handlers {
		mapHandlers[handler.Prefix()] = handler
		util.SetAction(actions, handler)

		login_field_name := fmt.Sprintf("%s_login", handler.Prefix())
		loginActions[login_field_name] = reflect.ValueOf(handler.GetAuthorizer()).MethodByName("Login")

		reconnect_field_name := fmt.Sprintf("%s_reconnect", handler.Prefix())
		reconnectActions[reconnect_field_name] = reflect.ValueOf(handler.GetAuthorizer()).MethodByName("Reconnect")
	}

	return &Dispatcher{Handlers: handlers, mapHandlers: mapHandlers, Actions: actions, LoginActions: loginActions, ReconnectActions: reconnectActions}
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

func (d *Dispatcher) HookInitialize(database map[string]interface{}) {
	for _, handler := range d.Handlers {
		database[handler.Prefix()] = map[string]interface{}{}
		handler.HookInitialize(database)
	}
}
func (d *Dispatcher) HookDestroy(database map[string]interface{}) {
	for _, handler := range d.Handlers {
		handler.HookDestroy(database)
	}
}
func (d *Dispatcher) HookDisconnect(conn net.Conn) {
	for _, handler := range d.Handlers {
		handler.HookDisconnect(conn)
	}
}
