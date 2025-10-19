package user

import (
	"context"
	"djiroutine-go-clean-architecture/internal/entity"
)

type Repository interface {
	ListUsers(ctx context.Context, param *entity.RequestList) ([]*entity.UserResponse, error)
	GetTotalUsers(ctx context.Context, param *entity.RequestList) (int64, error)
}
