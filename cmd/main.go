package main

import (
	"fmt"
	"github.com/Ivan-Martins-DevProjects/HTTP-Server/cmd/internal"
	"log"
	"net"
)

type Result struct {
	Req *internal.Request
	Err error
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	IpList := internal.CreateIpList()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		go func(c net.Conn) {
			defer c.Close()

			// Função para capturar os headers de forma assíncrona
			headersChannel := make(chan Result)
			go func() {
				req, err := internal.GetHeadersRequest(c)
				headersChannel <- Result{Req: req, Err: err}
			}()

			// Coleta o endereço IP
			ipAddr, _, _ := net.SplitHostPort(c.RemoteAddr().String())

			// Verifica rate limit
			if err := IpList.AddAndCheckIP(ipAddr); err != nil {
				log.Fatalf("IP %s bloqueado: %v", ipAddr, err)
			}
			log.Printf("IP %s logado\n", ipAddr)

			// Imprime os headers capturados
			result := <-headersChannel
			if result.Err != nil {
				log.Fatal("error", "error", result.Err)
			}

			params := result.Req.Params
			fmt.Printf("Method: %s, Path: %s, Version: %s\n", params.Method, params.Path, params.Version)
			for _, header := range result.Req.Headers {
				fmt.Printf("%s: %s\n", header.Key, header.Value)
			}
		}(conn)
	}
}
