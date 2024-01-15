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

type ValidationResponse struct {
	OkMessage    string `json:"okMessage"`
	Code         string `json:"code"`
	Advice       string `json:"advice,omitempty"` // eg; an exiting entry already contains filed ""
	ErrorMessage string `json:"error,omitempty"`
	IsError      bool   `json:"isError"`
}

func main() {
	client, err := unixclient.NewUnixClient("/tmp/unix.sock")
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}
	validationResponse := &ValidationResponse{}
	resp, err := client.Send("validation/ip", "192.168.", 5, &validationResponse, "any")
	Print(resp)
	Print(validationResponse)

	client.Close()
}

func Print(i interface{}) {
	fmt.Printf("%+v\n", i)
	return
}
