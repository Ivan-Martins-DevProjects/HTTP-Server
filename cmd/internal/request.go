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
	Headers map[string]any
	Body    map[string]any
	Content bool
	mu      sync.Mutex
}

func (r *Request) AppendHeaderToRequest(header *Header) {
	r.Headers[strings.ToLower(header.key)] = header.value
}

func (r *Request) AppendBodyToRequest(body *Body) {
	r.Body[strings.ToLower(body.key)] = body.value
}

type Params struct {
	Method  string
	Path    string
	Version string
	Params  map[string]string
}

func (p *Params) SetParams() {
	finalParams := make(map[string]string)
	parameters := strings.SplitN(p.Path, "?", 2)
	if len(parameters) > 0 {
		for _, param := range parameters {
			individualParams := strings.SplitN(param, "=", 2)
			if len(individualParams) == 2 {
				key := strings.TrimSpace(individualParams[0])
				value := strings.TrimSpace(individualParams[0])
				finalParams[key] = value
			}
		}
	}

	p.Params = finalParams
}

func GetRequestContext(f io.ReadCloser) (*Request, error) {
	req := &Request{
		Content: false,
		Headers: make(map[string]any),
		Body:    make(map[string]any),
	}
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

	contentType, ok := req.Headers["content-type"]
	if !ok {
		return nil, fmt.Errorf("Content-Type não especificado")
	}
	if str, ok := contentType.(string); ok {
		if strings.EqualFold(str, "application/json") {
			req.Content = true
		}
	} else {
		return nil, fmt.Errorf("Content-Type deve ser uma string")
	}

	var contentLength int64
	if req.Content {
		length, ok := req.Headers["content-length"]
		if !ok {
			return nil, fmt.Errorf("Content-Length não especificado")
		}
		strLength, _ := length.(string)

		convertedLength, err := strconv.Atoi(strLength)
		if err != nil {
			return nil, fmt.Errorf("Content-Lenght deve ser um inteiro")

		}

		contentLength = int64(convertedLength)

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
