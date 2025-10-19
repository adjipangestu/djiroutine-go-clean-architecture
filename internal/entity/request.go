package entity

import (
	"djiroutine-clean-architecture/pkg"
	"djiroutine-clean-architecture/pkg/helper"
)

// request
type RequestList struct {
	Page   *int    `json:"page"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
	Search *string `json:"search"`
}

func (request *RequestList) MappingToGlobalValidation() pkg.GlobalValidation {
	res := pkg.GlobalValidation{
		RequiredValidation: []pkg.RequiredValidation{
			{
				Key:   "Page",
				Value: helper.IntToString(*request.Page),
			},
			{
				Key:   "Limit",
				Value: helper.IntToString(*request.Limit),
			},
		},
		DataTypeNumberIntValidation: []pkg.DataTypeNumberIntValidation{
			{
				Key:   "Page",
				Value: helper.IntToString(*request.Page),
			},
			{
				Key:   "Limit",
				Value: helper.IntToString(*request.Limit),
			}},
		DataTypeNumberFloatValidation: []pkg.DataTypeNumberFloatValidation{},
		PageLimitValidation: []pkg.PageLimitValidation{
			{
				Key:   "Page",
				Value: helper.IntToString(*request.Page),
			},
			{
				Key:   "Limit",
				Value: helper.IntToString(*request.Limit),
			},
		},
		MaxMinNumberValidation: []pkg.MaxMinNumberValidation{},
	}

	return res
}

func (req RequestList) MappingDefaultPage() *RequestList {
	res := &RequestList{
		Limit:  helper.IntToIntNullable(100),
		Page:   helper.IntToIntNullable(1),
		Offset: helper.IntToIntNullable(0),
	}
	return res
}

// param
type ParamList struct {
	Page   *int    `json:"page"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
	Search *string `json:"search"`
}
