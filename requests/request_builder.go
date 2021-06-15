package requests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type RequestBuilder struct {
	HttpSession *HttpSession
	Request     *Request
}

func (rb *RequestBuilder) WithHttpRequest(request *http.Request) *RequestBuilder {
	defer request.Body.Close()
	body, _ := ioutil.ReadAll(request.Body)

	header := make(map[string]string)
	for k := range request.Header {
		v := request.Header.Get(k)
		header[k] = v
	}

	rb.Request.HttpRequest = request
	rb.Request.Url = request.URL.RequestURI()
	rb.Request.Method = request.Method
	rb.Request.Headers = header
	rb.Request.Body = string(body)
	return rb
}

func (rb *RequestBuilder) Send() *ResponseParser {
	if rb.HttpSession != nil {
		cookie, ok := rb.Request.Headers["Cookie"]
		if !ok {
			cookie = rb.Request.Headers["cookie"]
		}
		if cookie != "" {
			rb.Header("Cookie", strings.Trim(fmt.Sprintf("%s;%s", cookie, rb.HttpSession.GetCookieString()), ";"))
		} else {
			rb.Header("Cookie", rb.HttpSession.GetCookieString())
		}
	}
	rp := &ResponseParser{}
	rp.Response, rp.Err = SendRequest(rb.Request, rb.HttpSession)

	if rb.HttpSession != nil && rp.Err == nil {
		rb.HttpSession.SetCookie(rp.Response.HttpResponse.Cookies())
	}

	return rp
}

func (rb *RequestBuilder) SendWithRequest(req *Request) *ResponseParser {
	rb.Request = req
	return rb.Send()
}

func (rb *RequestBuilder) Url(url string) *RequestBuilder {
	rb.Request.Url = url
	return rb
}

func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.Request.Method = method
	return rb
}

func (rb *RequestBuilder) Get(url string) *RequestBuilder {
	rb.Request.Method = "GET"
	rb.Request.Url = url
	return rb
}

func (rb *RequestBuilder) Post(url string) *RequestBuilder {
	rb.Request.Method = "POST"
	rb.Request.Url = url
	return rb
}

func (rb *RequestBuilder) Body(body string) *RequestBuilder {
	rb.Request.Body = body
	return rb
}

func (rb *RequestBuilder) Param(key string, value string) *RequestBuilder {
	rb.Request.Params[key] = value
	return rb
}

func (rb *RequestBuilder) Params(params map[string]string) *RequestBuilder {
	rb.Request.Params = params
	return rb
}

func (rb *RequestBuilder) ParamString(params string) *RequestBuilder {
	if params == "" {
		return rb
	}
	ps := strings.Split(params, "&")
	for _, p := range ps {
		kv := strings.Split(p, "=")
		if len(kv) != 2 {
			return rb
		}
		rb.Param(kv[0], kv[1])
	}
	return rb
}

func (rb *RequestBuilder) Header(key string, value string) *RequestBuilder {
	rb.Request.Headers[key] = value
	return rb
}

func (rb *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	rb.Request.Headers = headers
	return rb
}

func (rb *RequestBuilder) AllowRedirect(allow bool) *RequestBuilder {
	rb.Request.AllowRedirect = allow
	return rb
}

func (rb *RequestBuilder) Timeout(timeout int) *RequestBuilder {
	rb.Request.Timeout = timeout
	return rb
}

func (rb *RequestBuilder) File(name string, file *os.File) *RequestBuilder {
	if rb.Request.UploadFiles == nil {
		rb.Request.UploadFiles = make(map[string]*os.File)
	}
	rb.Request.UploadFiles[name] = file
	return rb
}

func (rb *RequestBuilder) Files(files map[string]*os.File) *RequestBuilder {
	rb.Request.UploadFiles = files
	return rb
}

func (rb *RequestBuilder) GetRequest() *Request {
	return rb.Request
}
