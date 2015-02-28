package util

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"net"
	"reflect"
	"regexp"
	"strings"
)

func GenUUID() (string, error) {
	u4, err := uuid.NewV4()
	return u4.String(), err
}

func EmptyPort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	defer l.Close()

	if err != nil {
		return 0, fmt.Errorf("Fail to listen empty port")
	}

	addr := l.Addr()
	port := addr.(*net.TCPAddr).Port
	return port, nil
}

func SetAction(actions map[string]reflect.Value, in interface{}) {
	t := reflect.TypeOf(in)
	prefix := reflect.ValueOf(in).MethodByName("Prefix").Call([]reflect.Value{})[0]
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		match, _ := regexp.MatchString("^Action", method.Name)
		if !match {
			continue
		}
		action := reflect.ValueOf(in).MethodByName(method.Name)
		field_name := fmt.Sprintf("%s_%s", prefix, strings.TrimLeft(method.Name, "Action"))
		actions[field_name] = action
	}
}
