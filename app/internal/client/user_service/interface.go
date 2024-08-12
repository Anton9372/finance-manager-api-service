package user_service

import "context"

type UserService interface {
	Create(ctx context.Context, dto SignUpUserDTO) (User, error)
	GetByUUID(ctx context.Context, uuid string) (User, error)
	GetByEmailAndPassword(ctx context.Context, email, password string) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
}
