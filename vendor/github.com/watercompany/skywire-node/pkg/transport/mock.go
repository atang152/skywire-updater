package transport

import (
	"context"
	"io"
	"net"
	"time"

	"errors"

	"github.com/skycoin/skycoin/src/cipher"
)

type fConn struct {
	net.Conn
	cipher.PubKey
}

// MockFactory implements Factory over net.Pipe connections.
type MockFactory struct {
	local    cipher.PubKey
	connChan chan *fConn
}

// NewMockFactory constructs a pair of MockFactories.
func NewMockFactory(local, remote cipher.PubKey) (*MockFactory, *MockFactory) {
	connChan := make(chan *fConn)
	return &MockFactory{local, connChan}, &MockFactory{remote, connChan}
}

// Accept waits for new net.Conn notification from another MockFactory.
func (f *MockFactory) Accept(ctx context.Context) (Transport, error) {
	conn, more := <-f.connChan
	if !more {
		return nil, errors.New("factory: closed")
	}

	return NewMockTransport(conn, f.local, conn.PubKey), nil
}

// Dial creates pair of net.Conn via net.Pipe and passes one end to another MockFactory.
func (f *MockFactory) Dial(ctx context.Context, remote cipher.PubKey) (Transport, error) {
	in, out := net.Pipe()
	f.connChan <- &fConn{in, f.local}
	return NewMockTransport(out, f.local, remote), nil
}

// Close closes notification channel between a pair of MockFactories.
func (f *MockFactory) Close() error {
	select {
	case <-f.connChan:
	default:
		close(f.connChan)
	}
	return nil
}

// Local returns a local PubKey of the Factory.
func (f *MockFactory) Local() cipher.PubKey {
	return f.local
}

// Type returns type of the Factory.
func (f *MockFactory) Type() string {
	return "mock"
}

// MockTransport is a transport that accepts custom writers and readers to use them in Read and Write
// operations
type MockTransport struct {
	rw      io.ReadWriteCloser
	local   cipher.PubKey
	remote  cipher.PubKey
	context context.Context
}

// NewMockTransport creates a transport with the given secret key and remote public key, taking a writer
// and a reader that will be used in the Write and Read operation
func NewMockTransport(rw io.ReadWriteCloser, local, remote cipher.PubKey) *MockTransport {
	return &MockTransport{rw, local, remote, context.Background()}
}

// Read implements reader for mock transport
func (m *MockTransport) Read(p []byte) (n int, err error) {
	select {
	case <-m.context.Done():
		return 0, ErrTransportCommunicationTimeout
	default:

		return m.rw.Read(p)
	}
}

// Write implements writer for mock transport
func (m *MockTransport) Write(p []byte) (n int, err error) {
	select {
	case <-m.context.Done():
		return 0, ErrTransportCommunicationTimeout
	default:
		return m.rw.Write(p)
	}
}

// Close implements closer for mock transport
func (m *MockTransport) Close() error {
	return m.rw.Close()
}

// Local returns the local static public key
func (m *MockTransport) Local() cipher.PubKey {
	return m.local
}

// Remote returns the remote public key fo the mock transport
func (m *MockTransport) Remote() cipher.PubKey {
	return m.remote
}

// SetDeadline sets a deadline for the write/read operations of the mock transport
func (m *MockTransport) SetDeadline(t time.Time) error {
	// nolint
	ctx, cancel := context.WithDeadline(m.context, t)
	m.context = ctx

	go func(cancel context.CancelFunc) {
		time.Sleep(time.Until(t))
		cancel()
	}(cancel)

	return nil
}

// Type returns the type of the mock transport
func (m *MockTransport) Type() string {
	return "mock"
}
