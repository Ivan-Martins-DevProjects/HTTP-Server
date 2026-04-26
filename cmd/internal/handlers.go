package internal

func TestResponse(req *Request, response *Response) {
	res := map[string]string{
		"Status":  "Sucesso",
		"Message": "Request realizada",
	}
	response.SendResponse(200, res)
}

func EchoResponse(req *Request, response *Response) {
	res := map[string]any{
		"Método":  req.Method(),
		"Params":  req.Query(),
		"Headers": req.Header(),
		"Body":    req.Data(),
	}

	response.SendResponse(200, res)
}
