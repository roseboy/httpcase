package requests

import (
	"github.com/roseboy/httpcase/json"
)

type ResponseParser struct {
	Response *Response
	Err      error
}

func (rp *ResponseParser) ReadToText() (string, error) {
	if rp.Err != nil {
		return "", rp.Err
	}
	return rp.Response.Body, nil
}

func (rp *ResponseParser) ReadToJsonObject() (*json.Object, error) {
	if rp.Err != nil {
		return nil, rp.Err
	}
	return json.NewJsonObject(rp.Response.Body)
}
