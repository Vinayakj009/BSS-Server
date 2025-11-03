package main

import (
	"bss/src/database"
	"bss/src/server"
	"context"
	"fmt"
	"os"
)

func main() {
	db, err := database.NewDb(context.Background())
	if err != nil {
		panic(err)
	}
	server := server.NewServer(db)
	addr := ":" + os.Getenv("APP_PORT")
	fmt.Println("Starting BSS Server... on port", addr)
	if err := server.Start(addr); err != nil {
		panic(err)
	}
}
