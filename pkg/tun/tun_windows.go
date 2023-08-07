package tun

func New() (*Device, error) {

}

func (a *Adapter) Name() string {

}

func (a Adapter) Read() ([]byte, error) {
	a.StartSession()
	return nil, nil
}

var (
	TUNAME   = "fvpn0"
	TUN_TYPE = "fvpn0"
)

type NativeTun struct {
	wt        *wintun.Adapter
	name      string
	handle    windows.Handle
	rate      rateJuggler
	session   wintun.Session
	readWait  windows.Handle
	events    chan Event
	running   sync.WaitGroup
	closeOnce sync.Once
	close     atomic.Bool
	forcedMTU int
	outSizes  []int
}

// CreateTUNWithRequestedGUID creates a Wintun interface with the given name and
// a requested GUID. Should a Wintun interface with the same name exist, it is reused.
func CreateTUNWithRequestedGUID(ifname string, requestedGUID *windows.GUID, mtu int) (Device, error) {
	wt, err := wintun.CreateAdapter(ifname, WintunTunnelType, requestedGUID)
	if err != nil {
		return nil, fmt.Errorf("Error creating interface: %w", err)
	}

	forcedMTU := 1420
	if mtu > 0 {
		forcedMTU = mtu
	}

	tun := &NativeTun{
		wt:        wt,
		name:      ifname,
		handle:    windows.InvalidHandle,
		events:    make(chan Event, 10),
		forcedMTU: forcedMTU,
	}

	tun.session, err = wt.StartSession(0x800000) // Ring capacity, 8 MiB
	if err != nil {
		tun.wt.Close()
		close(tun.events)
		return nil, fmt.Errorf("Error starting session: %w", err)
	}
	tun.readWait = tun.session.ReadWaitEvent()
	return tun, nil
}

// Note: Read() and Write() assume the caller comes only from a single thread; there's no locking.

func (tun *NativeTun) Read(bufs [][]byte, sizes []int, offset int) (int, error) {
	tun.running.Add(1)
	defer tun.running.Done()
retry:
	if tun.close.Load() {
		return 0, os.ErrClosed
	}
	start := nanotime()
	shouldSpin := tun.rate.current.Load() >= spinloopRateThreshold && uint64(start-tun.rate.nextStartTime.Load()) <= rateMeasurementGranularity*2
	for {
		if tun.close.Load() {
			return 0, os.ErrClosed
		}
		packet, err := tun.session.ReceivePacket()
		switch err {
		case nil:
			packetSize := len(packet)
			copy(bufs[0][offset:], packet)
			sizes[0] = packetSize
			tun.session.ReleaseReceivePacket(packet)
			tun.rate.update(uint64(packetSize))
			return 1, nil
		case windows.ERROR_NO_MORE_ITEMS:
			if !shouldSpin || uint64(nanotime()-start) >= spinloopDuration {
				windows.WaitForSingleObject(tun.readWait, windows.INFINITE)
				goto retry
			}
			procyield(1)
			continue
		case windows.ERROR_HANDLE_EOF:
			return 0, os.ErrClosed
		case windows.ERROR_INVALID_DATA:
			return 0, errors.New("Send ring corrupt")
		}
		return 0, fmt.Errorf("Read failed: %w", err)
	}
}

func (tun *NativeTun) Write(bufs [][]byte, offset int) (int, error) {
	tun.running.Add(1)
	defer tun.running.Done()
	if tun.close.Load() {
		return 0, os.ErrClosed
	}

	for i, buf := range bufs {
		packetSize := len(buf) - offset
		tun.rate.update(uint64(packetSize))

		packet, err := tun.session.AllocateSendPacket(packetSize)
		switch err {
		case nil:
			// TODO: Explore options to eliminate this copy.
			copy(packet, buf[offset:])
			tun.session.SendPacket(packet)
			continue
		case windows.ERROR_HANDLE_EOF:
			return i, os.ErrClosed
		case windows.ERROR_BUFFER_OVERFLOW:
			continue // Dropping when ring is full.
		default:
			return i, fmt.Errorf("Write failed: %w", err)
		}
	}
	return len(bufs), nil
}
