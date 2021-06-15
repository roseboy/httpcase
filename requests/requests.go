package requests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/roseboy/httpcase/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	HttpRequest   *http.Request `json:"-"`
	Method        string
	Url           string
	Body          string
	Params        map[string]string
	Headers       map[string]string
	AllowRedirect bool
	Timeout       int
	UploadFiles   map[string]*os.File
}

type Response struct {
	HttpResponse   *http.Response      `json:"-"`
	ResponseWriter http.ResponseWriter `json:"-"`
	Body           string
	Headers        map[string]string
	Status         int
	Time           int64
	Proto          string
	ContentLength  int64
}

func NewRequest() *RequestBuilder {
	return &RequestBuilder{Request: &Request{Headers: make(map[string]string), Params: make(map[string]string), AllowRedirect: true}}
}

func NewHttpSession() *RequestBuilder {
	return &RequestBuilder{
		HttpSession: &HttpSession{Cookies: make(map[string]*http.Cookie), Transport: &http.Transport{}},
		Request:     &Request{Headers: make(map[string]string), Params: make(map[string]string), AllowRedirect: true},
	}
}

func Get(url string) *RequestBuilder {
	return NewRequest().Get(url)
}

func Post(url string) *RequestBuilder {
	return NewRequest().Post(url)
}

func SendRequest(req *Request, session *HttpSession) (*Response, error) {
	var (
		resp       *http.Response
		response   = &Response{}
		transport  *http.Transport
		bodyBuffer io.Reader
		begin      = util.NowMillisecond()
	)

	if session != nil {
		transport = session.Transport
	} else {
		transport = &http.Transport{}
	}
	if req.Timeout > 0 {
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(network, addr, time.Duration(req.Timeout)*time.Millisecond)
			if err != nil {
				return nil, err
			}
			_ = conn.SetDeadline(time.Now().Add(time.Duration(req.Timeout) * time.Millisecond))
			return conn, nil
		}
	}

	if req.Method == "POST" && len(req.UploadFiles) > 0 {
		fileFiled := "file"
		for k := range req.UploadFiles {
			fileFiled = k
			break
		}
		file := req.UploadFiles[fileFiled]
		postBody := new(bytes.Buffer)
		writer := multipart.NewWriter(postBody)
		formFile, err := writer.CreateFormFile(fileFiled, file.Name())
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(formFile, file)
		if err != nil {
			return nil, err
		}

		for key, val := range req.Params {
			_ = writer.WriteField(key, val)
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}
		bodyBuffer = postBody
		req.Headers["Content-Type"] = writer.FormDataContentType()
	} else if req.Method == "POST" && req.Body != "" {
		if _, ok := req.Headers["Content-Type"]; !ok {
			req.Headers["Content-Type"] = "application/json"
		}
		if _, ok := req.Headers["Accept"]; !ok {
			req.Headers["Accept"] = "application/json"
		}
		bodyBuffer = strings.NewReader(req.Body)
	} else if req.Method == "POST" && len(req.Params) > 0 {
		if _, ok := req.Headers["Content-Type"]; !ok {
			req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
		}
		params := ""
		for k, v := range req.Params {
			params = fmt.Sprintf("%s&%s=%s", params, k, v)
		}
		bodyBuffer = strings.NewReader(params)
	} else if req.Method == "GET" {
		params := ""
		for k, v := range req.Params {
			params = fmt.Sprintf("%s&%s=%s", params, k, v)
		}
		if params != "" {
			if strings.Contains(req.Url, "?") {
				req.Url = fmt.Sprintf("%s&%s", strings.Trim(req.Url, "&"), strings.Trim(params, "&"))
			} else {
				req.Url = fmt.Sprintf("%s?%s", strings.Trim(req.Url, "?"), strings.Trim(params, "&"))
			}
		}
	}

	request, err := http.NewRequest(req.Method, req.Url, bodyBuffer)
	if err != nil {
		response.Time = util.NowMillisecond() - begin
		return response, err
	}

	for k, v := range req.Headers {
		//request.Header.Set(k, v)
		request.Header[k] = []string{v}
	}

	resp, err = transport.RoundTrip(request)

	if err != nil {
		response.Time = util.NowMillisecond() - begin
		return response, err
	}
	if req.AllowRedirect && resp.StatusCode == 302 {
		location := resp.Header.Get("Location")
		if strings.HasPrefix(location, "http") {
			req.Url = location
		} else if strings.HasPrefix(location, "/") {
			host := resp.Request.URL.String()[:len(resp.Request.URL.String())-len(resp.Request.URL.RequestURI())]
			req.Url = fmt.Sprintf("%s%s", host, location)
		} else if location != "" {
			host := resp.Request.URL.String()[:len(resp.Request.URL.String())-len(resp.Request.URL.RequestURI())]
			path := resp.Request.URL.Path
			path = path[:strings.LastIndex(path, "/")]
			req.Url = fmt.Sprintf("%s%s/%s", host, path, location)
		}
		return SendRequest(req, session)
	}

	header := make(map[string]string)
	for k := range resp.Header {
		v := resp.Header.Get(k)
		header[k] = v
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.Time = util.NowMillisecond() - begin
		return response, err
	}
	response.HttpResponse = resp
	response.Headers = header
	response.Body = string(body)
	response.Status = resp.StatusCode
	response.Proto = resp.Proto
	response.ContentLength = resp.ContentLength
	response.Time = util.NowMillisecond() - begin
	return response, nil
}

// 上传文件
// url                请求地址
// params        post form里数据
// nameField  请求地址上传文件对应field
// fileName     文件名
// file               文件
func PostFile(url string, params map[string]string, nameField string, fileName string, file io.Reader) ([]byte, error) {
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func IsHttpMethod(methodName string) (bool, string) {
	var httpMethods = []string{"GET", "POST", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "PATCH", "TRACE"}
	upperLine := strings.ToUpper(methodName)
	for _, m := range httpMethods {
		if strings.HasPrefix(upperLine, m) {
			return true, m
		}
	}
	return false, ""
}
