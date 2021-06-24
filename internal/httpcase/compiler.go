package httpcase

import (
	"bufio"
	"fmt"
	"github.com/roseboy/httpcase/requests"
	"github.com/roseboy/httpcase/util"
	"os"
	"path/filepath"
	"strings"
)

type Compiler struct {
	Codes       []*CodeLine
	testContext *TestContext
}

type TestRequest struct {
	TestCaseFile *CodeFile
	TestCases    []*TestCase
}

type TestResult struct {
	TestCaseFile *CodeFile
	TestCases    []*TestCase
	Passed       int
	Failed       int
	Skipped      int
	Duration     int64
	Total        int
}

type TestCase struct {
	Name      string //用例名称
	Request   *requests.Request
	Response  *requests.Response
	Res       map[string]interface{} `json:"-"` //返回内容,包含返回header
	Functions []*Function            //测试函数
	IsGlobal  bool                   //本用例只有函数执行
	Codes     []*CodeLine            `json:"-"` //用例原始代码
	While     bool                   `json:"-"` //循环条件
	If        bool                   `json:"-"` //是否执行条件
	Skip      bool                   //是否跳过
	Pass      bool                   //是否通过
	Time      int64                  //耗时ms
	Err       error                  //执行错误
}

type CodeLine struct {
	Ok   bool `json:"-"`
	Line int
	Code string
	File *CodeFile
}

type CodeFile struct {
	File     *os.File `json:"-"`
	Name     string
	FileName string
	Dir      string
}

type Function struct {
	Name       string
	Args       []string
	ArgsValue  []string
	Skip       bool
	IsBefore   bool
	Return     string
	Err        error
	ReturnName string
	Code       *CodeLine `json:"-"`
	Functions  []*Function
}

func NewCompiler(testCtx *TestContext, codes []*CodeLine) *Compiler {
	return &Compiler{
		Codes:       codes,
		testContext: testCtx,
	}
}

func (c *Compiler) Compile() (*TestRequest, error) {
	var (
		err         error
		lines       = c.Codes
		testCases   = make([]*TestCase, 0)
		testRequest = &TestRequest{}
	)

	if len(lines) > 0 {
		testRequest.TestCaseFile = lines[0].File
	}

	//匹配全局函数
	globalFuncs, err := c.parseGlobalFuncs(lines)
	if err != nil {
		return nil, err
	}

	//匹配case
	var testCase *TestCase
	var fileName = ""
	for _, line := range lines {
		if fileName != line.File.Name {
			fileName = line.File.Name
			if gFuns, ok := globalFuncs[fileName]; ok && len(gFuns) > 0 {
				delete(globalFuncs, fileName)
				funcCase := &TestCase{
					Codes:     make([]*CodeLine, 0),
					Name:      "Global(" + line.File.FileName + ")",
					Functions: gFuns,
					Pass:      true,
					If:        true,
					IsGlobal:  true,
					Skip:      false,
				}
				testCases = append(testCases, funcCase)
			}
		}

		if line.Ok {
			continue
		}

		if strings.HasPrefix(line.Code, "@") {
			testCase = &TestCase{Codes: make([]*CodeLine, 0), Pass: true, If: true, IsGlobal: false, Skip: false}
			testCases = append(testCases, testCase)
			testCase.Codes = append(testCase.Codes, line)
		} else if testCase != nil {
			testCase.Codes = append(testCase.Codes, line)
		}
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("can't read test file")
	}

	//匹配case内容
	for _, testCase := range testCases {
		if testCase.IsGlobal {
			continue
		}

		codes := testCase.Codes
		testCase.Name = codes[0].Code[1:]
		codes[0].Ok = true

		method, url, err := c.parseUrl(codes)
		if err != nil {
			return nil, err
		}

		headers, err := c.parseHeader(codes)
		if err != nil {
			return nil, err
		}

		params, err := c.parseParam(codes)
		if err != nil {
			return nil, err
		}

		var body string
		if params == "" {
			body, err = c.parseBody(codes)
			if err != nil {
				return nil, err
			}
		}

		testCase.Request = requests.NewRequest().AllowRedirect(false).Url(url).Method(method).ParamString(params).
			Body(body).Headers(headers).GetRequest()

		testCase.Functions, err = c.parseFunctions(codes)
		if err != nil {
			return nil, err
		}
	}

	//匹配要执行的case
	if c.testContext.Tags != "" {
		tags := strings.Split(c.testContext.Tags, ",")
		for _, testCase := range testCases {
			if testCase.IsGlobal {
				continue
			}
			testCase.Skip = true
			for _, tag := range tags {
				if strings.Contains(testCase.Name, tag) {
					testCase.Skip = false
					break
				}
			}
		}
	}

	testRequest.TestCases = testCases
	return testRequest, nil
}

func ReadCaseFile(path string) ([]*CodeLine, error) {
	var (
		codes = make([]*CodeLine, 0)
	)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	codeFile := &CodeFile{
		File: f,
		Name: f.Name(),
	}
	codeFile.Dir, codeFile.FileName = filepath.Split(f.Name())

	lineNum := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		line = util.Trim(line)
		line = strings.Trim(line, "#")

		if strings.Contains(line, " //") {
			line = line[:strings.Index(line, " //")]
		}
		line = util.Trim(line)
		if line == "" {
			continue
		}

		//include
		lineLower := strings.ToLower(line)
		if strings.HasPrefix(lineLower, "include") || strings.HasPrefix(lineLower, "!include") ||
			strings.HasPrefix(lineLower, "(include") || strings.HasPrefix(lineLower, "!(include") {
			codes = append(codes, &CodeLine{File: codeFile, Line: lineNum, Code: line, Ok: true})
			lineTrim := strings.Trim(line, "(")
			lineTrim = strings.Trim(lineTrim, ")")
			path := strings.Split(lineTrim, " ")
			if len(path) > 1 {
				dir, _ := filepath.Split(f.Name())
				includeCase, err := ReadCaseFile(dir + path[1])
				if err != nil {
					return codes, nil
				}
				codes = append(codes, includeCase...)
			}
		} else {
			codes = append(codes, &CodeLine{File: codeFile, Line: lineNum, Code: line, Ok: strings.HasPrefix(line, "//")})
		}
	}

	if err := scanner.Err(); err != nil {
		return codes, err
	}

	return codes, nil
}

//匹配全局配置
func (c *Compiler) parseGlobalFuncs(lines []*CodeLine) (map[string][]*Function, error) {
	var (
		fileName        = ""
		inCase          = false
		globalFuncCodes = make(map[string][]*CodeLine)
		globalFuncs     = make(map[string][]*Function)
	)
	for _, line := range lines {
		if line.Ok {
			continue
		}
		if fileName != line.File.Name {
			fileName = line.File.Name
			inCase = false
		}
		if inCase {
			continue
		}
		if strings.HasPrefix(line.Code, "@") {
			inCase = true
		}
		if !inCase {
			if _, ok := globalFuncCodes[line.File.Name]; !ok {
				globalFuncCodes[line.File.Name] = make([]*CodeLine, 0)
			}
			globalFuncCodes[line.File.Name] = append(globalFuncCodes[line.File.Name], line)
		}
	}

	for fileName, codes := range globalFuncCodes {
		funcs, err := c.parseFunctions(codes)
		if err != nil {
			return nil, err
		}
		globalFuncs[fileName] = funcs
	}

	return globalFuncs, nil

}

func (c *Compiler) parseUrl(codes []*CodeLine) (string, string, error) {
	for _, line := range codes {
		if line.Ok {
			continue
		}
		if ok, method := requests.IsHttpMethod(line.Code); ok {
			url := strings.Trim(line.Code[len(method):], " ")
			line.Ok = true
			return method, url, nil
		}
	}
	return "", "", fmt.Errorf("can't parse the http url (file:%s line:%d)", codes[0].File.Name, codes[0].Line)
}

func (c *Compiler) parseHeader(codes []*CodeLine) (map[string]string, error) {
	header := make(map[string]string)
	for _, line := range codes {
		if line.Ok {
			continue
		}
		if strings.HasPrefix(line.Code, "{") {
			break
		}
		if strings.Contains(line.Code, ": ") && !strings.Contains(line.Code, ":\"") {
			hd := strings.Split(line.Code, ": ")
			header[util.Trim(hd[0])] = util.Trim(hd[1])
			line.Ok = true
		} else if len(header) > 0 {
			break
		}
	}

	return header, nil
}

func (c *Compiler) parseParam(codes []*CodeLine) (string, error) {
	var (
		param string
	)
	for _, line := range codes {
		if line.Ok {
			continue
		}
		if strings.Contains(line.Code, "=") &&
			!strings.Contains(line.Code, " ") &&
			!strings.Contains(line.Code, "\"") &&
			!strings.Contains(line.Code, ":") {
			param = fmt.Sprintf("%s%s", param, line.Code)
			line.Ok = true
		}
	}
	return param, nil
}

func (c *Compiler) parseBody(codes []*CodeLine) (string, error) {
	var (
		body   string
		inBody = false
		stack  = &util.Stack{}
	)

	for _, line := range codes {
		if line.Ok {
			continue
		}

		if strings.HasPrefix(line.Code, "{") {
			inBody = true
		}

		if !inBody {
			continue
		}

		line.Ok = true

		for _, ch := range line.Code {
			s := string(ch)
			if s == "{" || (s == "}" && stack.Top() == "}") {
				stack.Push(s)
			} else if s == "}" && stack.Top() == "{" {
				stack.Pop()
				body = fmt.Sprintf("%s%s", body, s)
			}

			if !stack.IsEmpty() && s != "}" {
				body = fmt.Sprintf("%s%s", body, s)
			}

			if stack.IsEmpty() {
				return body, nil
			}
		}
	}

	return "", nil
}

func (c *Compiler) parseFunctions(codes []*CodeLine) ([]*Function, error) {
	var (
		functions      = make([]*Function, 0)
		funcStr        = make([]string, 0)
		funcRetrunName = make(map[int]string)
		funcCodeLine   = make(map[int]*CodeLine)
		stack          = &util.Stack{}
		codeFile       *CodeFile
	)

	for _, line := range codes {
		if codeFile != nil && codeFile.Name != line.File.Name {
			break
		}
		codeFile = line.File
		if line.Ok {
			continue
		}
		line.Ok = true
		lineCode := line.Code
		isBefore := strings.HasPrefix(lineCode, "!")

		if isBefore {
			lineCode = lineCode[1:]
		}
		if !strings.Contains(lineCode, "(") && !strings.Contains(lineCode, ")") {
			funcStr = append(funcStr, fmt.Sprintf("%s%s", util.IfStr(isBefore, "!", ""), lineCode))
			funcCodeLine[len(funcStr)-1] = line
			continue
		}

		chs := ""
		ignoreBrackets := false
		stack.Empty()
		for i := 0; i < len(lineCode); i++ {
			s := lineCode[i : i+1]
			if ignoreBrackets {
				if i == len(lineCode)-stack.Length() && s == ")" && stack.Top() == "(" {
					ignoreBrackets = false
				}
			}

			if s == "(" && !ignoreBrackets {
				if chs != "" {
					stack.Push(chs)
					chs = ""
				}
				stack.Push(s)
			} else if s == ")" && !ignoreBrackets {
				pps := stack.Pop()
				for pps != "(" && !stack.IsEmpty() {
					chs = fmt.Sprintf("%s%s", pps, chs)
					pps = stack.Pop()
				}

				funcStr = append(funcStr, fmt.Sprintf("%s%s", util.IfStr(isBefore, "!", ""), chs))
				funcCodeLine[len(funcStr)-1] = line
				chs = ""
				if !stack.IsEmpty() {
					c.testContext.VarIndex++
					varName := fmt.Sprintf("${&0x%X}", c.testContext.VarIndex)
					stack.Push(varName)
					funcRetrunName[len(funcStr)-1] = fmt.Sprintf("&0x%X", c.testContext.VarIndex)
				}
			} else if s == ">" && strings.HasSuffix(chs, "-") {
				chs = fmt.Sprintf("%s%s", chs, s)
				ignoreBrackets = true
			} else {
				chs = fmt.Sprintf("%s%s", chs, s)
			}
		}
		if chs != "" {
			stack.Push(chs)
			chs = ""
		}
		for !stack.IsEmpty() {
			chs = fmt.Sprintf("%s%s", stack.Pop(), chs)
		}
		if chs != "" {
			funcStr = append(funcStr, fmt.Sprintf("%s%s", util.IfStr(isBefore, "!", ""), chs))
			funcCodeLine[len(funcStr)-1] = line
		}
	}

	for i, fs := range funcStr {
		fs = util.DelRepeat(fs, " ")
		isBefore := strings.HasPrefix(fs, "!")
		fs = util.IfStr(isBefore, fs[1:], fs)

		if strings.Index(strings.Split(fs, " ")[0], ".") > 0 {
			fs = fmt.Sprintf("RunJs %s", strings.Replace(fs, " ", ",", -1))
		}
		if ok, funName := GetFunctionName(fs); ok {
			function := &Function{
				Name:       funName,
				ArgsValue:  make([]string, 0),
				ReturnName: funcRetrunName[i],
				Code:       funcCodeLine[i],
				IsBefore:   isBefore,
			}
			argsStr := util.Trim(fs[len(funName):])
			if argsStr != "" {
				function.Args = strings.Split(argsStr, " ")
			}
			functions = append(functions, function)
		}
	}

	for _, fun := range functions {
		for i, arg := range fun.Args {
			if arg == "->" {
				codeStrs := strings.Split(strings.Join(fun.Args[i+1:], " "), "|")
				subCodes := make([]*CodeLine, 0)
				for _, codeStr := range codeStrs {
					subCodes = append(subCodes, &CodeLine{File: codeFile, Line: fun.Code.Line, Ok: false, Code: strings.Trim(codeStr, " ")})
				}
				subFuncs, err := c.parseFunctions(subCodes)
				if err != nil {
					return nil, err
				}
				fun.Functions = subFuncs
				fun.Args = fun.Args[:i]
				break
			}
		}
	}
	return functions, nil
}
