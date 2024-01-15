package unixclient

import (
	"bufio"
	"fmt"
	"time"
)

func (uc *UnixClient) Get(path string, timeoutInSeconds int, expectedResponse interface{}, expectedType string) (*Response, error) {
	deadline := time.Now().Add(time.Duration(timeoutInSeconds) * time.Second)
	uc.conn.SetDeadline(deadline)

	if _, err := uc.conn.Write([]byte(path + "\n")); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	reader := bufio.NewReader(uc.conn)
	return uc.processResponse(reader, expectedResponse, expectedType)
}

func (uc *UnixClient) GetString(path string, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Get(path, timeoutInSeconds, nil, "string")
	return resp, err
}

func (uc *UnixClient) GetBool(path string, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Get(path, timeoutInSeconds, nil, "bool")
	return resp, err
}

func (uc *UnixClient) GetNumber(path string, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Get(path, timeoutInSeconds, nil, "number")
	return resp, err
}

func (uc *UnixClient) GetMap(path string, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Get(path, timeoutInSeconds, nil, "map")
	if err != nil {
		return nil, fmt.Errorf("error getting map: %w", err)
	}

	return resp, nil
}

func (uc *UnixClient) GetArray(path string, timeoutInSeconds int) (*Response, error) {
	var resp Response
	var data []interface{}
	_, err := uc.Get(path, timeoutInSeconds, &data, "array")
	if err != nil {
		return nil, fmt.Errorf("error getting array: %w", err)
	}
	resp.Data = data
	return &resp, nil
}
