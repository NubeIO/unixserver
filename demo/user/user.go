package user

import (
	"fmt"
	"github.com/NubeIO/unixserver"
	"net"
	"time"
)

// User struct represents a user with a Name
type User struct {
	Name string `json:"name"`
}

type UserService struct {
	us *unixserver.UnixService
}

func NewUserService(us *unixserver.UnixService) *UserService {
	return &UserService{us: us}
}

// UserAddHandler handles adding a new user
func (us *UserService) UserAddHandler(conn net.Conn, model interface{}) {
	user, ok := model.(*User)
	if !ok {
		us.us.Response(conn, nil, fmt.Errorf("invalid model type"))
		return
	}
	fmt.Println("ADDED USER", user.Name)

	us.us.Response(conn, user, nil)
}

type Person struct {
	Name string `json:"name"`
}

func (us *UserService) UserSendNumber(conn net.Conn, model interface{}) {
	dataIn, ok := model.(*float64)
	if !ok {
		us.us.Response(conn, nil, fmt.Errorf("invalid data type, expected number"))
		return
	}

	responseData := fmt.Sprintf("Echo: %f", *dataIn) // Format as float
	us.us.Response(conn, responseData, nil)
}

func (us *UserService) UserSendString(conn net.Conn, model interface{}) {
	dataIn, ok := model.(*string)
	if !ok {
		us.us.Response(conn, nil, fmt.Errorf("invalid data type, expected string"))
		return
	}

	responseData := fmt.Sprintf("Echo: %f", *dataIn) // Format as float
	us.us.Response(conn, responseData, nil)
}

func (us *UserService) UserGetHandler(conn net.Conn, model interface{}) {
	user := User{Name: "Aidan P"}
	us.us.Response(conn, user, nil)
}

func (us *UserService) UserGetDateHandler(conn net.Conn, model interface{}) {

	us.us.Response(conn, time.Now(), nil)
}

func (us *UserService) UserSendMap(conn net.Conn, model interface{}) {
	dataIn, ok := model.(*map[string]interface{})
	if !ok {
		us.us.Response(conn, nil, fmt.Errorf("invalid data type, expected map"))
		return
	}
	us.us.Response(conn, dataIn, nil)
}

func (us *UserService) UserGetMap(conn net.Conn, model interface{}) {
	// Example data to be sent back
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	us.us.Response(conn, data, nil)
}

func (us *UserService) UserGetArray(conn net.Conn, model interface{}) {
	// Example data to be sent back
	data := []interface{}{"value1", 123, true}
	us.us.Response(conn, data, nil)
}
