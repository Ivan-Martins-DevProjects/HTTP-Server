package main

import (
	"bytes"
	// "fmt"
	"io"
	"log"
	"net"

	"github.com/Ivan-Martins-DevProjects/HTTP-Server/cmd/internal"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			}

			str += string(data)
		}

		if len(str) != 0 {
			out <- str
		}
	}()

	return out
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

			ipAddr, _, _ := net.SplitHostPort(c.RemoteAddr().String())

			// Verifica rate limit
			if err := IpList.AddAndCheckIP(ipAddr); err != nil {
				log.Printf("IP %s bloqueado: %v", ipAddr, err)
				return
			} else {
				log.Printf("IP %s logado", ipAddr)
			}

			// Processa as linhas
			// for line := range getLinesChannel(c) {
			// 	fmt.Printf("[%s] Read: %s\n", ipAddr, line)
			// }
		}(conn)
	}

}
