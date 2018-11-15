package libp2p

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	ma "github.com/multiformats/go-multiaddr"
)

type Libp2p struct {
	host host.Host
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

func (l *Libp2p) Connect(pinfo *PeerInfo) error {
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
	fmt.Println("read data: ", max, n, b[:n])
	return b[:n], err
}

func (s *Stream) Write(b []byte) (int, error) {
	fmt.Println("writing data: ", b)
	return s.s.Write(b)
}

func (s *Stream) Close() error {
	return s.s.Close()
}

func (s *Stream) Reset() error {
	return s.s.Reset()
}

func (l *Libp2p) NewStream(pid *PeerID, proto string) (*Stream, error) {
	s, err := l.host.NewStream(context.TODO(), pid.pid, protocol.ID(proto))
	if err != nil {
		return nil, err
	}

	return &Stream{s}, nil
}

func New() (*Libp2p, error) {
	h, err := libp2p.New(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Libp2p{host: h}, nil
}
