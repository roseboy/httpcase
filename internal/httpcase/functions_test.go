package httpcase

import (
	"fmt"
	"github.com/roseboy/httpcase/internal/httpcase/functions"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestAssertsBefore(t *testing.T) {
	a := &functions.Asserts{}
	fmt.Println(a.AssertBefore("2020-03-01", "2020-03-02 12:12:33"))
	fmt.Println(a.AssertBefore("2020-03-01 00:00:00", "2020-03-02 00:00:00"))
}

func TestAssertsTrim(t *testing.T) {
	a := &functions.Strings{}
	fmt.Println(a.TrimRight("\r\n abc 哈defg\r\n "))
	fmt.Println(a.TrimLeft("\r\n abc defg\r\n "))
	fmt.Println(a.Trim("\r\n abc  写大meinvdefg\r\n "))
}

func TestSubStr2(t *testing.T) {
	str := `<input name="xhq" value="xhqqq" />`

	a := &functions.Strings{}
	fmt.Println(a.SubStr2(str, "value=\"", "\""))

}

func TestFuncNames(t *testing.T) {
	f := new(TestFunctions)
	typ := reflect.TypeOf(f)
	size := typ.NumMethod()
	for i := 0; i < size; i++ {
		fmt.Println(typ.Kind(), typ.Method(i).Name)
	}
}

func TestMatch(t *testing.T) {
	ps := []string{"abcd*cd", "abc*", "*cdba", "abcd*cd*", "*abcd*cd", "*abcd*cd*"}
	for _, p := range ps {
		p1 := strings.Replace(p, "*", "(.*?)", -1)
		regx := regexp.MustCompile(p1)
		fmt.Println(p, p1, regx.MatchString("abcdcd"))
	}
}
