package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHitGrpc(t *testing.T) {
	assert := assert.New(t)

	t.Run("normal", func(t *testing.T) {
		ctx := context.Background()
		res, err := HitGrpc(ctx, func(ctx context.Context) error {
			stat := GetRpcStat(ctx)
			stat.RecvBytes = 123
			stat.SentBytes = 456
			return nil
		})

		assert.Nil(err)
		assert.Equal("", res.Error)
		assert.Equal(uint16(200), res.Code)
		assert.Equal(uint64(123), res.RecvBytes)
		assert.Equal(uint64(456), res.SentBytes)
	})

	t.Run("gRPC error", func(t *testing.T) {
		ctx := context.Background()
		res, err := HitGrpc(ctx, func(ctx context.Context) error {
			stat := GetRpcStat(ctx)
			stat.RecvBytes = 123
			stat.SentBytes = 456

			return status.Error(codes.InvalidArgument, "wao")
		})

		assert.EqualError(err, "rpc error: code = InvalidArgument desc = wao")
		assert.Equal("rpc error: code = InvalidArgument desc = wao", res.Error)
		assert.Equal(uint16(400), res.Code)
		assert.Equal(uint64(123), res.RecvBytes)
		assert.Equal(uint64(456), res.SentBytes)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		res, err := HitGrpc(ctx, func(ctx context.Context) error {
			stat := GetRpcStat(ctx)
			stat.RecvBytes = 123
			stat.SentBytes = 456

			return errors.New("wao")
		})

		assert.EqualError(err, "wao")
		assert.Equal("wao", res.Error)
		assert.Equal(uint16(500), res.Code)
		assert.Equal(uint64(123), res.RecvBytes)
		assert.Equal(uint64(456), res.SentBytes)
	})
}
