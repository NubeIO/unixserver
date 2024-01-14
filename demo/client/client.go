package main

import (
	"fmt"
	unixclient "github.com/NubeIO/unixserver/unixcleint"
	"time"
)

type User struct {
	Name string `json:"name"`
}

type Person struct {
	Name2 string `json:"name"`
}

func main() {
	client, err := unixclient.NewUnixClient("/tmp/unix.sock")
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	// No defer client.Close() here

	// Example of sending a user
	userToSend := User{Name: "Alice C"}
	p := &Person{}
	resp, err := client.Send("user/add", &userToSend, 5, p)
	if err != nil {
		fmt.Println("Error sending user:", err)
		return
	}

	Print(resp)

	fmt.Println("Got back from send:", p.Name2)

	// Example of getting a user

	resp, err = client.Get("user/get", 1, p)
	if err != nil {
		fmt.Println("Error getting user:", err)
		return
	}
	Print(resp)

	var pong string
	resp, err = client.Send("server/ping", nil, 5, &pong)
	if err != nil {
		fmt.Println("ping err", err)
		// handle error
	}

	fmt.Println("PING", pong)

	var timeBack time.Time
	resp, err = client.Send("user/date", nil, 5, &timeBack)
	if err != nil {
		fmt.Println("ping err", err)
		// handle error
	}
	Print(resp)
	fmt.Println("timeBack", timeBack)

	client.Close()
}

func Print(i interface{}) {
	fmt.Printf("%+v\n", i)
	return
}
