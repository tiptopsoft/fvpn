package packet

type FrameType int

const (
	ARP FrameType = 0x0806

	IP FrameType = 0x0800

	IPV6Type FrameType = 0x86DD
)
