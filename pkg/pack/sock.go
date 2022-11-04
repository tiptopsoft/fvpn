package pack

type StarSock struct {
	Family uint8
	Type   uint8
	Port   uint16
	Addr   [128]byte
}
