package pack

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode(bs []byte, rem, idx int) (interface{}, error)
}

type EnDecoder interface {
	Encode() ([]byte, error)
	Decode(bs []byte, rem, idx int) (interface{}, error)
}
