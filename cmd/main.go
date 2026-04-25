package main

import (
	"github.com/Ivan-Martins-DevProjects/HTTP-Server/cmd/internal"
	"log"
	"net"
)

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("error", "error", err)
	}
	defer listener.Close()

	log.Printf("Servidor rodando em 127.0.0.1:%s\n", port)
	if err := internal.Run(listener); err != nil {
		log.Fatalf("Erro interno: %v", err)
	}

}
