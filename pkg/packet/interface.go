package packet

type Interface interface {
	Encode() ([]byte, error)
	Decode([]byte) (Interface, error)
}
