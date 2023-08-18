package tun

import (
	"errors"
	"fmt"
	"github.com/tiptopsoft/fvpn/pkg/tun/winipcfg"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wintun"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	_ "unsafe"
)

func New() (Device, error) {
	return CreateTun(name, requestGUID, 0)
}

func (tun *adapter) Name() string {
	return tun.name
}

func (tun *adapter) SetIP(net1, ip string) error {
	var err error
	var ipp netip.Prefix
	tun.ip = net.ParseIP(ip)
	link := winipcfg.LUID(tun.wt.LUID())
	if !strings.Contains(ip, "/") {
		ip = fmt.Sprintf("%s/24", ip)
		ipp, err = netip.ParsePrefix(ip)
	}
	if err != nil {
		return err
	}
	err = link.SetIPAddresses([]netip.Prefix{ipp})

	return err
}

func (tun *adapter) SetMTU(mtu int) error {
	tun.MTU = mtu
	return nil
}

func (tun *adapter) IPToString() string {
	return tun.ip.String()
}

func (tun *adapter) Addr() net.IP {
	return tun.ip
}

const (
	name    = "fvpn0"
	tunType = "fvpn0"
)

var (
	_           Device = (*adapter)(nil)
	requestGUID *windows.GUID
)

type adapter struct {
	wt        *wintun.Adapter
	ip        net.IP
	name      string
	handle    windows.Handle
	session   wintun.Session
	readWait  windows.Handle
	running   sync.WaitGroup
	closeOnce sync.Once
	close     atomic.Bool
	MTU       int
}

// CreateTun creates a Wintun interface with the given name and
// a requested GUID. Should a Wintun interface with the same name exist, it is reused.
func CreateTun(ifname string, requestedGUID *windows.GUID, mtu int) (Device, error) {
	wt, err := wintun.CreateAdapter(name, tunType, requestedGUID)
	if err != nil {
		return nil, fmt.Errorf("Error creating interface: %w", err)
	}

	forcedMTU := 1420
	if mtu > 0 {
		forcedMTU = mtu
	}

	tun := &adapter{
		wt:     wt,
		name:   ifname,
		handle: windows.InvalidHandle,
		MTU:    forcedMTU,
	}

	tun.session, err = wt.StartSession(0x800000) // Ring capacity, 8 MiB
	if err != nil {
		tun.wt.Close()
		return nil, fmt.Errorf("Error starting session: %w", err)
	}
	tun.readWait = tun.session.ReadWaitEvent()
	return tun, nil
}

// Note: Read() and Write() assume the caller comes only from a single thread; there's no locking.

func (tun *adapter) Read(buff []byte) (int, error) {
	if tun.close.Load() {
		return 0, os.ErrClosed
	}
	for {
		if tun.close.Load() {
			return 0, os.ErrClosed
		}
		packet, err := tun.session.ReceivePacket()
		switch err {
		case nil:
			size := len(packet)
			copy(buff[:], packet)
			tun.session.ReleaseReceivePacket(packet)
			return size, nil
		case windows.ERROR_NO_MORE_ITEMS:
			windows.WaitForSingleObject(tun.readWait, windows.INFINITE)
			continue
		case windows.ERROR_HANDLE_EOF:
			return 0, os.ErrClosed
		case windows.ERROR_INVALID_DATA:
			return 0, errors.New("send ring corrupt")
		}
		return 0, fmt.Errorf("read failed: %w", err)
	}
}

func (tun *adapter) Write(buff []byte) (int, error) {
	tun.running.Add(1)
	defer tun.running.Done()
	if tun.close.Load() {
		return 0, os.ErrClosed
	}

	packetSize := len(buff)

	_, err := tun.session.AllocateSendPacket(packetSize)
	switch err {
	case nil:
		// TODO: Explore options to eliminate this copy.
		tun.session.SendPacket(buff)
		return len(buff), nil
	case windows.ERROR_HANDLE_EOF:
		return 0, os.ErrClosed
	case windows.ERROR_BUFFER_OVERFLOW:
	default:
		return 0, fmt.Errorf("write failed: %w", err)
	}
	return len(buff), nil
}

//go:linkname nanotime runtime.nanotime
func nanotime() int64
