package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Response struct {
	Conn   net.Conn
	Status int
	Body   any
}

func CreateResponse(c net.Conn) *Response {
	return &Response{
		Conn:   c,
		Status: 0,
		Body:   nil,
	}
}

func (r *Response) SendResponse(status int, body any) {
	r.Status = status
	r.Body = body

	r.WriteResponse(r.Conn)
}

func (r *Response) WriteResponse(w io.Writer) {
	data, _ := json.Marshal(r.Body)

	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", r.Status, http.StatusText(r.Status))
	fmt.Fprintf(w, "Content-Type: application/json\r\n")
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(data))
	fmt.Fprintf(w, "Connection: close\r\n")
	fmt.Fprintf(w, "\r\n")

	w.Write(data)
}
