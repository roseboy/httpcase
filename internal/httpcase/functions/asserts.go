package functions

import (
	"fmt"
	"strings"
	"time"
)

type Asserts struct {
}

func (f *Asserts) Assert(v1 string, op string, v2 string) (bool, error) {
	switch strings.ToLower(op) {
	case "==":
		return f.AssertEq(v1, v2)
	case "!=":
		return f.AssertNe(v1, v2)
	case "<>":
		return f.AssertNe(v1, v2)
	case ">":
		return f.AssertGt(v1, v2)
	case ">=":
		return f.AssertGe(v1, v2)
	case "<":
		return f.AssertLt(v1, v2)
	case "<=":
		return f.AssertLe(v1, v2)
	case "like":
		return f.AssertContains(v1, v2)
	default:
		return false, fmt.Errorf("%s error", op)
	}
}

func (f *Asserts) AssertEq(v1 string, v2 string) (bool, error) {
	if v1 == v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertEq return false")
	}

}
func (f *Asserts) AssertNe(v1 string, v2 string) (bool, error) {
	if v1 != v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertNe return false")
	}
}

func (f *Asserts) AssertLt(v1 string, v2 string) (bool, error) {
	if v1 < v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertLt return false")
	}
}

func (f *Asserts) AssertGt(v1 string, v2 string) (bool, error) {
	if v1 > v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertGt return false")
	}
}
func (f *Asserts) AssertLe(v1 string, v2 string) (bool, error) {
	if v1 <= v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertLe return false")
	}
}

func (f *Asserts) AssertGe(v1 string, v2 string) (bool, error) {
	if v1 >= v2 {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertGe return false")
	}
}

func (f *Asserts) AssertContains(v1 string, v2 string) (bool, error) {
	if strings.Contains(v1, v2) {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertContains return false")
	}
}

func (f *Asserts) AssertMatch(v1 string, regStr string) (bool, error) {
	strFunc := &Strings{}
	if strFunc.Match(v1, regStr) {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertMatch return false")
	}
}

func (f *Asserts) AssertBefore(v1 string, v2 string) (bool, error) {
	d1, err := f.parseTime(v1)
	if err != nil {
		return false, err
	}
	d2, err := f.parseTime(v2)
	if err != nil {
		return false, err
	}
	if d1.Before(d2) {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertBefore return false")
	}

}

func (f *Asserts) AssertAfter(v1 string, v2 string) (bool, error) {
	d1, err := f.parseTime(v1)
	if err != nil {
		return false, err
	}
	d2, err := f.parseTime(v2)
	if err != nil {
		return false, err
	}

	if d1.After(d2) {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertAfter return false")
	}
}

func (f *Asserts) AssertEmpty(v1 string) (bool, error) {
	if v1 == "" {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertEmpty return false")
	}
}

func (f *Asserts) AssertNotEmpty(v1 string) (bool, error) {
	if v1 != "" {
		return true, nil
	} else {
		return false, fmt.Errorf("AssertNotEmpty return false")
	}
}

func (f *Asserts) parseTime(timeStr string) (t time.Time, err error) {
	timeFormats := []string{"2006-01-02 15:04:05", "2006-01-02", "15:04:05", "2006/01/02", "2006/01/02 15:04:05", "01-02-2006", "2006.01.02 15:04:05",
		"2006#01#02", "2006-01-02 15:04:05.000000", "2006-01-02 15:04:05.000000000",
		"2006-1-02", "2006-January-02", "2006-Jan-02", "2006-Jan-02", "06-Jan-02", "2006-01-02 15:04:05 Monday",
		"2006-01-02 Mon", "Mon 2006-01-2", "2006-01-02 3:4:5", "2006-01-02 3:4:5 PM", "2006-01-02 3:4:5 pm",
		time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850,
		time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano,
		time.Kitchen, time.Stamp, time.StampMilli, time.StampMicro, time.StampNano}
	for _, tf := range timeFormats {
		t, err = time.Parse(tf, timeStr)
		if err == nil {
			return t, nil
		}
	}
	return t, err
}
