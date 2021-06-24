package requests

import (
	"github.com/roseboy/httpcase/json"
	"io"
	"io/ioutil"
	"os"
)

type ResponseParser struct {
	Response *Response
	Err      error
}

func (rp *ResponseParser) ReadToText() (string, error) {
	if rp.Err != nil {
		return "", rp.Err
	}
	return rp.body()
}

func (rp *ResponseParser) ReadToJsonObject() (*json.Object, error) {
	if rp.Err != nil {
		return nil, rp.Err
	}
	body, err := rp.body()
	if err != nil {
		return nil, err
	}
	return json.NewJsonObject(body)
}

func (rp *ResponseParser) SaveToFile(path string) error {
	if rp.Err != nil {
		return rp.Err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, *rp.Response.Body)

	return err
}

func (rp *ResponseParser) body() (string, error) {
	var reader = *rp.Response.Body
	defer func() { _ = reader.Close() }()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
