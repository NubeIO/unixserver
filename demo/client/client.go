package main

import (
	"fmt"
	unixclient "github.com/NubeIO/unixserver/unixcleint"
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

	resp, err := client.SendString("user/send/string", "hello", 5)
	if err != nil {
		fmt.Println("Error sending user:", err)
		return
	}
	Print(resp)

	resp, err = client.SendNumber("user/send/number", 1234.66, 5)
	if err != nil {
		fmt.Println("Error sending user:", err)
		return
	}
	Print(resp)

	resultMap := make(map[string]interface{})
	resultMap["a"] = "abc"
	resultMap["num"] = 123.66

	resp, err = client.SendMap("user/send/map", resultMap, 5)
	if err != nil {
		fmt.Println("Error sending user:", err)
		return
	}
	Print(resp)

	resp, err = client.GetString("user/get", 5)
	if err != nil {
		fmt.Println("Error sending user:", err)
		return
	}
	Print(resp)

	resp, err = client.GetMap("user/get/map", 5)
	if err != nil {
		fmt.Println("Error getting map:", err)
		return
	}
	Print(resp)

	resp, err = client.GetArray("user/get/array", 5)
	if err != nil {
		fmt.Println("Error getting map:", err)
		return
	}
	Print(resp)

	client.Close()
}

func Print(i interface{}) {
	fmt.Printf("%+v\n", i)
	return
}
