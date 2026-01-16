package handler

type Response struct {
	Err  string `json:"error"`
	Data any    `json:"data"`
}

func sendResponse(err string, data any) Response {
	return Response{err, data}
}

func (a *AppHandler) SendResponse(data any) Response {
	return sendResponse("", data)
}

func (a *AppHandler) SendError(err error) Response {
	return sendResponse(err.Error(), nil)
}
