package usercase

import (
	"context"
	"djiroutine-go-clean-architecture/internal/entity"
	"djiroutine-go-clean-architecture/internal/modules/user"
	"djiroutine-go-clean-architecture/pkg/logger"
	"time"
)

type UserUsecase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
	log            logger.Logger
}

func NewUserUsecase(userRepo user.Repository, timeout time.Duration, log logger.Logger) user.UseCase {
	return &UserUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
		log:            log,
	}
}

func (u UserUsecase) ListUsers(ctx context.Context, request *entity.RequestList) (res []*entity.UserResponse, total int64, err error) {
	log := "modules.master.usecase.ListPegawai: %s"

	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.ListUsers(ctx, request)
	if err != nil {
		u.log.Error(log+"list users - ", err.Error())

		return nil, 0, err
	}

	total, err = u.userRepo.GetTotalUsers(ctx, request)
	if err != nil {
		u.log.Error(log+"count total users - ", err.Error())

		return nil, 0, err
	}

	return res, total, err
}
