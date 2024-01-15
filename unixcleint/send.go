package unixclient

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

func (uc *UnixClient) Send(path string, dataToSend interface{}, timeoutInSeconds int, expectedResponse interface{}, expectedType string) (*Response, error) {
	deadline := time.Now().Add(time.Duration(timeoutInSeconds) * time.Second)
	uc.conn.SetDeadline(deadline)

	data, err := json.Marshal(dataToSend)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %w", err)
	}

	message := fmt.Sprintf("%s\n%s\n", path, string(data))
	_, err = uc.conn.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("error writing to connection: %w", err)
	}

	reader := bufio.NewReader(uc.conn)
	return uc.processResponse(reader, expectedResponse, expectedType)
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

func (uc *UnixClient) SendString(path string, data string, timeoutInSeconds int) *Response {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "string")
	fmt.Println(resp, err)
	return errorResp(resp, err)
}

func (uc *UnixClient) SendBool(path string, data string, timeoutInSeconds int) *Response {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "bool")
	return errorResp(resp, err)
}

func (uc *UnixClient) SendNumber(path string, data float64, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "number")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

func (uc *UnixClient) SendMap(path string, data map[string]interface{}, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "map")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

func (uc *UnixClient) SendArray(path string, data []interface{}, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "array")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}
