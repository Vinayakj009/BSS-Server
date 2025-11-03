package main

import (
	"bss/src/server"
	"fmt"
	"os"
)

func main() {
	server := server.NewServer()
	addr := ":" + os.Getenv("APP_PORT")
	fmt.Println("Starting BSS Server... on port", addr)
	if err := server.Start(addr); err != nil {
		panic(err)
	}
}
