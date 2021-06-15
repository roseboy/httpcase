package functions

import (
	"errors"
	"fmt"
	"github.com/roseboy/httpcase/json"
	"github.com/roseboy/httpcase/util"
	"regexp"
	"strconv"
	"strings"
)

type Strings struct {
}

func (f *Strings) Len(t string) int {
	return len(t)
}

func (f *Strings) Replace(str string, old string, new string) string {
	return strings.Replace(str, old, new, -1)
}

func (f *Strings) ToLower(str string) string {
	return strings.ToLower(str)
}

func (f *Strings) ToUpper(str string) string {
	return strings.ToUpper(str)
}

func (f *Strings) TrimLeft(str string) string {
	return util.TrimLeft(str)
}

func (f *Strings) TrimRight(str string) string {
	return util.TrimLeft(str)
}

func (f *Strings) Trim(str string) string {
	return util.Trim(str)
}

func (f *Strings) Match(str string, regStr string) bool {
	regx := regexp.MustCompile(regStr)
	return regx.MatchString(str)
}

func (f *Strings) IndexOf(str string, substr string) int {
	return strings.Index(str, substr)
}

func (f *Strings) SubStr(str string, indexStr string, indexStr2 string) (string, error) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", err
	}
	index2, err := strconv.Atoi(indexStr2)
	if err != nil {
		return "", err
	}
	if index > len(str) || index2 > len(str) || index > index2 {
		return "", errors.New("index1 most grant index2")
	}
	return str[index:index2], nil
}

func (f *Strings) SubStr2(str string, beginStr string, endStr string) (string, error) {
	index := strings.Index(str, beginStr) + len(beginStr)
	index2 := strings.Index(str[index:], endStr) + index
	if index > len(str) || index2 > len(str) || index > index2 {
		return "", errors.New("index1 most grant index2")
	}
	return str[index:index2], nil
}

func (f *Strings) Concat(str1 string, str2 string) string {
	return fmt.Sprintf("%s%s", str1, str2)
}

func (f *Strings) Filter(listJson string, field, where, value string) (ret string, err error) {
	jsonObj, err := json.NewJsonObject(fmt.Sprintf(`{"list":%s}`, listJson))
	if err != nil {
		return "", err
	}

	jsonObj.GetArray("list").ForEach(func(i int, object *json.Object) {
		if object.Get(field).(string) == where {
			ret = fmt.Sprintf("%v", object.Get(value))
			return
		}
	})

	return ret, nil
}
