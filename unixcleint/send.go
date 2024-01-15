package unixclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"
)

func (uc *UnixClient) Send(path string, model interface{}, timeoutInSeconds int, expectedResponse interface{}, expectedType string) (*Response, error) {
	deadline := time.Now().Add(time.Duration(timeoutInSeconds) * time.Second)
	uc.conn.SetDeadline(deadline)

	data, err := json.Marshal(model)
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

func (uc *UnixClient) SendString(path string, data string, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "string")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
}

func (uc *UnixClient) SendBool(path string, data bool, timeoutInSeconds int) (*Response, error) {
	resp, err := uc.Send(path, &data, timeoutInSeconds, nil, "bool")
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	return resp, nil
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
