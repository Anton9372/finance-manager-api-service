package user_service_grpc

import (
	"finance-manager-api-service/internal/client/user_service"
	protoUserService "github.com/Anton9372/user-service-contracts/gen/go/user_service/v1"
)

func NewCreateUserRequest(dto user_service.SignUpUserDTO) *protoUserService.CreateRequest {
	return &protoUserService.CreateRequest{
		Name:             dto.Name,
		Email:            dto.Email,
		Password:         dto.Password,
		RepeatedPassword: dto.RepeatedPassword,
	}
}

func NewUserResponse(resp *protoUserService.UserResponse) user_service.User {
	return user_service.User{
		UUID:     resp.User.Uuid,
		Name:     resp.User.Name,
		Email:    resp.User.Email,
		Password: resp.User.Password,
	}
}

func NewUpdateUserRequest(dto user_service.UpdateUserDTO) *protoUserService.UpdateRequest {
	return &protoUserService.UpdateRequest{
		Uuid:                dto.UUID,
		Name:                dto.Name,
		Email:               dto.Email,
		Password:            dto.Password,
		NewPassword:         dto.NewPassword,
		RepeatedNewPassword: dto.RepeatedPassword,
	}
}
