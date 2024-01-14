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

// UserGetHandler handles retrieving a user
func (us *UserService) UserGetHandler(conn net.Conn, model interface{}) {
	user := User{Name: "Aidan P"}
	us.us.Response(conn, user, nil)
}

func (us *UserService) UserGetDateHandler(conn net.Conn, model interface{}) {

	us.us.Response(conn, time.Now(), nil)
}
