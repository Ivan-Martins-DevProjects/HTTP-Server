package internal

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

type Request struct {
	Headers []*Header
	mu      sync.Mutex
}

func (r *Request) AppendHeaderToRequest(header *Header) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Headers = append(r.Headers, header)
}

type Header struct {
	Key   string
	Value string
}

func CreateHeader(key, value string) *Header {
	return &Header{
		Key:   key,
		Value: value,
	}
}

func GetHeadersRequest(f io.ReadCloser) (*Request, error) {
	req := &Request{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}

		KeyAndValue := strings.SplitN(line, ":", 2)
		if len(KeyAndValue) == 2 {
			key := strings.TrimSpace(KeyAndValue[0])
			value := strings.TrimSpace(KeyAndValue[1])

			header := CreateHeader(key, value)
			req.AppendHeaderToRequest(header)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return req, nil
}
