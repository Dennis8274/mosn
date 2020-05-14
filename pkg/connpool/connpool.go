package connpool

import (
	"container/list"
	"context"
	"mosn.io/api"
	"mosn.io/mosn/pkg/types"
	"sync"
	"sync/atomic"
)

// Connpool is a connection pool
type Connpool struct {
	idleClients []*activeClient

	host types.Host
	mux  sync.Mutex

	multiplexMode bool

	totalClientCount uint64 // total clients

	wantConnQueue list.List
}

// NewConnPool creates a new connection pool
func NewConnPool(ctx context.Context, host types.Host) *Connpool {
	return &Connpool{
		host : host,
	}
}

// GetConn get a conn from conn pool
func (p *Connpool) GetConn(ctx context.Context) (*activeClient, types.PoolFailureReason) {
	maxConns := p.host.ClusterInfo().ResourceManager().Connections().Max()

	p.mux.Lock()
	defer p.mux.Unlock()
	idleClientCount := len(p.idleClients)

	// there is no available client
	if idleClientCount == 0 {
		// max conns is 0 means no limit
		if maxConns == 0 || atomic.LoadUint64(&p.totalClientCount) < maxConns {
			// create new conn
			return p.newConnLocked(ctx)
		}

		return nil, types.Overflow
	}

	// conn array len > 0
	var (
		usedConns = atomic.LoadUint64(&p.totalClientCount) - uint64(idleClientCount)
		lastIdx   = idleClientCount - 1
	)

	// Only refuse extra connection, keepalive-connection is closed by timeout
	if maxConns != 0 && usedConns > maxConns {
		return nil, types.Overflow
	}

	c := p.idleClients[lastIdx] // return the last conn
	p.idleClients[lastIdx] = nil

	if !p.multiplexMode {
		p.idleClients = p.idleClients[:lastIdx]
	}

	return c, ""
}

// newConnLocked creates a new conn
// should acquire pool lock before entering
func (p *Connpool) newConnLocked(ctx context.Context) (*activeClient, types.PoolFailureReason) {
	ac := &activeClient{
		pool: p,
		host: p.host.CreateConnection(ctx),
	}

	if err := ac.host.Connection.Connect(); err != nil {
		return nil, types.ConnectionFailure
	}

	atomic.AddUint64(&p.totalClientCount, 1)
	p.idleClients = append(p.idleClients, ac)

	return ac, ""
}

func (p *Connpool) tryPutIdleConn(ac *activeClient) {
	if p.multiplexMode {
		// do nothing
		return
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	if atomic.LoadUint32(&ac.goaway) == 1 {
		return
	}

	p.idleClients = append(p.idleClients, ac)
}

// Shutdown graceful shut down the conn pool
func (p *Connpool) Shutdown() {
	p.mux.Lock()
	defer p.mux.Unlock()

	for _, client := range p.idleClients {
		client.GoAway()
	}
}

type activeClient struct {
	index int // index in connpool idleClients array
	pool  *Connpool
	host  types.CreateConnectionData

	goaway uint32
}

// only for multiplex mode !
func (p *Connpool) clearIdleClients() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.idleClients = p.idleClients[:0]
}

// Close the conn, if there is no error, the conn will be put back to pool
func (ac *activeClient) Close(err error) {
	if err != nil {
		atomic.StoreUint32(&ac.goaway, 1)

		if ac.pool.multiplexMode {
			// in multiplex mode, there is only one client, clear it
			ac.pool.clearIdleClients()
		}

		// TODO... is it proper to close it here
		ac.host.Connection.Close(api.NoFlush, api.LocalClose)

		// To subtract a signed positive constant value c from x, do AddUint64(&x, ^uint64(c-1)).
		atomic.AddUint64(&ac.pool.totalClientCount, uint64(0))

		return
	}

	// put conn back to pool
	ac.pool.tryPutIdleConn(ac)
}

// GoAway go away the client
func (ac *activeClient) GoAway() {
	atomic.StoreUint32(&ac.goaway, 1)
}
