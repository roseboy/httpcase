package httpcase

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/roseboy/httpcase/json"
	"github.com/roseboy/httpcase/requests"
	"github.com/roseboy/httpcase/util"
)

type TestContext struct {
	Session            *requests.RequestBuilder
	VarIndex           int
	Attr               map[string]*AttrValue
	JSVms              map[string]*otto.Otto
	CallbackUrl        string
	CallbackJsFunction string
	CommonHeader       map[string]string

	Env string
	Out string

	testRequest    *TestRequest
	responseValues []map[string]interface{}
}

const (
	AttrValueTypeSet  = 1
	AttrValueTypeEnv  = 2
	AttrValueTypeTemp = 3
	AttrValueTypeData = 4
)

type AttrValue struct {
	Name  string
	Env   string
	Value interface{}
	Type  int
}

func NewTestContext() *TestContext {
	t := &TestContext{
		Session:        requests.NewHttpSession(),
		JSVms:          make(map[string]*otto.Otto),
		Attr:           make(map[string]*AttrValue),
		CommonHeader:   make(map[string]string),
		responseValues: make([]map[string]interface{}, 0),
		VarIndex:       1500000,
	}
	TestHolder.TestContext = t
	return t
}

func (t *TestContext) Init(testRequest *TestRequest) *TestContext {
	t.testRequest = testRequest
	return t
}

func (t *TestContext) SetAttr(name string, vType int, value interface{}) {
	t.Attr[name] = &AttrValue{
		Name:  name,
		Env:   "",
		Value: value,
		Type:  vType,
	}
}

func (t *TestContext) SetEnvAttr(name string, env string, value interface{}) {
	t.Attr[name] = &AttrValue{
		Name:  name,
		Env:   env,
		Value: value,
		Type:  AttrValueTypeEnv,
	}
}

func (t *TestContext) GetAttr(name string) interface{} {
	return t.Attr[name]
}

func (t *TestContext) GetAttrStr(name string) string {
	if v, ok := t.Attr[name]; ok {
		return v.Value.(string)
	}
	return ""
}

func (t *TestContext) Run() (*TestResult, error) {
	var (
		cases = t.testRequest.TestCases
	)

	//执行case
	for _, cas := range cases {
		TestHolder.TestCase = cas

		util.Println()
		util.Log(util.Green(fmt.Sprintf("[%s]", cas.Name)))

		//运行全局函数
		if cas.IsGlobal {
			t.RunFunction(cas, true, false)
			if !cas.Pass {
				return nil, fmt.Errorf("global error")
			}
			TestHolder.TestCase = nil
			continue
		}

		cas = t.RunFunction(cas, true, false)
		if !cas.If {
			TestHolder.TestCase = nil
			continue
		}
		loopCount := 0
		for cas.While || loopCount == 0 {
			loopCount++

			cas = t.RunCase(cas)
			if cas.Err != nil {
				cas.Pass = false
				cas = t.RunFunction(cas, false, true)
				util.Log(util.Red("ERROR:"), cas.Err.Error())
				util.Log(util.Blue("RESULT:"), "Pass =>", cas.Pass)
				continue
			}
			cas = t.RunFunction(cas, false, false)

			util.Log(util.Blue("RESULT:"), "Pass =>", cas.Pass)
		}
		TestHolder.TestCase = nil
	}

	var (
		total    = 0
		pass     = 0
		fail     = 0
		skip     = 0
		duration = int64(0)
	)

	for _, cas := range cases {
		if cas.IsGlobal {
			continue
		}
		total++
		if cas.Pass {
			pass++
		} else {
			fail++
		}
		duration = duration + cas.Time
	}

	return &TestResult{
		TestCases:    cases,
		Total:        total,
		Passed:       pass,
		Failed:       fail,
		Skipped:      skip,
		Duration:     duration,
		TestCaseFile: t.testRequest.TestCaseFile,
	}, nil
}

func (t *TestContext) RunCase(cas *TestCase) *TestCase {
	cas.Request.Url = t.RenderValueString(cas.Request.Url)
	cas.Request.Body = t.RenderValueString(cas.Request.Body)
	params := make(map[string]string)
	for k, v := range cas.Request.Params {
		params[t.RenderValueString(k)] = t.RenderValueString(v)
	}
	cas.Request.Params = params

	headers := make(map[string]string)
	for k, v := range t.CommonHeader {
		headers[t.RenderValueString(k)] = t.RenderValueString(v)
	}
	for k, v := range cas.Request.Headers {
		headers[t.RenderValueString(k)] = t.RenderValueString(v)
	}
	cas.Request.Headers = headers

	rp := t.Session.SendWithRequest(cas.Request)
	util.Log(util.Cyan(cas.Request.Method+":"), cas.Request.Url, "=>", util.IfStr(rp.Err != nil, "ERROR", fmt.Sprintf("%d", rp.Response.Status)))

	if rp.Err != nil {
		cas.Err = rp.Err
	}

	var err error
	cas.Response = rp.Response
	cas.Time = cas.Time + rp.Response.Time
	cas.Res, err = t.parseRes(rp.Response)
	if cas.Err == nil {
		cas.Err = err
	}

	t.responseValues = append(t.responseValues, cas.Res)
	return cas
}

func (t *TestContext) RunFunction(cas *TestCase, isBefore, skipAll bool) *TestCase {
	for _, f := range cas.Functions {
		if skipAll {
			f.Skip = true
			continue
		}
		if f.IsBefore == isBefore {
			RunFuncByName(t, cas, f)
		}
	}
	return cas
}

func (t *TestContext) GetValueFromBody(key string) interface{} {
	var (
		res map[string]interface{}
	)

	if len(t.responseValues) > 0 {
		res = t.responseValues[len(t.responseValues)-1]
	}
	if res == nil {
		return key
	}

	if strings.ToLower(key) == "${res}" {
		return json.NewJsonObjectWithMap(res).ToString()
	}
	if strings.HasPrefix(strings.ToLower(key), "${res.") {
		obj := json.NewJsonObjectWithMap(res).GetByKeys(key[6 : len(key)-1])
		if obj == nil {
			return ""
		}
		t := reflect.TypeOf(obj)
		if t.Name() == "" {
			objStr, _ := json.Marshal(obj)
			return string(objStr)
		}
		return fmt.Sprintf("%v", obj)
	}
	return key
}

func (t *TestContext) RenderValueString(str string) string {
	if str == "" || strings.HasPrefix(strings.ToLower(str), "${res.") || strings.ToLower(str) == "${res}" {
		return str
	}
	exp := regexp.MustCompile(`\$\{(.*?)\}`)
	params := exp.FindAllStringSubmatch(str, -1)
	for _, p := range params {
		if strings.HasPrefix(strings.ToLower(p[0]), "${res.") || strings.ToLower(p[0]) == "${res}" {
			continue
		}
		var v string
		if strings.ToLower(p[1]) == "space" {
			v = " "
		} else {
			v = t.GetAttrStr(p[1])
		}
		str = strings.Replace(str, p[0], v, -1)
	}
	return str
}

func (t *TestContext) parseRes(response *requests.Response) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	res["time"] = response.Time
	res["status"] = response.Status
	res["protocol"] = response.Proto
	res["length"] = len(response.Body)

	cookie := make(map[string]string)
	for _, c := range t.Session.HttpSession.Cookies {
		cookie[c.Name] = c.Value
	}
	res["cookie"] = cookie

	header := make(map[string]string)
	for k, v := range response.Headers {
		header[k] = v
	}
	res["header"] = header
	res["text"] = response.Body

	body := make(map[string]interface{})
	err := json.Unmarshal([]byte(response.Body), &body)
	res["body"] = body
	if err != nil {
		res["body"] = response.Body
	}

	return res, nil
}

func RunFuncByName(ctx *TestContext, cas *TestCase, fun *Function) {
	var (
		args = make([]reflect.Value, 0)
		tf   = &TestFunctions{}
	)

	TestHolder.Function = fun
	defer func() {
		TestHolder.Function = nil
	}()

	args = append(args, reflect.ValueOf(tf))
	fun.ArgsValue = make([]string, 0)
	for _, param := range fun.Args {
		arg := ctx.GetValueFromBody(ctx.RenderValueString(param))
		args = append(args, reflect.ValueOf(fmt.Sprintf("%v", arg)))
		fun.ArgsValue = append(fun.ArgsValue, fmt.Sprintf("%v", arg))
	}

	method, valid := reflect.TypeOf(tf).MethodByName(fun.Name)
	if !valid {
		return
	}

	if !cas.Pass {
		fun.Skip = true
		util.Log(util.Cyan("TEST:"), fmt.Sprintf("%s%s", util.IfStr(fun.IsBefore, "!", ""), fun.Name), fun.Args, "=>", "skip")
		return
	}
	methodTypeStr := fmt.Sprintf("%v", method.Type)
	methodArgs := strings.Split(methodTypeStr[strings.Index(methodTypeStr, "(")+1:strings.Index(methodTypeStr, ")")], ", ")

	if len(methodArgs) == len(args) {
		values := method.Func.Call(args)
		fun.Return = fmt.Sprintf("%v", values[0])
		if len(values) == 2 {
			if !values[1].IsNil() {
				fun.Err = values[1].Interface().(error)
			}
		}
		if fun.ReturnName != "" {
			ctx.SetAttr(fun.ReturnName, AttrValueTypeTemp, fun.Return)
		}
		cas.Pass = cas.Pass && fun.Err == nil
	} else if len(methodArgs) > len(args) {
		cas.Pass = false
		fun.Err = fmt.Errorf("call with to few input arguments")
	} else if len(methodArgs) < len(args) {
		cas.Pass = false
		fun.Err = fmt.Errorf("call with to many input arguments")
	}

	util.Log(util.Cyan("TEST:"), fmt.Sprintf("%s%s", util.IfStr(fun.IsBefore, "!", ""), fun.Name), fun.Args, "=>",
		util.IfStr(fun.Err == nil, fun.Return, fmt.Sprintf("Err: %s", fun.Err)))

}
