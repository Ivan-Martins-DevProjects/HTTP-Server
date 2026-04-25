package internal

import (
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

	SendResponse(
		200,
		"application/json",
		req.Body,
		c,
	)

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
