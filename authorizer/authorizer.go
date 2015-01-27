package authorizer

import (
	"github.com/tomyhero/ore_server/context"
)

type IAuthorizer interface {
	Login(c *context.Context) bool
	Auth(c *context.Context) bool
}
