package internal

import (
	"fmt"
	"log"
	"net"
)

func HandleConnections(c net.Conn, IpList *IpList) {
	defer c.Close()

	req, err := GetRequestContext(c)
	if err != nil {
		log.Fatal("error", "error", err)
		return
	}

	ipAddr := c.RemoteAddr().String()
	if err := IpList.AddAndCheckIP(ipAddr); err != nil {
		log.Fatal("error", "error", err)
		return
	}

	log.Printf("IP %s logado\n", ipAddr)

	// Imprime os headers capturados
	params := req.Params
	fmt.Println("PARAMS")
	fmt.Printf("Method: %s, Path: %s, Version: %s\n\n", params.Method, params.Path, params.Version)

	fmt.Println("HEADERS")
	for _, header := range req.Headers {
		fmt.Printf("%s: %s\n", header.GetKey(), header.GetValue())
	}

	fmt.Printf("\nBODY\n")
	for _, body := range req.Body {
		fmt.Printf("%s: %s\n", body.GetKey(), body.GetValue())
	}

}

func Run(listener net.Listener) error {
	ipList := CreateIpList()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go HandleConnections(conn, ipList)
	}

}
