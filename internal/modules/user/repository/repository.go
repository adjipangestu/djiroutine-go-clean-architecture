package repository

import (
	"context"
	"djiroutine-go-clean-architecture/internal/entity"
	"djiroutine-go-clean-architecture/pkg/config"
	"djiroutine-go-clean-architecture/pkg/logger"
)

type UserRepository struct {
	db  config.DBService
	log logger.Logger
}

func NewUserRepository(db config.DBService, log logger.Logger) *UserRepository {
	return &UserRepository{
		db:  db,
		log: log,
	}
}

func (r *UserRepository) ListUsers(ctx context.Context, param *entity.RequestList) ([]*entity.UserResponse, error) {
	log := "modules.master.repository.ListUser: %s"

	var res []*entity.UserResponse
	employee := new([]entity.User)
	query := r.db.GetConnection().WithContext(ctx)

	if param.Search != nil || *param.Search != "" {
		searchPattern := "%" + *param.Search + "%"
		query = query.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	if param.Limit != nil && param.Offset != nil {
		query = query.Limit(*param.Limit).Offset(*param.Offset)
	}

	err := query.Model(&employee).Find(&res).Error
	if err != nil {
		r.log.Error(log, err)

		return nil, err
	}

	return res, nil
}

func (r *UserRepository) GetTotalUsers(ctx context.Context, param *entity.RequestList) (int64, error) {
	var total int64
	log := "modules.master.repository.CountTotalPegawai: %s"

	employee := new([]entity.User)
	query := r.db.GetConnection().WithContext(ctx)

	if param.Search != nil {
		searchPattern := "%" + *param.Search + "%"
		query = query.Where("LOWER(first_name) LIKE LOWER(?) OR LOWER(last_name) LIKE LOWER(?)", searchPattern, searchPattern)
	}

	err := query.Model(&employee).Count(&total).Error
	if err != nil {
		r.log.Error(log, err)

		return 0, err
	}

	return total, nil
}
