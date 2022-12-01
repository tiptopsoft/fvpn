package syscall

type EpollEvent struct {
	Events uint32
	_      int32
	Fd     int32
	Pad    int32
}
