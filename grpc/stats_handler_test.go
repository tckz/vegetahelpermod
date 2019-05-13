package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/stats"
)

func TestRpcStatsHandler(t *testing.T) {
	assert := assert.New(t)

	t.Run("With rpcstat", func(t *testing.T) {
		handler := &RpcStatsHandler{}
		ctx := context.Background()
		stat := &RpcStat{}

		ctx = SetRpcStat(ctx, stat)

		// nothing happen
		handler.HandleConn(ctx, &stats.ConnBegin{})
		handler.TagConn(ctx, &stats.ConnTagInfo{})
		handler.TagRPC(ctx, &stats.RPCTagInfo{})

		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 11,
		})
		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 12,
		})
		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 13,
		})
		handler.HandleRPC(ctx, &stats.OutPayload{
			Length: 21,
		})
		handler.HandleRPC(ctx, &stats.OutPayload{
			Length: 22,
		})

		assert.Equal(uint64(43), stat.SentBytes)
		assert.Equal(uint64(2), stat.SentCount)
		assert.Equal(uint64(36), stat.RecvBytes)
		assert.Equal(uint64(3), stat.RecvCount)
	})

	t.Run("Without rpcstat", func(t *testing.T) {
		handler := &RpcStatsHandler{}
		ctx := context.Background()

		// nothing happen
		handler.HandleConn(ctx, &stats.ConnBegin{})
		handler.TagConn(ctx, &stats.ConnTagInfo{})
		handler.TagRPC(ctx, &stats.RPCTagInfo{})

		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 11,
		})
		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 12,
		})
		handler.HandleRPC(ctx, &stats.InPayload{
			Length: 13,
		})
		handler.HandleRPC(ctx, &stats.OutPayload{
			Length: 21,
		})
		handler.HandleRPC(ctx, &stats.OutPayload{
			Length: 22,
		})

		// nothing happen, don't panic
	})
}

func TestSetGetRpcStat(t *testing.T) {
	assert := assert.New(t)

	t.Run("normal", func(t *testing.T) {
		st := &RpcStat{}
		ctx := context.Background()
		ctx = SetRpcStat(ctx, st)
		ret := GetRpcStat(ctx)
		if ret != st {
			t.Errorf("Set address must be retrieved")
		}
	})

	t.Run("no set", func(t *testing.T) {
		ctx := context.Background()
		ret := GetRpcStat(ctx)
		assert.Nil(ret)
	})
}
