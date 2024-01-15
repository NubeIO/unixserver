package unixserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strings"
)

type UnixHandlerFunc func(conn net.Conn, model interface{})

type RouteInfo struct {
	Type             string `json:"type"`
	Path             string `json:"path"`
	ExpectedResponse string `json:"expectedResponse"`
	ExpectsPayload   bool   `json:"expectsPayload"`
}

type Route struct {
	routePath      string
	model          interface{}
	handler        UnixHandlerFunc
	expectsPayload bool
}

type UnixService struct {
	socketPath string
	listener   net.Listener
	routes     map[string]Route
	routesInfo map[string]RouteInfo
}

// NewUnixService creates a new UnixService instance
func NewUnixService(socketPath string) *UnixService {
	return &UnixService{
		socketPath: socketPath,
		routes:     make(map[string]Route),
		routesInfo: make(map[string]RouteInfo),
	}
}

func (us *UnixService) NewUnixRoute(routePath string, model interface{}, handler UnixHandlerFunc, expectsPayload bool, expectedResponse interface{}) {
	us.routes[routePath] = Route{
		routePath:      routePath,
		model:          model,
		handler:        handler,
		expectsPayload: expectsPayload,
	}

	responseType := "nil"
	if expectedResponse != nil {
		responseType = reflect.TypeOf(expectedResponse).String()
	}

	var routeType = "Get"
	if expectsPayload {
		routeType = "Send"
	}
	us.routesInfo[routePath] = RouteInfo{
		Type:             routeType, // Send or "Get" based on your implementation
		Path:             routePath,
		ExpectedResponse: responseType,
		ExpectsPayload:   expectsPayload,
	}

	fmt.Println("Registered route:", routePath)
}

func (us *UnixService) Start() {
	os.Remove(us.socketPath)

	listener, err := net.Listen("unix", us.socketPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	us.listener = listener
	defer us.listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Move defer conn.Close() here to keep the connection open
		go func() {
			defer conn.Close()
			us.handleConnection(conn)
		}()
	}
}

func (us *UnixService) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := us.readMessage(reader)
		if err != nil {
			us.handleReadError(err)
			break
		}

		if handled := us.handleRoute(conn, msg, reader); !handled {
			break
		}
	}
}

func (us *UnixService) readMessage(reader *bufio.Reader) (string, error) {
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(msg), nil
}

func (us *UnixService) handleReadError(err error) {
	if err == io.EOF {
		fmt.Println("Client closed the connection")
	} else {
		fmt.Printf("Error reading message: %v\n", err)
	}
}

func (us *UnixService) handleRoute(conn net.Conn, msg string, reader *bufio.Reader) bool {
	if msg == "" || msg == "null" {
		fmt.Printf("Received empty or null message: '%s'\n", msg)
		return true // continue handling other messages
	}

	fmt.Printf("Received message: '%s'\n", msg)

	if routeInfo, exists := us.routes[msg]; exists {
		us.processRoute(conn, routeInfo, reader)
		return true
	}

	fmt.Printf("No matching route found for message: '%s'\n", msg)
	us.Response(conn, nil, fmt.Errorf("no matching route found"))
	return false
}

func (us *UnixService) processRoute(conn net.Conn, routeInfo Route, reader *bufio.Reader) {
	var modelInstance interface{}
	var err error

	if routeInfo.expectsPayload {
		var payload string
		payload, err = us.readPayload(reader, routeInfo)
		if err != nil {
			us.Response(conn, nil, fmt.Errorf("error reading payload: %v", err))
			return
		}

		modelInstance = us.unmarshalPayload(payload, routeInfo)
		if err != nil {
			us.Response(conn, nil, fmt.Errorf("error unmarshalling payload: %v", err))
			return
		}
	}

	routeInfo.handler(conn, modelInstance)
}

func (us *UnixService) readPayload(reader *bufio.Reader, routeInfo Route) (string, error) {
	payload, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading payload for route %s: %v\n", routeInfo, err)
		return "", err
	}
	return strings.TrimSpace(payload), nil
}

func (us *UnixService) unmarshalPayload(payload string, routeInfo Route) interface{} {
	// Check if model is nil or a pointer type
	if routeInfo.model == nil || reflect.TypeOf(routeInfo.model).Kind() != reflect.Ptr {
		return nil // No model to unmarshal into, or model is not a pointer
	}

	modelInstance := reflect.New(reflect.TypeOf(routeInfo.model).Elem()).Interface()
	err := json.Unmarshal([]byte(payload), modelInstance)
	if err != nil {
		fmt.Printf("Error unmarshalling payload for route '%s': %v\n", routeInfo.routePath, err)
		return nil
	}
	return modelInstance
}

type Response struct {
	Status string // Success Error
	Data   interface{}
	Error  interface{}
}

func (us *UnixService) Response(conn net.Conn, data interface{}, err error) {
	resp := us.prepareResponse(data, err)
	us.writeResponse(conn, resp)
}

func (us *UnixService) prepareResponse(data interface{}, err error) Response {
	if err != nil {
		return Response{Status: "Error", Error: err.Error()}
	}
	return Response{Status: "Success", Data: data}
}

func (us *UnixService) writeResponse(conn net.Conn, resp Response) {
	jsonData, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		fmt.Println("Error marshalling response:", jsonErr)
		return
	}

	response := string(jsonData) + "\n"
	if _, writeErr := conn.Write([]byte(response)); writeErr != nil {
		fmt.Println("Error writing response:", writeErr)
	}
}

func (us *UnixService) RoutesHandler(conn net.Conn, model interface{}) {
	us.Response(conn, us.routesInfo, nil)
}

func (us *UnixService) PingHandler(conn net.Conn, model interface{}) {
	us.Response(conn, "pong", nil)
}
