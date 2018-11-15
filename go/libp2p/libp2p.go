package libp2p

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type Host struct {
	host host.Host
}

func (h *Host) PeerInfo() *PeerInfo {
	return &PeerInfo{
		pinfo: &pstore.PeerInfo{
			ID:    h.host.ID(),
			Addrs: h.host.Addrs(),
		},
	}
}

type PeerInfo struct {
	pinfo *pstore.PeerInfo
}

func (pi *PeerInfo) ID() *PeerID {
	return &PeerID{pi.pinfo.ID}
}

func ParseMultiaddrString(a string) (*PeerInfo, error) {
	addr, err := ma.NewMultiaddr(a)
	if err != nil {
		return nil, err
	}

	pi, err := pstore.InfoFromP2pAddr(addr)
	if err != nil {
		return nil, err
	}

	return &PeerInfo{pi}, nil
}

func (l *Host) Connect(pinfo *PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := l.host.Connect(ctx, *pinfo.pinfo); err != nil {
		return err
	}

	fmt.Println("Connected to: ", pinfo.pinfo.ID)

	return nil
}

// more of the same, need to wrap most of the interesting types...
// this one is because of the 'loggable' method on peer.ID
type PeerID struct {
	pid peer.ID
}

func (pid PeerID) String() string {
	return pid.pid.Pretty()
}

// have to wrap 'Stream' as the gomobile binder can't handle the 'Conn' type on
// the 'Conn' method of streams...
type Stream struct {
	s net.Stream
}

// currently, gomobile doesnt allow using byte arrays as out parameters
// so, for now. We have to do this dumb thing
func (s *Stream) ReadData(max int) ([]byte, error) {
	b := make([]byte, max)
	n, err := io.ReadFull(s.s, b)
	if err != nil {
		return nil, err
	}
	return b[:n], err
}

func (s *Stream) Write(b []byte) (int, error) {
	return s.s.Write(b)
}

func (s *Stream) Close() error {
	return s.s.Close()
}

func (s *Stream) Reset() error {
	return s.s.Reset()
}

func (l *Host) NewStream(pid *PeerID, proto string) (*Stream, error) {
	s, err := l.host.NewStream(context.TODO(), pid.pid, protocol.ID(proto))
	if err != nil {
		return nil, err
	}

	return &Stream{s}, nil
}

func New(t Transport) (*Host, error) {
	tptopt := libp2p.Transport(&transportConverter{t})
	h, err := libp2p.New(context.TODO(), tptopt)
	if err != nil {
		return nil, err
	}

	return &Host{host: h}, nil
}

type Multiaddr struct {
	addr ma.Multiaddr

	net  string
	host string
	port uint16
}

func (m *Multiaddr) GetHost() (string, error) {
	if m.net == "" {
		if err := m.parse(); err != nil {
			return "", err
		}
	}

	return m.host, nil
}

func (m *Multiaddr) GetPort() (int, error) {
	if m.net == "" {
		if err := m.parse(); err != nil {
			return 0, err
		}
	}

	return int(m.port), nil
}

func (m *Multiaddr) parse() error {
	net, host, err := manet.DialArgs(m.addr)
	if err != nil {
		return err
	}

	parts := strings.Split(host, ":")
	if len(parts) != 2 {
		return fmt.Errorf("expected two parts to host from manet")
	}

	port, err := strconv.ParseUint(parts[1], 10, 16)
	if err != nil {
		return err
	}

	m.net = net
	m.host = host
	m.port = uint16(port)
	return nil
}

type Conn struct {
}

type Listener struct {
}

type Transport interface {
	Dial(raddr *Multiaddr, p *PeerID) (*Conn, error)

	CanDial(addr *Multiaddr) bool

	Listen(laddr *Multiaddr) (*Listener, error)

	//ForEachProtocol(func(int))
	//Protocols() []int

	Proxy() bool
}

type transportConverter struct {
	t Transport
}

func (tc *transportConverter) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.Conn, error) {
	c, err := tc.t.Dial(&Multiaddr{addr: raddr}, &PeerID{p})
	if err != nil {
		return nil, err
	}

	_ = c
	fmt.Println("Need to convert the fake Conn's back up!")
	return nil, nil
}

func (tc *transportConverter) CanDial(addr ma.Multiaddr) bool {
	return tc.t.CanDial(&Multiaddr{addr: addr})
}

func (tc *transportConverter) Listen(laddr ma.Multiaddr) (transport.Listener, error) {
	list, err := tc.t.Listen(&Multiaddr{addr: laddr})
	if err != nil {
		return nil, err
	}

	_ = list
	panic("figure out how to 'unwrap' this")
}

func (tc *transportConverter) Proxy() bool {
	return tc.t.Proxy()
}

func (tc *transportConverter) Protocols() []int {
	var out []int
	/* TODO: figure out how to get an array of things back...
	tc.t.ForEachProtocol(func(i int) {
		out = append(out, i)
	})
	*/
	return out
}

var _ transport.Transport = (*transportConverter)(nil)

var _ = transport.AcceptTimeout
