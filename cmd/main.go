package main

import (
	server "github.com/retinotopic/pokerGO/internal/app"
)

func main() {
	srv := server.NewServer("localhost:8080")
	srv.Run()
}
