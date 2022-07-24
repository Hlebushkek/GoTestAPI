package main

import (
	"fmt"
	"github.com/Hlebushkek/GoTestAPI/api"
)

const (
	port = "3000"
)

func main() {
	fmt.Printf("Starting api server on port %s...\n", port)
	err := api.StartServer(port)
	if err != nil {
		panic(err)
	}
}
