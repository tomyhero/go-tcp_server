package util

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func GetMethods(actions map[string]reflect.Value, in interface{}) {
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
