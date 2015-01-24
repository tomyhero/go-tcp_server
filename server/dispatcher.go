package server

type Dispatcher struct {
	Handlers *interface{}
}

func NewDispatcher(handlers *interface{}) *Dispatcher {
	return &Dispatcher{Handlers: handlers}
}
