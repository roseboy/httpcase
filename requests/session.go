package requests

import (
	"fmt"
	"net/http"
)

type HttpSession struct {
	Cookies   map[string]*http.Cookie
	Transport *http.Transport
}

func (session *HttpSession) SetCookie(cookies []*http.Cookie) {
	if session.Cookies == nil {
		session.Cookies = make(map[string]*http.Cookie)
	}
	for _, c := range cookies {
		session.Cookies[c.Name] = c
	}
}

func (session *HttpSession) GetCookieString() string {
	str := ""
	for _, c := range session.Cookies {
		str = fmt.Sprintf("%s; %s=%s", str, c.Name, c.Value)
	}
	if len(str) > 2 {
		return str[2:]
	}
	return str
}
