package http

type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"data"`
	Message string      `json:"msg"`
}
