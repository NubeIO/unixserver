package unixclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type UnixClient struct {
	socketPath string
	conn       net.Conn
}

// NewUnixClient creates a new UnixClient instance
func NewUnixClient(socketPath string) (*UnixClient, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Unix socket: %w", err)
	}

	return &UnixClient{
		socketPath: socketPath,
		conn:       conn,
	}, nil
}

func (uc *UnixClient) Reconnect() error {
	if uc.conn != nil {
		uc.conn.Close()
	}
	var err error
	uc.conn, err = net.Dial("unix", uc.socketPath)
	return err
}

func (uc *UnixClient) Close() {
	if uc.conn != nil {
		uc.conn.Close()
	}
}

type Response struct {
	Status string // Success Error
	Data   interface{}
	Error  interface{}
}

func (uc *UnixClient) Send(path string, model interface{}, timeoutInSeconds int, expectedResponse interface{}) (*Response, error) {
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
	return uc.processResponse(reader, expectedResponse)
}

func (uc *UnixClient) Get(path string, timeoutInSeconds int, expectedResponse interface{}) (*Response, error) {
	deadline := time.Now().Add(time.Duration(timeoutInSeconds) * time.Second)
	uc.conn.SetDeadline(deadline)

	if _, err := uc.conn.Write([]byte(path + "\n")); err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	reader := bufio.NewReader(uc.conn)
	return uc.processResponse(reader, expectedResponse)
}

func (uc *UnixClient) processResponse(reader *bufio.Reader, expectedResponse interface{}) (*Response, error) {
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var resp Response
	err = json.Unmarshal([]byte(responseStr), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if resp.Status == "Success" && expectedResponse != nil && resp.Data != nil {
		dataBytes, err := json.Marshal(resp.Data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling response data: %w", err)
		}
		err = json.Unmarshal(dataBytes, expectedResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling response data: %w", err)
		}
	}

	return &resp, nil
}
