package httpcase

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"os"
	"path/filepath"
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
	ctx *TestContext
	cas *TestCase
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
	f.ctx.SetAttr(key, AttrValueTypeSet, value)
	return true
}

func (f *TestFunctions) Env(envName string) bool {
	f.ctx.Env = envName
	return true
}

func (f *TestFunctions) EnvSet(env, key, value string) bool {
	if f.ctx.Env != "" && f.ctx.Env == env {
		f.ctx.SetEnvAttr(key, env, value)
	}
	return true
}

func (f *TestFunctions) Callback(url string) bool {
	f.ctx.CallbackUrl = url
	return true
}

func (f *TestFunctions) CallbackWithFunction(jsFunctionName, url string) bool {
	f.ctx.CallbackJsFunction = jsFunctionName
	f.ctx.CallbackUrl = url
	return true
}

func (f *TestFunctions) LoadData(file string) (bool, error) {
	text, err := util.ReadText(f.ctx.testRequest.TestCaseFile.Dir + file)
	if err != nil {
		return false, err
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal([]byte(text), &jsonMap)
	if err != nil {
		return false, err
	}
	for k, v := range jsonMap {
		f.ctx.SetAttr(k, AttrValueTypeData, fmt.Sprintf("%v", v))
	}
	return true, nil
}

func (f *TestFunctions) Import(name, file string) (bool, error) {
	script, err := util.ReadText(f.ctx.testRequest.TestCaseFile.Dir + file)
	if err != nil {
		return false, err
	}

	jsVm := otto.New()
	_, err = jsVm.Run(script)
	if err != nil {
		return false, err
	}
	f.ctx.JSVms[name] = jsVm
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
	if jsVm, ok := f.ctx.JSVms[jsLib]; ok {
		enc, err := jsVm.Call(jsFunc, nil, args...)
		if err != nil {
			return "", err
		}
		return enc.ToString()
	}
	return "", fmt.Errorf("can't find js import:%s", jsLib)
}

func (f *TestFunctions) Header(key, value string) bool {
	f.ctx.CommonHeader[key] = value
	return true
}

func (f *TestFunctions) Body(body string) bool {
	f.cas.Request.Body = body
	return true
}

func (f *TestFunctions) Param(key, value string) bool {
	f.cas.Request.Params[key] = value
	return true
}

func (f *TestFunctions) File(name, path string) (bool, error) {
	dir, _ := filepath.Split(f.ctx.testRequest.TestCaseFile.Name)
	file, err := os.Open(dir + path)
	if err != nil {
		return false, nil
	}
	if f.cas.Request.UploadFiles == nil {
		f.cas.Request.UploadFiles = make(map[string]*os.File)
	}
	f.cas.Request.UploadFiles[name] = file
	return true, nil
}

func (f *TestFunctions) AllowRedirect() bool {
	f.cas.Request.AllowRedirect = true
	return true
}

func (f *TestFunctions) While(v1 string, op string, v2 string) bool {
	f.cas.While, _ = f.Assert(v1, op, v2)
	return f.cas.While
}

func (f *TestFunctions) If(v1 string, op string, v2 string) bool {
	f.cas.If, _ = f.Assert(v1, op, v2)
	return f.cas.If
}

func (f *TestFunctions) Foreach(v1 string, op string, v2 string) bool {

	return true
}
