package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	myjson "github.com/roseboy/httpcase/json"
	"github.com/roseboy/httpcase/util"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Common struct {
}

func (f *Common) Print(o interface{}) string {
	var (
		ret = fmt.Sprintf("%v", o)
	)
	t := reflect.TypeOf(o)
	if t.Name() == "" { //obj
		str, err := json.Marshal(o)
		if err == nil {
			ret = string(str)
		}
	}

	if strings.Contains(ret, "\n") {
		return fmt.Sprintf("\n%s", ret)
	}
	return ret
}

func (f *Common) PrintJson(o interface{}) string {
	var (
		str  bytes.Buffer
		data = f.Print(o)
	)
	err := json.Indent(&str, []byte(data), "", "  ")
	if err != nil {
		return data
	}
	data = str.String()

	if util.IsWindows() {
		return data
	}

	jsonRaw, err := myjson.Highlight([]byte(data))
	if err != nil {
		return data
	}

	return fmt.Sprintf("\n%s", string(jsonRaw))
}

func (f *Common) Sleep(t string) bool {
	s, err := strconv.Atoi(t)
	if err != nil {
		return false
	}
	time.Sleep(time.Millisecond * time.Duration(s))
	return true
}
