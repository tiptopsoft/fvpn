package node

type RouterManager interface {
	AddRouter(ip string) error
	RemoveRouter(ip string) error
}

type router struct {
	ip   string
	name string
}

func NewRouter(ip, name string) RouterManager {
	return &router{
		ip:   ip,
		name: name,
	}
}
