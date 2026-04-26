package internal

import (
	"log"
	"net"
)

type Content struct {
	Request *Request
	IpAddr  string
}

func GetContext(c net.Conn, ipList *IpList) *Content {
	req, err := GetRequestContext(c)
	if err != nil {
		log.Printf("Erro: %v", err)
		return nil
	}

	ipAddr := c.RemoteAddr().String()
	if err := ipList.AddAndCheckIP(ipAddr); err != nil {
		log.Printf("Erro: %v", err)
		return nil
	}

	return &Content{
		Request: req,
		IpAddr:  ipAddr,
	}
}

func HandleConnections(c net.Conn) {
	defer c.Close()

	ipList := CreateIpList()
	context := GetContext(c, ipList)

	log.Printf("IP %s logado\n", context.IpAddr)

	router := CreateRouter(c)
	response := CreateResponse(c)

	router.AddRoute("/api", HandlerFunc(TestResponse))
	router.AddRoute("/echo", HandlerFunc(EchoResponse))

	router.Serve(context.Request, response)
}

func Run(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go HandleConnections(conn)
	}

}
