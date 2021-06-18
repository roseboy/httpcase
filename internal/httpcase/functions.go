package httpcase

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"os"
	"reflect"
	"strings"

	"github.com/roseboy/httpcase/internal/httpcase/functions"
	"github.com/roseboy/httpcase/util"
)

var (
	FunctionNameMap map[string]string
)

type TestFunctions struct {
	functions.Common
	functions.Strings
	functions.Numbers
	functions.Files
	functions.Asserts
	functions.Mock
}

func init() {
	FunctionNameMap = make(map[string]string)
	f := new(TestFunctions)
	typ := reflect.TypeOf(f)
	size := typ.NumMethod()
	for i := 0; i < size; i++ {
		FunctionNameMap[strings.ToLower(typ.Method(i).Name)] = typ.Method(i).Name
	}
}

func GetFunctionName(line string) (bool, string) {
	str := strings.Split(line, " ")
	name, ok := FunctionNameMap[strings.ToLower(str[0])]
	return ok, name
}

func (f *TestFunctions) Set(key string, value string) bool {
	TestHolder.TestContext.SetAttr(key, AttrValueTypeSet, value)
	return true
}

func (f *TestFunctions) Env(envName string) bool {
	if TestHolder.TestContext.Env == "" {
		TestHolder.TestContext.Env = envName
	}
	return true
}

func (f *TestFunctions) EnvSet(env, key, value string) bool {
	if TestHolder.TestContext.Env != "" && TestHolder.TestContext.Env == env {
		TestHolder.TestContext.SetEnvAttr(key, env, value)
	}
	return true
}

func (f *TestFunctions) Callback(url string) bool {
	TestHolder.TestContext.CallbackUrl = url
	return true
}

func (f *TestFunctions) CallbackWithFunction(jsFunctionName, url string) bool {
	TestHolder.TestContext.CallbackJsFunction = jsFunctionName
	TestHolder.TestContext.CallbackUrl = url
	return true
}

func (f *TestFunctions) LoadData(file string) (bool, error) {
	filePath := TestHolder.Function.Code.File.Dir + file
	if !util.IsExist(filePath) {
		filePath = file
	}
	text, err := util.ReadText(filePath)
	if err != nil {
		return false, err
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(text), &jsonMap)
	if err != nil {
		return false, err
	}
	for k, v := range jsonMap {
		TestHolder.TestContext.SetAttr(k, AttrValueTypeData, fmt.Sprintf("%v", v))
	}
	return true, nil
}

func (f *TestFunctions) Import(name, file string) (bool, error) {
	filePath := TestHolder.Function.Code.File.Dir + file
	if !util.IsExist(filePath) {
		filePath = file
	}
	script, err := util.ReadText(filePath)
	if err != nil {
		return false, err
	}

	jsVm := otto.New()
	_, err = jsVm.Run(script)
	if err != nil {
		return false, err
	}
	TestHolder.TestContext.JSVms[name] = jsVm
	return true, nil
}

func (f *TestFunctions) RunJs(argsStr string) (string, error) {
	args := strings.Split(argsStr, ",")
	argsInterface := make([]interface{}, 0)
	for _, v := range args {
		argsInterface = append(argsInterface, v)
	}
	return f.RunJsWithName(args[0], argsInterface[1:]...)
}

func (f *TestFunctions) RunJsWithName(functionName string, args ...interface{}) (string, error) {
	jsLib := functionName[:strings.Index(functionName, ".")]
	jsFunc := functionName[strings.Index(functionName, ".")+1:]
	if jsVm, ok := TestHolder.TestContext.JSVms[jsLib]; ok {
		enc, err := jsVm.Call(jsFunc, nil, args...)
		if err != nil {
			return "", err
		}
		return enc.ToString()
	}
	return "", fmt.Errorf("can't find js import:%s", jsLib)
}

func (f *TestFunctions) Header(key, value string) bool {
	TestHolder.TestContext.CommonHeader[key] = value
	return true
}

func (f *TestFunctions) Body(body string) bool {
	TestHolder.TestCase.Request.Body = body
	return true
}

func (f *TestFunctions) Param(key, value string) bool {
	TestHolder.TestCase.Request.Params[key] = value
	return true
}

func (f *TestFunctions) File(name, path string) (bool, error) {
	filePath := TestHolder.Function.Code.File.Dir + path
	if !util.IsExist(filePath) {
		filePath = path
	}
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	if TestHolder.TestCase.Request.UploadFiles == nil {
		TestHolder.TestCase.Request.UploadFiles = make(map[string]*os.File)
	}
	TestHolder.TestCase.Request.UploadFiles[name] = file
	return true, nil
}

func (f *TestFunctions) AllowRedirect() bool {
	TestHolder.TestCase.Request.AllowRedirect = true
	return true
}

func (f *TestFunctions) While(v1 string, op string, v2 string) bool {
	TestHolder.TestCase.While, _ = f.Assert(v1, op, v2)
	return TestHolder.TestCase.While
}

func (f *TestFunctions) If(v1 string, op string, v2 string) bool {
	TestHolder.TestCase.If, _ = f.Assert(v1, op, v2)
	return TestHolder.TestCase.If
}

func (f *TestFunctions) Foreach(v1 string, op string, v2 string) bool {

	return true
}
