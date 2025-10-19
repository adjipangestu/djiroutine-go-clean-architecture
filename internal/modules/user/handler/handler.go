package handler

import (
	"context"
	"djiroutine-go-clean-architecture/internal/entity"
	"djiroutine-go-clean-architecture/internal/modules/user"
	"djiroutine-go-clean-architecture/pkg"
	"djiroutine-go-clean-architecture/pkg/errors"
	"djiroutine-go-clean-architecture/pkg/helper"
	"djiroutine-go-clean-architecture/pkg/logger"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Log         logger.Logger
	UserUsecase user.UseCase
}

func NewUserHandler(log logger.Logger, userUseCase user.UseCase) *UserHandler {
	return &UserHandler{
		Log:         log,
		UserUsecase: userUseCase,
	}
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	log := "auth.handler.MasterHandler.ListUser: %s"

	response := new(pkg.ResponseWithPaginator)
	request := new(entity.RequestList)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	dec, err := helper.QueryParamDecode(c, request)
	if err != nil {
		response.MappingResponseError(helper.GetStatusCode(errors.ErrBadParamInput), err.Error())

		return c.JSON(response.Code, response)
	}
	_, _, offset := helper.Pagination(helper.IntToString(*request.Page), helper.IntToString(*request.Limit))
	request = dec.(*entity.RequestList)

	if request.Limit == nil {
		request = request.MappingDefaultPage()
	}

	validation := request.MappingToGlobalValidation()
	checkQueryparams, message := helper.GlobalValidationQueryParams(validation)

	if !checkQueryparams {
		response.MappingResponseError(helper.GetStatusCode(errors.ErrBadParamInput), message)
		return c.JSON(response.Code, response)
	}

	request.Offset = helper.IntToIntNullable(offset)

	res, total, err := h.UserUsecase.ListUsers(ctx, request)
	if err != nil {
		h.Log.Error("["+helper.ErrId()+"]  "+log, err.Error())
		response.MappingResponseError(helper.GetStatusCode(errors.ErrBadParamInput), err.Error())

		return c.JSON(response.Code, response)
	}

	h.Log.Info(log, helper.JsonString(res))

	response.MappingResponseSuccess("Get users list successfull", res)
	response.MappingPagination(int32(*request.Page), int32(*request.Limit), int(total), len(res), response.Response)

	return c.JSON(response.Code, response)
}
