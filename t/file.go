package main

import "os"

func main() {
	f := "/dev/net/tun"

	file, err := os.OpenFile(f, os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}

	s := Socket{FileDescriptor: int(file.Fd())}

	b := make([]byte, 2048)

	if _, err := s.Read(b); err != nil {
		panic(err)
	}
}
