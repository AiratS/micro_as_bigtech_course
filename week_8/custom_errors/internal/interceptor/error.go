package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys"
	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys/codes"
	"github.com/AiratS/micro_as_bigtech_course/platform_common/pkg/sys/validate"
	grpcCodes "google.golang.org/grpc/codes"
)

type GRPCStatusInterface interface {
	GRPCStatus() *status.Status
}

func ErrorCodesInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := handler(ctx, req)
	if err == nil {
		return res, err
	}

	switch {
	case sys.IsCommonError(err):
		commErr := sys.GetCommonError(err)
		code := toGRPCCode(commErr.Code())

		err = status.Error(code, commErr.Error())

	case validate.IsValidationError(err):
		err = status.Error(grpcCodes.InvalidArgument, err.Error())

	default:
		var se GRPCStatusInterface
		if errors.As(err, &se) {
			return nil, se.GRPCStatus().Err()
		} else {
			if errors.Is(err, context.DeadlineExceeded) {
				err = status.Error(grpcCodes.DeadlineExceeded, err.Error())
			} else if errors.Is(err, context.Canceled) {
				err = status.Error(grpcCodes.Canceled, err.Error())
			} else {
				err = status.Error(grpcCodes.Internal, "internal error")
			}
		}
	}

	return res, err
}

func toGRPCCode(code codes.Code) grpcCodes.Code {
	var res grpcCodes.Code

	switch code {
	case codes.OK:
		res = grpcCodes.OK
	case codes.Canceled:
		res = grpcCodes.Canceled
	case codes.InvalidArgument:
		res = grpcCodes.InvalidArgument
	case codes.DeadlineExceeded:
		res = grpcCodes.DeadlineExceeded
	case codes.NotFound:
		res = grpcCodes.NotFound
	case codes.AlreadyExists:
		res = grpcCodes.AlreadyExists
	case codes.PermissionDenied:
		res = grpcCodes.PermissionDenied
	case codes.ResourceExhausted:
		res = grpcCodes.ResourceExhausted
	case codes.FailedPrecondition:
		res = grpcCodes.FailedPrecondition
	case codes.Aborted:
		res = grpcCodes.Aborted
	case codes.OutOfRange:
		res = grpcCodes.OutOfRange
	case codes.Unimplemented:
		res = grpcCodes.Unimplemented
	case codes.Internal:
		res = grpcCodes.Internal
	case codes.Unavailable:
		res = grpcCodes.Unavailable
	case codes.DataLoss:
		res = grpcCodes.DataLoss
	case codes.Unauthenticated:
		res = grpcCodes.Unauthenticated
	default:
		res = grpcCodes.Unknown
	}

	return res
}
