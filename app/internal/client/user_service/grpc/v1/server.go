package user_service_grpc

import (
	"context"
	"finance-manager-api-service/internal/client/user_service"
	"finance-manager-api-service/pkg/logging"
	"fmt"
	protoUserService "github.com/Anton9372/user-service-contracts/gen/go/user_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

const requestWaitTime = 5 * time.Second

type client struct {
	grpcClient protoUserService.UserServiceClient
	Conn       *grpc.ClientConn
	logger     *logging.Logger
}

func NewClient(grpcServerHostPort string, logger *logging.Logger) (user_service.UserService, error) {
	conn, err := grpc.NewClient(grpcServerHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("can not connect to gRPC server: %v", err)
		return nil, fmt.Errorf("can not connect to gRPC server: %v", err)
	}

	grpcClient := protoUserService.NewUserServiceClient(conn)

	return &client{
		grpcClient: grpcClient,
		Conn:       conn,
		logger:     logger,
	}, nil
}

func (c *client) Create(ctx context.Context, dto user_service.SignUpUserDTO) (user_service.User, error) {
	c.logger.Debug("Create user")
	req := NewCreateUserRequest(dto)

	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()

	resp, err := c.grpcClient.Create(reqCtx, req)
	if err != nil {
		c.logger.Error("failed to create user: %v", err)
		return user_service.User{}, HandleGrpcServerError(err)
	}

	user, err := c.GetByUUID(reqCtx, resp.Uuid)
	if err != nil {
		return user_service.User{}, err
	}

	return user, nil
}

func (c *client) GetByUUID(ctx context.Context, uuid string) (user_service.User, error) {
	c.logger.Debug("Get user by uuid")

	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()

	resp, err := c.grpcClient.GetByUUID(reqCtx, &protoUserService.GetByUUIDRequest{Uuid: uuid})
	if err != nil {
		c.logger.Error("failed to get user by uuid: %v", err)
		return user_service.User{}, HandleGrpcServerError(err)
	}
	return NewUserResponse(resp), nil
}

func (c *client) GetByEmailAndPassword(ctx context.Context, email, password string) (user_service.User, error) {
	c.logger.Debug("Get user by email and password")

	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()

	resp, err := c.grpcClient.GetByEmailAndPassword(reqCtx,
		&protoUserService.GetByEmailAndPasswordRequest{Email: email, Password: password})
	if err != nil {
		c.logger.Error("failed to get user by email and password: %v", err)
		return user_service.User{}, HandleGrpcServerError(err)
	}
	return NewUserResponse(resp), nil
}

func (c *client) Update(ctx context.Context, dto user_service.UpdateUserDTO) error {
	c.logger.Debug("Update user")
	req := NewUpdateUserRequest(dto)

	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()

	_, err := c.grpcClient.Update(reqCtx, req)
	if err != nil {
		c.logger.Error("failed to update user: %v", err)
		return HandleGrpcServerError(err)
	}
	return nil
}

func (c *client) Delete(ctx context.Context, uuid string) error {
	c.logger.Debug("Delete user")

	reqCtx, cancel := context.WithTimeout(ctx, requestWaitTime)
	defer cancel()

	_, err := c.grpcClient.Delete(reqCtx, &protoUserService.DeleteRequest{Uuid: uuid})
	if err != nil {
		c.logger.Error("failed to delete user: %v", err)
		return HandleGrpcServerError(err)
	}
	return nil
}
