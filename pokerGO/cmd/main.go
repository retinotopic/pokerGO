package main

import (
	server "pokerGO/internal/app"
)

func main() {
	srv := server.NewServer("localhost:8080")
	srv.Run()
}
