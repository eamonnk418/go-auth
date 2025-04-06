package main

import (
	"github.com/eamonnk418/go-auth/internal/auth"
	"github.com/eamonnk418/go-auth/internal/server"
)

func init() {
	auth.NewAuth()
}

func main() {
	server := server.NewServer(nil)

	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
