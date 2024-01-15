package unixclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
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

func (uc *UnixClient) processResponse(reader *bufio.Reader, expectedResponse interface{}, expectedType string) (*Response, error) {
	responseStr, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	var resp Response
	err = json.Unmarshal([]byte(responseStr), &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	// Handling different expected types
	if resp.Status == "Success" && expectedResponse != nil && resp.Data != nil {
		switch expectedType {
		case "string":
			if strPtr, ok := expectedResponse.(*string); ok {
				*strPtr = resp.Data.(string)
				*strPtr = resp.AsString

			}
		case "number":
			if numPtr, ok := expectedResponse.(*float64); ok {
				*numPtr = resp.Data.(float64)
				*numPtr = resp.AsNumber
			}
		case "bool":
			if boolPtr, ok := expectedResponse.(*bool); ok {
				*boolPtr = resp.Data.(bool)
				*boolPtr = resp.AsBool
			}
		case "map":
			if mapPtr, ok := expectedResponse.(*map[string]interface{}); ok {
				*mapPtr = resp.Data.(map[string]interface{})
			}
		case "array":
			if arrayPtr, ok := expectedResponse.(*[]interface{}); ok {
				*arrayPtr = resp.Data.([]interface{})
			}
		default:
			dataBytes, err := json.Marshal(resp.Data)
			if err != nil {
				return nil, fmt.Errorf("error marshalling response data: %w", err)
			}
			err = json.Unmarshal(dataBytes, expectedResponse)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling response data: %w", err)
			}
		}
	}

	return &resp, nil
}
