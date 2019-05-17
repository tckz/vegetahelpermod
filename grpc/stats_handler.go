package grpc

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc/stats"
)

type contextKeyRpcStatMarker struct{}

var contextKeyRpcStat = &contextKeyRpcStatMarker{}

type RpcStat struct {
	SentCount uint64
	SentBytes uint64
	RecvCount uint64
	RecvBytes uint64
}

// SetRpcStat sets intstance of RpcStat to context.
func SetRpcStat(ctx context.Context, s *RpcStat) context.Context {
	return context.WithValue(ctx, contextKeyRpcStat, s)
}

// GetRpcStat returns intstance of RpcStat from the context.
// If there is no RpcStat is set, returns nil.
func GetRpcStat(ctx context.Context) *RpcStat {
	if v, ok := ctx.Value(contextKeyRpcStat).(*RpcStat); ok {
		return v
	}
	return nil
}

// RpcStatsHandler treats Sent/Recv Count/Bytes
type RpcStatsHandler struct {
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (r *RpcStatsHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	// no-op
}

// TagConn exists to satisfy gRPC stats.Handler.
func (r *RpcStatsHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (r *RpcStatsHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	switch st := rs.(type) {
	case *stats.OutPayload:
		if s := GetRpcStat(ctx); s != nil {
			atomic.AddUint64(&s.SentBytes, uint64(st.Length))
			atomic.AddUint64(&s.SentCount, 1)
		}
	case *stats.InPayload:
		if s := GetRpcStat(ctx); s != nil {
			atomic.AddUint64(&s.RecvBytes, uint64(st.Length))
			atomic.AddUint64(&s.RecvCount, 1)
		}
	}
}

// TagRPC implements per-RPC context management.
func (r *RpcStatsHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	return ctx
}
