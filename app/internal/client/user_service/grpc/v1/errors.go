package user_service_grpc

import (
	"finance-manager-api-service/internal/apperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleGrpcServerError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return apperror.APIError("UNKNOWN", "An unknown error occurred", err.Error())
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return apperror.BadRequestError(st.Message())
	case codes.NotFound:
		return apperror.ErrNotFound
	default:
		return err
	}
}
