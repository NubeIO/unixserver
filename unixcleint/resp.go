package unixclient

import (
	"errors"
	"fmt"
)

type Response struct {
	Status   string // Success Error
	Data     interface{}
	Error    interface{}
	AsString string
	AsNumber float64
	AsBool   bool
}

func (r *Response) GetString() string {
	return r.AsString
}

func (r *Response) GetNumber() float64 {
	return r.AsNumber
}

func (r *Response) GetBool() bool {
	return r.AsBool
}

func (r *Response) IsError() bool {
	if r.Error != "" {
		return true
	}
	return false
}

func (r *Response) IsOk() bool {
	if r.Error != "" {
		return false
	}
	return true
}

func (r *Response) GetError() error {
	if r.Error != "" {
		return errors.New(fmt.Sprint(r.Error))
	}
	return nil
}

// GetData returns the Data field from the response.
func (r *Response) GetData() interface{} {
	return r.Data
}

func errorResp(resp *Response, err error) *Response {
	if resp == nil {
		resp = &Response{}
		resp.Error = fmt.Sprintf("reponse was empty")
		return resp
	}
	if err != nil {
		resp.Error = fmt.Sprintf("error sending request: %v", err)
		return resp
	}
	return resp
}
