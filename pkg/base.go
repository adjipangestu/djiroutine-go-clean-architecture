package pkg

import (
	"math"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithPaginator struct {
	Response
	Paginator interface{} `json:"paginator"`
}

type Paginator struct {
	CurrentPage  int32 `json:"current_page"`
	PerPage      int32 `json:"limit_per_page"`
	PreviousPage int32 `json:"back_page"`
	NextPage     int32 `json:"next_page"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
}

func Pagination(qpage, qperPage int32) (limit, page, offset int32) {
	page = qpage
	limit = qperPage
	if page == 0 && limit == 0 {
		page = 1
		limit = 10
	}
	offset = (page - 1) * limit

	return limit, page, offset
}

func (p Paginator) MappingPaginator(page, limit int32, totalAllRecords, countData int) Paginator {
	totalPage := int(math.Ceil(float64(totalAllRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if int(page) != totalPage {
		next = page + 1
	}

	p = Paginator{
		CurrentPage:  page,
		PerPage:      int32(countData),
		PreviousPage: prev,
		NextPage:     next,
		TotalRecords: int32(totalAllRecords),
		TotalPages:   int32(totalPage),
	}

	return p
}

func (r *Response) MappingResponseSuccess(message string, data interface{}) {
	r.Code = http.StatusOK
	r.Status = "success"
	r.Message = message
	r.Data = data
}

func (r *ResponseWithPaginator) MappingPagination(page, limit int32, totalAllRecords, countData int, response Response) {
	paginator := new(Paginator)

	r.Response = response
	r.Paginator = paginator.MappingPaginator(page, limit, totalAllRecords, countData)
}

func (r *Response) MappingResponseError(code int, message string) {
	r.Code = code
	r.Status = "error"
	r.Message = message
	r.Data = nil
}

type GlobalValidation struct {
	RequiredValidation                []RequiredValidation                `json:"required_validation"`
	ValueAbleValidation               []ValueAbleValidation               `json:"value_able_validation"`
	DataTypeNumberDateMonthValidation []DataTypeNumberDateMonthValidation `json:"data_type_number_date_month_validation"`
	DataTypeNumberDateValidation      []DataTypeNumberDateValidation      `json:"data_type_number_date_validation"`
	DataTypeNumberIntValidation       []DataTypeNumberIntValidation       `json:"data_type_number_int_validation"`
	DataTypeNumberFloatValidation     []DataTypeNumberFloatValidation     `json:"data_type_number_float_validation"`
	MaxMinLonglatValidation           []MaxMinLonglatValidation           `json:"max_min_longlat_validation"`
	PageLimitValidation               []PageLimitValidation               `json:"page_limit_validation"`
	MaxMinNumberValidation            []MaxMinNumberValidation            `json:"max_min_number_validation"`
}

type ValueAbleValidation struct {
	Key            string   `json:"key"`
	Value          string   `json:"value"`
	AvailableValue []string `json:"available_value"`
}

type DataTypeNumberDateMonthValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequiredValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DataTypeNumberDateValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DataTypeNumberIntValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DataTypeNumberFloatValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MaxMinLonglatValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PageLimitValidation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MaxMinNumberValidation struct {
	Key            string  `json:"key"`
	Value          string  `json:"value"`
	ValueMaxNumber float64 `json:"value_max_number"`
	ValueMinNumber float64 `json:"value_min_number"`
}
