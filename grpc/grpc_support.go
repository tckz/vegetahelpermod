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

func MapCode2Status(code codes.Code) uint16 {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusInternalServerError
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.ResourceExhausted:
		return http.StatusInternalServerError
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusInternalServerError
	case codes.OutOfRange:
		return http.StatusInternalServerError
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
