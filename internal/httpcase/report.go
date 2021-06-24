package httpcase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/roseboy/httpcase/requests"
	"html/template"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/roseboy/httpcase/util"
)

func PrintReport(testCtx *TestContext, result *TestResult) error {
	if result.Skipped > 0 {
		return nil
	}
	err := PrintTextReport(result)
	if err != nil {
		return err
	}

	if testCtx.Out != "" {
		err := PrintHtmlReport(testCtx, result)
		if err != nil {
			return err
		}
	}

	if testCtx.CallbackUrl != "" {
		err := ReportCallback(testCtx, result)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintTextReport(result *TestResult) error {
	var (
		cases = result.TestCases
	)

	util.Println()
	util.Println("---------------------------------------------------------------------")
	util.Println(fmt.Sprintf("  Test Result (Total:%s, Pass:%s, Fail:%s, Skip:%s, Duration:%dms)",
		util.Blue(result.Total), util.Green(result.Passed), util.Red(result.Failed), util.Cyan(result.Skipped), result.Duration))
	util.Println("---------------------------------------------------------------------")
	i := 0
	for _, cas := range cases {
		if cas.IsGlobal {
			continue
		}
		i++
		util.Println(fmt.Sprintf("[%d]", i), util.IfStr(cas.Pass, "["+util.Green("Pass")+"]", "["+util.Red("Fail")+"]"), cas.Name, cas.Request.Method, cas.Request.Url)
	}
	if i == 0 {
		return nil
	}
	util.Println("---------------------------------------------------------------------")
	return nil
}

func PrintHtmlReport(testCtx *TestContext, result *TestResult) error {
	var (
		cases = result.TestCases
	)

	data := make(map[string]interface{})
	data["Cases"] = cases
	data["Date"] = util.FmtDate(time.Now(), "2006-01-02 15:04:05")
	data["FileName"] = result.TestCaseFile.FileName
	data["Env"] = testCtx.Env
	data["Vars"] = testCtx.Attr
	data["Total"] = result.Total
	data["Pass"] = result.Passed
	data["Fail"] = result.Failed
	data["Skip"] = result.Skipped
	data["Duration"] = result.Duration

	html := getHtmlTemplate()
	//isProd := strings.Contains(html, "<!--prod-->")
	html = trimHtml(html)
	html, err := renderHtml(html, data)
	if err != nil {
		return err
	}

	//html = testCtx.RenderValueString(html)

	htmlPath := testCtx.Out
	if !strings.HasSuffix(strings.ToLower(html), ".html") && !strings.HasSuffix(strings.ToLower(html), ".htm") {
		htmlPath = fmt.Sprintf("%s.html", testCtx.Out)
	}

	err = util.WriteText(html, htmlPath)
	if err != nil {
		return err
	}
	return nil
}

func ReportCallback(testCtx *TestContext, result *TestResult) error {
	var (
		url            = testCtx.CallbackUrl
		jsFunctionName = testCtx.CallbackJsFunction
	)

	resultBody, err := json.Marshal(result)
	if err != nil {
		return err
	}
	resultBodyStr := string(resultBody)

	if jsFunctionName != "" {
		tf := &TestFunctions{}
		resultBodyStr, err = tf.RunJsWithName(jsFunctionName, resultBodyStr)
		if err != nil {
			return err
		}
	}

	rb := requests.Post(url)
	if strings.HasPrefix(resultBodyStr, "?") {
		rb.ParamString(resultBodyStr)
	} else {
		rb.Body(resultBodyStr)
	}
	res := rb.Send()

	if res.Err != nil {
		return res.Err
	}
	if res.Response.Status != 200 {
		return fmt.Errorf("ReportCallback:%s return code %d", url, res.Response.Status)
	}

	return nil
}

//渲染html
func renderHtml(html string, data map[string]interface{}) (string, error) {
	funcs := template.FuncMap{"mod": func(i int, m int) int { return i % m },
		"html": func(str string) string {
			return str
		},
		"rpl": func(str string, old string, new string) string {
			return strings.Replace(str, old, new, -1)
		},
		"ifstr": func(ok bool, a string, b string) string {
			if ok {
				return a
			} else {
				return b
			}
		},
		"json": func(v interface{}) string {
			str, _ := json.Marshal(v)
			return string(str)
		},
		"fmtjson": func(data string) string {
			var str bytes.Buffer
			err := json.Indent(&str, []byte(data), "", "  ")
			if err != nil {
				return data
			}
			return str.String()
		},
		"fmtnum": func(a int64, n int, m int) string {
			return fmt.Sprintf(fmt.Sprintf("%s.%df", "%", n), float64(a)*math.Pow(float64(10), float64(m)))
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	}
	t, err := template.New("report.html").Funcs(funcs).Parse(html)
	if err != nil {
		return "", err
	}

	htmlData := bytes.Buffer{}
	err = t.Execute(&htmlData, data)
	if err != nil {
		return "", err
	}
	return htmlData.String(), nil
}

func getHtmlTemplate() string {
	html, err := util.ReadText("resource/report-tpl.html")
	if err != nil {
		html = HtmlTemplate + "\r<!--prod-->"
	}

	return html
}

func trimHtml(html string) string {
	html = strings.ReplaceAll(html, "\r", "")
	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "\t", "")
	reg := regexp.MustCompile(`\s+`)
	html = reg.ReplaceAllString(html, " ")
	html = strings.ReplaceAll(html, "> <", "><")
	return html
}

const HtmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Test Report ({{.Date}} - HttpCase v1.0</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style type="text/css">
        html {
            -ms-text-size-adjust: 100%;
            -webkit-text-size-adjust: 100%;
        }

        body {
            margin: 10px 0 0 0;
            background: #E9EDF0;
        }

        table {
            border-collapse: collapse;
            border-spacing: 0;
        }

        td, th {
            padding: 0;
        }

        .content {
            width: 1100px;
            margin: 0 auto;
            background: #fff;
            border: 1px solid #E9EDF0;
            padding: 20px;
        }

        .title {
            font-size: 22px;
            padding: 10px 10px 10px 0;
            font-weight: bold;
        }

        .title2 {
            font-size: 20px;
            padding: 0 5px 2px 0;
            margin-top: 30px;
            font-weight: bold;
            margin-bottom: 10px;
        }

        .title3 {
            font-size: 16px;
            padding: 0 5px 5px 0;
            font-weight: bold;
            color: #464646;
            line-height: 28px;
        }

        .title4 {
            padding: 3px;
            color: #6c6c6d;
        }

        .pure-table {
            border-collapse: collapse;
            border-spacing: 0;
            empty-cells: show;
            border: 1px solid #cdcdcd;
            border-width: 1px 0 1px 0;
        }

        .pure-table caption {
            color: #000;
            font: italic 85%/1 arial, sans-serif;
            padding: 1em 0;
            text-align: center;
        }

        .pure-table td, .pure-table th {
            border-bottom: 1px solid #e8e8e8;
            font-size: inherit;
            margin: 0;
            overflow: visible;
            padding: .3em 1em;
        }

        .pure-table thead {
            color: #0e0e0e;
            text-align: left;
            vertical-align: bottom;
        }

        .pure-table td {
            background-color: #FEFFFE;
            color: #505050;
            word-break: break-all;
        }

        .pass {
            color: #2dbf6e;
            font-weight: bold;
        }

        .unpass {
            color: #ff0000;
            font-weight: bold;
        }

        .skip {
            color: #3E9EFF;
            font-weight: bold;
        }

        .blue {
            color: #00f;
        }

        .text {
            color: #0e0e0e;
            padding: 0 10px 10px 10px;
        }

        .footer {
            text-align: center;
            padding: 15px;
        }

        a {
            color: #505050;
        }

        pre {
            display: block;
            padding: 10px;
            line-height: 16px;
            font-size: 13px;
            border: 1px solid #d9d9d9;
            white-space: pre-wrap;
            background: #f6f6f6;
            overflow: auto;
            max-height: 200px;
            margin: 0;
            margin-bottom: 10px;
        }

        .flow-block {
            margin-top: 10px;
            padding: 10px 10px 0 10px;
            border: 1px solid #e5e5e5;
            position: relative;
            display: inline-block;
            width: 1080px;
            margin-bottom: 10px;
        }

        .flow-block > h3:first-child {
            display: inline-block;
            padding: 5px;
            position: absolute;
            top: -18px;
            left: 15px;
            font-size: inherit;
            font-weight: bold;
            margin: 0;
            background: #fff;
            color: #0e0e0e;
        }

        hr {
            width: 1140px;
            margin-left: -20px;
            border: 2px solid #e9edf0;
        }

        .big-hr {
            border: 5px solid #e9edf0;
        }
    </style>
</head>
<body>
<div class="content">
    <div class="title">Test Report for [{{.FileName}}] ({{.Date}})</div>
    <div class="title2">Enviroment</div>
    <table class="pure-table" width="100%">
        <thead>
        <tr>
            <th width="20%">Var</th>
            <th>Value</th>
        </tr>
        </thead>
        <tr>
            <td>env</td>
            <td>{{.Env}}</td>
        </tr>
        {{ range $key,$value := .Vars }}
        {{ if eq $value.Type 2}}
        <tr>
            <td>{{ $key }}</td>
            <td>{{ $value.Value }}</td>
        </tr>
        {{ end }}
        {{ end }}
        {{ range $key,$value := .Vars }}
            {{ if eq $value.Type 1}}
                <tr>
                    <td>{{ $key }}</td>
                    <td>{{ $value.Value }}</td>
                </tr>
            {{ end }}
        {{ end }}
    </table>
    <br/>
    <hr class="big-hr"/>
    <div class="title2">Summary</div>
    <div class="text">Total: <span class="blue">{{ .Total}}</span> , Pass: <span class="pass">{{ .Pass}}</span> , Fail:
        <span class="unpass">{{ .Fail}}</span> , Skip: <span class="skip">{{ .Skip}}</span> , Duration: {{ .Duration }}
        ms .
    </div>
    <table class="pure-table" width="100%">
        <thead>
        <tr>
            <th width="9%">Result</th>
            <th width="25%">Name</th>
            <th>Url</th>
            <th width="10%">Duration</th>
        </tr>
        </thead>
        {{ range $i,$item := .Cases}}
        {{ if not .IsGlobal }}
        <tr>
            <td>{{html (ifstr .Pass "[ <span class=\"pass\">Pass</span> ]" "[ <span class=\"unpass\">Fail</span> ]")}}
            </td>
            <td><a href="#{{.Name}}">{{.Name}}</a></td>
            <td>{{.Request.Method}} {{.Request.Url}}</td>
            <td>{{.Time}}ms</td>
        </tr>
        {{ end }}
        {{ end }}
    </table>
    <br/>
    <hr class="big-hr"/>
    <div class="title2">Detail</div>
    {{ range $i,$item := .Cases}}
    {{ if not .IsGlobal }}
    <div class="title3" id="{{.Name}}">
        {{ if .Pass}}
        [ <span class="pass">Pass</span> ] {{.Name}}
        {{ else }}
        [ <span class="unpass">Fail</span> ] {{.Name}}
        {{ end }}
        <br/>{{.Request.Method}} {{.Request.Url}} ==>
        {{ if .Err}}ERROR{{ else }}{{.Res.status}}{{ end }}
    </div>
    <table class="pure-table" width="100%">
        <thead>
        <tr>
            <th width="6%">Result</th>
            <th>Test</th>
        </tr>
        </thead>
        <tbody>
        {{ range $ii,$fun := .Functions}}
        {{ if eq .ReturnName ""}}
        <tr>
            {{if .Skip}}
            <td><span class="skip">Skip</span></td>
            {{ else if eq .Return "false"}}
            <td><span class="unpass">Fail</span></td>
            {{ else }}
            <td><span class="pass">OK</span></td>
            {{ end }}
            <td>{{.Code.Code}}</td>
        </tr>
        {{ end }}
        {{ end }}
        </tbody>
    </table>
    <br/>
    <div class="flow-block">
        <h3>Request</h3>
        <div class="title4">Headers</div>
        <pre>{{fmtjson (json .Request.Headers)}}</pre>
        <div class="title4">{{ifstr (eq .Request.Body "") "Params" "Body"}}</div>
        <pre>{{ifstr (eq .Request.Body "") (fmtjson (json .Request.Params)) (fmtjson .Request.Body)}}</pre>
    </div>
    <br/>
    <div class="flow-block">
        <h3>Response</h3>
        <div class="title4">Headers</div>
        <pre>{{fmtjson (json .Res.header)}}</pre>
        <div class="title4">Body</div>
        {{ if .Res }}
        <pre>{{fmtjson .Res.text}}</pre>
        {{ else }}
        <pre></pre>
        {{ end }}
    </div>
    <br/>
    <hr/>
    {{ end }}
    {{ end }}

</div>
<div class="footer">
    Report generated on {{.Date}} by
    <a target="_blank" href="https://github.com/roseboy/httpcase">HttpCase</a>
    v1.0
</div>
</body>
</html>
`
