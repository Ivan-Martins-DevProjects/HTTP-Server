package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Status      int
	ContentType string
	Body        any
}

func SendResponse(status int, contentType string, body any, w io.Writer) {
	response := &Response{
		Status:      status,
		ContentType: contentType,
		Body:        body,
	}

	response.WriteResponse(w)
}

func (r *Response) WriteResponse(w io.Writer) {
	data, _ := json.Marshal(r.Body)

	fmt.Fprintf(w, "HTTP/1.1 %d %s\r\n", r.Status, http.StatusText(r.Status))
	fmt.Fprintf(w, "Content-Type: %s\r\n", r.ContentType)
	fmt.Fprintf(w, "Content-Length: %d\r\n", len(data))
	fmt.Fprintf(w, "Connection: close\r\n")
	fmt.Fprintf(w, "\r\n")

	w.Write(data)
}
