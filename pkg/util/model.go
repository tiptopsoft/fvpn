package util

type Login struct {
	Auth string `json:"auth"`
}

type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"data"`
	Message string      `json:"msg"`
}

type JoinRequest struct {
	IP      string `json:"ip"`
	Network string `json:"network"`
}

type JoinResponse struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	NetworkId string `json:"networkId"`
}

func HttpOK(data interface{}) Response {
	return Response{
		Code:    200,
		Result:  data,
		Message: "",
	}
}

func HttpError(message string) Response {
	return Response{
		Code:    500,
		Result:  nil,
		Message: message,
	}
}
