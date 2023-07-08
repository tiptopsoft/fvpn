package relay

import (
	"github.com/topcloudz/fvpn/pkg/handler"
	"github.com/topcloudz/fvpn/pkg/node"
	"github.com/topcloudz/fvpn/pkg/packet"
	"github.com/topcloudz/fvpn/pkg/security"
	"github.com/topcloudz/fvpn/pkg/util"
	"net"
	"runtime"
	"sync"

	"github.com/topcloudz/fvpn/pkg/log"
)

var (
	logger = log.Log()
)

// RegServer use as server
type RegServer struct {
	*util.ServerConfig
	conn         *net.UDPConn
	cache        node.CacheFunc
	ws           sync.WaitGroup
	readHandler  handler.Handler
	writeHandler handler.Handler
	queue        struct {
		outBound *node.OutBoundQueue
		inBound  *node.InBoundQueue
	}

	key struct {
		privateKey security.NoisePrivateKey
		pubKey     security.NoisePublicKey
	}

	//every node has it's own key
	appIds map[string]string
}

func (r *RegServer) Start(address string) error {
	var err error
	r.queue.outBound = node.NewOutBoundQueue()
	r.queue.inBound = node.NewInBoundQueue()
	if r.key.privateKey, err = security.NewPrivateKey(); err != nil {
		return err
	}
	if err = r.start(address); err != nil {
		return err
	}

	r.readHandler = handler.WithMiddlewares(r.serverUdpHandler(), node.Decode())
	r.writeHandler = handler.WithMiddlewares(r.writeUdpHandler(), node.Encode())
	r.cache = node.NewCache()
	r.ws.Wait()
	return nil
}

// Peer register cache for net, and for user create client
func (r *RegServer) start(address string) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.IPv4zero, Port: 4000})
	if err != nil {
		return err
	}
	r.conn = conn
	logger.Debugf("server start at: %s", address)

	nums := runtime.NumCPU()
	for i := 0; i < nums/2; i++ {
		r.ws.Add(1)
		go r.RoutineInBound(i + 1)
		go r.RoutineOutBound(i + 1)
	}

	go r.ReadFromUdp()
	return nil
}

func (r *RegServer) PutPktToOutbound(frame *packet.Frame) {
	r.queue.outBound.PutPktToOutbound(frame)
}

//func (r *RegServer) GetPktFromOutbound() *packet.Frame {
//	return r.queue.outBound.GetPktFromOutbound()
//}

func (r *RegServer) PutPktToInbound(frame *packet.Frame) {
	r.queue.inBound.PutPktToInbound(frame)
}

func (r *RegServer) RoutineInBound(id int) {
	defer r.ws.Done()
	logger.Debugf("start routine %d to handle incomming udp packets", id)
	for {
		select {
		case pkt := <-r.queue.inBound.GetPktFromInbound():
			r.handleInPackets(pkt, id)
		default:

		}

	}
}

func (r *RegServer) handleInPackets(pkt *packet.Frame, id int) {
	pkt.Lock()
	defer func() {
		logger.Debugf("handing in packet success in %d routine finished", id)
		defer pkt.Unlock()
	}()

	err := r.readHandler.Handle(pkt.Context(), pkt)
	if err != nil {
		logger.Error(err)
	}
}

func (r *RegServer) RoutineOutBound(id int) {
	logger.Debugf("start route %d to handle outgoing udp packets", id)
	for {
		select {
		case pkt := <-r.queue.outBound.GetPktFromOutbound():
			r.handleOutPackets(pkt, id)
		default:

		}
	}
}

func (r *RegServer) handleOutPackets(pkt *packet.Frame, id int) {
	pkt.Lock()
	defer func() {
		logger.Debugf("handing out packet success in %d routine finished", id)
		defer pkt.Unlock()
	}()

	err := r.writeHandler.Handle(pkt.Context(), pkt)
	if err != nil {
		logger.Error(err)
	}
}
