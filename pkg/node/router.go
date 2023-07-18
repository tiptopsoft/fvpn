package node

type RouterManager interface {
	AddRouter(cidr string) error
	RemoveRouter(cidr string) error
}

type router struct {
	cidr string
	name string
}

func NewRouter(cidr, name string) RouterManager {
	return &router{
		cidr: cidr,
		name: name,
	}
}
