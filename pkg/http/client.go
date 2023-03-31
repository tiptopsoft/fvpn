package http

// Client fvpn client to get token from user center, then use token and group to transfer data.
type Client interface {
	Login(username, password string) (*Response, error)
	Logout(username string) (*Response, error)
}

type Request struct {
}

type Response struct {
	Code    int
	Data    interface{}
	Message string
}

type StarClient struct {
}

func (sc StarClient) Login(username, password string) (*Response, error) {
	//TODO implement
	return nil, nil
}

func (sc StarClient) Logout(username, password string) (*Response, error) {
	//TODO implemenet
	return nil, nil
}
