package internal

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

type Request struct {
	Params  *Params
	Headers []Content
	Body    []Content
	Content bool
	mu      sync.Mutex
}

func (r *Request) AppendHeaderToRequest(header *Header) {
	r.Headers = append(r.Headers, header)
}

func (r *Request) AppendBodyToRequest(body *Body) {
	r.Body = append(r.Body, body)
}

type Params struct {
	Method  string
	Path    string
	Version string
}

func GetRequestContext(f io.ReadCloser) (*Request, error) {
	req := &Request{Content: false}
	scanner := bufio.NewScanner(f)

	if scanner.Scan() {
		pathLine := scanner.Text()
		parts := strings.Fields(pathLine)
		if len(parts) >= 2 {
			method := &Params{
				Method:  parts[0],
				Path:    parts[1],
				Version: parts[2],
			}
			req.Params = method
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}

		KeyAndValue := strings.SplitN(line, ":", 2)
		if len(KeyAndValue) == 2 {
			key := strings.TrimSpace(KeyAndValue[0])
			value := strings.TrimSpace(KeyAndValue[1])
			header := &Header{
				key:   key,
				value: value,
			}
			req.AppendHeaderToRequest(header)
		}
	}

	for _, header := range req.Headers {
		if strings.EqualFold(header.GetKey(), "Content-Type") && strings.EqualFold(header.GetValue(), "application/json") {
			req.Content = true
		}
	}

	var contentLength int64
	if req.Content {
		for _, header := range req.Headers {
			if strings.EqualFold(header.GetKey(), "Content-Length") {
				number, err := strconv.Atoi(header.GetValue())
				contentLength = int64(number)
				if err != nil {
					return nil, fmt.Errorf("Erro ao extrair Content-Lenght")
				}
			}
		}

		if contentLength > 0 {
			limitedReader := io.LimitReader(f, contentLength)
			scanner := bufio.NewScanner(limitedReader)

			for scanner.Scan() {
				line := scanner.Text()
				body, err := getBodyContent(line)
				if err != nil {
					return nil, fmt.Errorf("Erro ao extrair conteúdo do body")
				}

				if body == nil {
					continue
				}
				req.AppendBodyToRequest(body)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return req, nil
}

func getBodyContent(input string) (*Body, error) {
	KeyAndValue := strings.SplitN(input, ":", 2)
	if len(KeyAndValue) != 2 {
		return nil, nil
	}

	key := strings.TrimSpace(KeyAndValue[0])
	value := strings.TrimSpace(KeyAndValue[1])
	body := &Body{
		key:   key,
		value: value,
	}

	return body, nil
}
