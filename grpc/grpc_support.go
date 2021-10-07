package grpc

import (
	"context"
	"net/http"

	"github.com/tckz/vegetahelper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HitGrpcFunc func(context.Context) error

func HitGrpc(ctx context.Context, f HitGrpcFunc) (*vegetahelper.HitResult, error) {
	stat := &RpcStat{}
	ctx = SetRpcStat(ctx, stat)
	err := f(ctx)

	result := &vegetahelper.HitResult{
		SentBytes: stat.SentBytes,
		RecvBytes: stat.RecvBytes,
	}
	if err == nil {
		result.Code = http.StatusOK
	} else {
		s, ok := status.FromError(err)
		if ok {
			result.Code = MapCode2Status(s.Code())
		} else {
			result.Code = http.StatusInternalServerError
		}
		result.Error = err.Error()
	}

	return result, err
}

var waoTable = map[codes.Code]uint16{}

var code2statusTable = map[codes.Code]uint16{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           http.StatusInternalServerError,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusRequestTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.ResourceExhausted:  http.StatusInternalServerError,
	codes.FailedPrecondition: http.StatusPreconditionFailed,
	codes.Aborted:            http.StatusInternalServerError,
	codes.OutOfRange:         http.StatusInternalServerError,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusServiceUnavailable,
	codes.DataLoss:           http.StatusInternalServerError,
	codes.Unauthenticated:    http.StatusUnauthorized,
}

func MapCode2Status(code codes.Code) uint16 {
	if v, ok := code2statusTable[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
