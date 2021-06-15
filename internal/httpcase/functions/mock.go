package functions

import (
	"time"
)

type Mock struct {
}

func (f *Mock) NowDate() string {
	return time.Now().Format("2006-01-02")
}

func (f *Mock) NowTime() string {
	return time.Now().Format("15:04:05")
}

func (f *Mock) NowDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (f *Mock) Now(fmt string) string {
	return time.Now().Format(fmt)
}

func (f *Mock) RandStr(length string) string {
	return "KJGUZJBCEOSUjho"
}

func (f *Mock) RandEnName() string {
	return "Tom"
}

func (f *Mock) RandChName() string {
	return "王二虎"
}
