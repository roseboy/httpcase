package json

import (
	"regexp"
)

type Rule struct {
	Expr    string
	Replace string
}

var Rules = []Rule{{"", ""},
	{`(?m)(\"[^\"]+\"):`, "\033[36m$1\033[39m:"},
	{`(?m)(^\s*[{}\]]{1}[,]*$|: [{\[][}\]]?)`, "\033[33m$1\033[39m"},
	{`: (\"[^\"]+\")`, ": \033[32m$1\033[39m$2"},
	{`(?m): ([\d][\d\.e+]*)([,]*)$`, ": \033[33m$1\033[39m$2"},
	{`(?m)(^\s+(?:[\d][\d\.e+]*|true|false|null))([,]*)$`, "\033[33m$1\033[39m$2"},
	{`(?::) ((?:true|false|null))`, ": \033[33m$1\033[39m"},
	{`(?m)^(true|false|null)$`, "\033[33m$1\033[39m"},
}

func Highlight(data []byte) ([]byte, error) {
	var (
		re  *regexp.Regexp
		err error
	)
	for _, rule := range Rules {
		if rule.Expr == "" {
			continue
		}

		re, err = regexp.Compile(rule.Expr)
		if err != nil {
			return nil, err
		}
		data = re.ReplaceAll(data, []byte(rule.Replace))
	}

	return data, nil
}
