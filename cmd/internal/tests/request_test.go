package tests

import (
	"net"
	"testing"

	"github.com/Ivan-Martins-DevProjects/HTTP-Server/cmd/internal"
)

func TestGetRequestContext(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go func() {
		clientConn.Write([]byte("POST / HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 15\r\n\r\n{\"id\": \"123\"}"))
		clientConn.Close()
	}()

	req, err := internal.GetRequestContext(serverConn)
	if err != nil {
		t.Fatalf("Erro na função: %v", err)
	}

	if req.Params.Method != "POST" {
		t.Errorf("Esperava POST, recebi %s", req.Params.Method)
	}
}
