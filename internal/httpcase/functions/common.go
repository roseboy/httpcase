package functions

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Common struct {
}

func (f *Common) Print(o interface{}) string {
	t := reflect.TypeOf(o)
	if t.Name() == "" {
		str, err := json.Marshal(o)
		if err != nil {
			return fmt.Sprintf("%v", o)
		}
		return string(str)
	}
	return fmt.Sprintf("%v", o)
}

func (f *Common) Sleep(t string) bool {
	s, err := strconv.Atoi(t)
	if err != nil {
		return false
	}
	time.Sleep(time.Millisecond * time.Duration(s))
	return true
}
