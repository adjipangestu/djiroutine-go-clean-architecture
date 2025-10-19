package user

import (
	"context"
	"djiroutine-go-clean-architecture/internal/entity"
)

type UseCase interface {
	ListUsers(ctx context.Context, request *entity.RequestList) (res []*entity.UserResponse, total int64, err error)
}
