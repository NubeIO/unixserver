package main

import (
	"fmt"
	"github.com/NubeIO/unixserver"
	"github.com/NubeIO/unixserver/demo/user"
	"time"
)

func main() {
	service := unixserver.NewUnixService("/tmp/unix.sock")

	userService := user.NewUserService(service)

	service.NewUnixRoute("server/routes", nil, service.RoutesHandler, false, nil)
	service.NewUnixRoute("server/ping", nil, service.PingHandler, false, nil)

	service.NewUnixRoute("user/add", &user.User{}, userService.UserAddHandler, true, &user.User{})
	service.NewUnixRoute("user/get", nil, userService.UserGetHandler, false, &user.User{})

	var dateResp time.Time
	service.NewUnixRoute("user/date", nil, userService.UserGetDateHandler, false, &dateResp)

	fmt.Println("Starting Unix service...")
	service.Start()
}
