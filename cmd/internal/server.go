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
	params.SetParams()
	for key, value := range params.Params {
		fmt.Printf("%s: %s\n", key, value)
	}

	fmt.Println("HEADERS")
	for key, value := range req.Headers {
		fmt.Printf("%s: %s\n", key, value)
	}

	fmt.Printf("\nBODY\n")
	for key, value := range req.Body {
		fmt.Printf("%s: %s\n", key, value)
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
