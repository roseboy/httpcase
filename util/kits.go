package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func InArrayInt(arr []int, ele int) bool {
	for _, e := range arr {
		if e == ele {
			return true
		}
	}
	return false
}

func InArrayString(arr []string, ele string) bool {
	for _, e := range arr {
		if e == ele {
			return true
		}
	}
	return false
}

func IfStr(ok bool, v1 string, v2 string) string {
	if ok {
		return v1
	}
	return v2
}

func TrimLeft(str string) string {
	chs := make([]int32, 0)
	ok := true
	for _, ch := range []rune(str) {
		if ch <= 32 && ok {
			continue
		}
		ok = false
		chs = append(chs, ch)
	}
	return string(chs)
}

func TrimRight(str string) string {
	chs := make([]int32, 0)
	rs := []rune(str)
	ok := true
	for i := len(rs) - 1; i >= 0; i-- {
		ch := rs[i]
		if ch <= 32 && ok {
			continue
		}

		ok = false
		chs = append(chs, ch)
	}

	for i, l := 0, len(chs); i < l/2; i++ {
		chs[i], chs[l-i-1] = chs[l-i-1], chs[i]
	}
	return string(chs)
}

func Trim(str string) string {
	return TrimLeft(TrimRight(str))
}

func DelRepeat(str string, ch string) string {
	chs := ""
	lastC := ""
	for _, c := range str {
		s := string(c)
		if s == ch && s == lastC {
			continue
		}
		chs = fmt.Sprintf("%s%s", chs, s)
		lastC = s
	}
	return chs
}

type Stack struct {
	stack map[int]string
}

func (s *Stack) Pop() string {
	v := s.stack[len(s.stack)-1]
	delete(s.stack, len(s.stack)-1)
	return v
}

func (s *Stack) Push(v string) {
	if s.stack == nil {
		s.stack = make(map[int]string)
	}
	s.stack[len(s.stack)] = v
}

func (s *Stack) Top() string {
	return s.stack[len(s.stack)-1]
}

func (s *Stack) Length() int {
	return len(s.stack)
}

func (s *Stack) IsEmpty() bool {
	return len(s.stack) == 0
}

func (s *Stack) Empty() {
	s.stack = make(map[int]string)
}

func FmtDate(time time.Time, fmt string) string {
	return time.Format(fmt)
}

func ReadText(path string) (string, error) {
	var text = ""

	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		text = fmt.Sprintf("%s%s\n", text, line)
	}
	return text, nil
}

func WriteText(text string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = file.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}

func ParseUrlParam(url string) map[string]string {
	param := make(map[string]string)
	if strings.Contains(url, "?") {
		url = url[strings.Index(url, "?"):]
	}

	kvs := strings.Split(url, "&")
	for _, kv := range kvs {
		p := strings.Split(kv, "=")
		param[p[0]] = p[1]
	}
	return param
}

func NowMillisecond() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}
