package helper

import (
	"bytes"
	"crypto/rand"
	pkg "djiroutine-go-clean-architecture/pkg"
	"djiroutine-go-clean-architecture/pkg/errors"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	mathRand "math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	DateTimeFormatDefault = "2006-01-02 15:04:05"
	DateFormatDefault     = "2006-01-02"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *mathRand.Rand = mathRand.New(
	mathRand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

func ConsoleLog(title string, data interface{}, content string) {
	fmt.Println()
	log.Println("=== " + title + " ===")
	log.Println(data)
	log.Println(content)
	log.Println("=== " + title + " ===")
	fmt.Println()
}

func RemoveTimeString(s string) string {
	if i := strings.Index(s, "T"); i != -1 {
		return s[:i]
	}

	if i := strings.Index(s, " "); i != -1 {
		return s[:i]
	}

	return s
}

func HandleError(message string, err interface{}) {
	log.Println()
	log.Println("========== Start Error Message ==========")
	log.Println("Message => " + message + ".")

	if err != nil {
		log.Println("Error => ", err)
	}

	log.Println("========== End Of Error Message ==========")
	log.Println()
}

func HandleSuccess(message string) {
	log.Println()
	log.Println("========== Start Message ==========")
	log.Println("Message => " + message + ".")
	log.Println("========== End Of Message ==========")
	log.Println()
}

func MapToString(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "\"%s\":\"%s\"\n", key, value)
	}
	return b.String()
}

func MapToString2(m map[string]interface{}) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "\"%s\":\"%s\"\n", key, value)
	}
	return b.String()
}

func ChangePhoneCode(s string) (newS string) {
	length := len([]rune(s))

	if s[0:1] == "0" {
		newS = "62" + s[1:length]
	}

	return
}

func ChangePhoneCodePlus(s string) (newS string) {
	length := len([]rune(s))

	if s[0:1] == "0" {
		newS = "+62" + s[1:length]
	} else if s[0:1] == "6" {
		newS = "+6" + s[1:length]
	} else {
		newS = s
	}

	return
}

func MarshalToUnmarshal(s string) (detail interface{}) {
	var dByte []byte = []byte(s)

	log.Println(string(dByte))

	u, err := strconv.Unquote(string(dByte))

	log.Println(err)

	err = json.Unmarshal([]byte(u), &detail)
	log.Println(err)

	return
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case errors.ErrInternalServerError:
		return http.StatusInternalServerError
	case errors.ErrForbidden:
		return http.StatusForbidden
	case errors.ErrNotFound:
		return http.StatusNotFound
	case errors.ErrUnAuthorize:
		return http.StatusUnauthorized
	case errors.ErrConflict:
		return http.StatusConflict
	case errors.ErrBadParamInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func JsonDecode(c echo.Context, request interface{}) (interface{}, error) {
	dec := json.NewDecoder(c.Request().Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func JsonString(object interface{}) string {
	res, _ := json.Marshal(object)

	return string(res)
}

func QueryParamDecode(c echo.Context, request interface{}) (interface{}, error) {
	params := c.QueryParams()
	if err := mapQueryParams(params, request); err != nil {
		return nil, err
	}
	return request, nil
}

func mapQueryParams(params url.Values, request interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true) // Ignore unknown keys to prevent errors
	return decoder.Decode(request, params)
}

func GlobalValidationQueryParams(request pkg.GlobalValidation) (bool, string) {
	checkQueryparams, message := ValidationQueryParamsRequired(request.RequiredValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsValueAble(request.ValueAbleValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsDataTypeDate(request.DataTypeNumberDateValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsDataTypeDateMonth(request.DataTypeNumberDateMonthValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsDataTypeNumberInt(request.DataTypeNumberIntValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsDataTypeNumberFloat(request.DataTypeNumberFloatValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryParamsPageLimit(request.PageLimitValidation)
	if checkQueryparams == false {
		return false, message
	}

	checkQueryparams, message = ValidationQueryMaxMinNumber(request.MaxMinNumberValidation)
	if checkQueryparams == false {
		return false, message
	}

	return true, ""
}

func ValidationQueryParamsValueAble(request []pkg.ValueAbleValidation) (bool, string) {

	for _, v := range request {
		found := InArray(v.Value, v.AvailableValue)
		if found == false {
			return false, v.Key + " " + errors.ErrInvalidValue.Error()
		}
	}
	return true, ""
}

func ValidationQueryParamsDataTypeNumberFloat(request []pkg.DataTypeNumberFloatValidation) (bool, string) {

	for _, v := range request {
		if _, err := strconv.ParseFloat(v.Value, 64); err != nil {
			return false, v.Key + " " + errors.ErrInvalidDataType.Error()
		}
	}
	return true, ""
}

func ValidationQueryParamsDataTypeNumberInt(request []pkg.DataTypeNumberIntValidation) (bool, string) {

	for _, v := range request {
		if _, err := strconv.Atoi(v.Value); err != nil {
			return false, v.Key + " " + errors.ErrInvalidDataType.Error()
		}
	}
	return true, ""
}

func ValidationQueryParamsDataTypeDate(request []pkg.DataTypeNumberDateValidation) (bool, string) {
	for _, v := range request {
		if reflect.TypeOf(v.Value).String() == "string" && v.Value != "" {
			valueString := reflect.ValueOf(v.Value).String()
			convertDate := StringToDate(valueString, DateFormatDefault)
			if convertDate.IsZero() == true {
				return false, v.Key + " " + errors.ErrInvalidDataType.Error()
			}
		}
	}
	return true, ""
}

func ValidationQueryParamsDataTypeDateMonth(request []pkg.DataTypeNumberDateMonthValidation) (bool, string) {
	for _, v := range request {
		if reflect.TypeOf(v.Value).String() == "string" && v.Value != "" {
			valueString := reflect.ValueOf(v.Value).String()
			convertDate := StringToDateWithFormat(valueString, "2006-01")
			if convertDate.IsZero() == true {
				return false, v.Key + " " + errors.ErrInvalidDataType.Error()
			}
		}
	}
	return true, ""
}

func ValidationQueryParamsRequired(request []pkg.RequiredValidation) (bool, string) {

	for _, v := range request {
		if v.Value == "" {
			return false, v.Key + " " + errors.ErrIsRequired.Error()
		}
	}
	return true, ""
}

func ValidationQueryParamsMaxMinLonglat(request []pkg.MaxMinLonglatValidation) (bool, string) {

	for _, v := range request {
		if v.Key == "latitude_now" {
			lat := StringToFloat(v.Value)
			if lat < -90 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}

			if lat > 90 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}
		}

		if v.Key == "longitude_now" {
			lat := StringToFloat(v.Value)
			if lat < -180 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}

			if lat > 180 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}
		}
	}

	return true, ""
}

func ValidationQueryMaxMinNumber(validation []pkg.MaxMinNumberValidation) (bool, string) {

	for _, v := range validation {
		val := StringToFloat(v.Value)
		if v.ValueMinNumber == -1 {
			continue
		}
		if val < v.ValueMinNumber {
			return false, v.Key + " " + errors.ErrInvalidValue.Error() + " min " + FloatToString(v.ValueMinNumber)
		}

		if val > v.ValueMaxNumber {
			return false, v.Key + " " + errors.ErrInvalidValue.Error() + " max " + FloatToString(v.ValueMaxNumber)
		}
	}

	return true, ""
}

func ValidationQueryParamsPageLimit(request []pkg.PageLimitValidation) (bool, string) {

	for _, v := range request {
		if v.Key == "page" {
			p := StringToInt(v.Value)
			if p < 1 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}
		}

		if v.Key == "limit" {
			l := StringToInt(v.Value)
			if l < 1 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}

			if l > 200 {
				return false, v.Key + " " + errors.ErrInvalidValue.Error()
			}
		}
	}

	return true, ""
}

func InArray(str string, list []string) bool {
	str = strings.ToLower(str)
	for _, v := range list {
		if strings.ToLower(v) == str {
			return true
		}
	}
	return false
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	if input_num != 0 {
		return strconv.FormatFloat(input_num, 'f', 0, 64)
	} else {
		return "0"
	}
}

func StringNullableToFloat(value *string) float64 {
	if value != nil {
		res, _ := strconv.ParseFloat(*value, 64)
		return res
	}
	return 0
}

func StringToFloat(value string) float64 {
	if value != "" {
		res, _ := strconv.ParseFloat(value, 64)
		return res
	}
	return 0
}

func FloatNUllableToString(input_num *float64) string {
	// to convert a float number to a string
	if input_num != nil {
		return strconv.FormatFloat(*input_num, 'f', 0, 64)
	} else {
		return ""
	}
}

func FloatNUllableToFloat(value *float64) float64 {
	if value != nil {
		return *value
	}
	return 0
}

func FloatToFloatNullable(value float64) *float64 {
	return &value
}

func DateTimeToDateTimeNullable(value time.Time) *time.Time {
	return &value
}

func DateTimeNullableToDateTime(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}

func IntToIntNullable(value int) *int {
	return &value
}

func IntNullableToInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}

func StringToStringNullable(value string) *string {
	return &value
}

func ObjectToString(value interface{}) string {
	result, _ := json.Marshal(value)
	return string(result)
}

func StringNullableToString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}

func IntNullableToStringNullable(value *int) *string {

	if value != nil {
		result := strconv.Itoa(*value)
		return &result
	}
	return nil
}

func IntNullableToString(value *int) string {

	if value != nil {
		result := strconv.Itoa(*value)
		return result
	}
	return "0"
}

func Int64NullableToString(value *int64) string {

	if value != nil {
		result := strconv.FormatInt(*value, 32)
		return result
	}
	return "0"
}

func Int64ToString(value int64) string {

	if value != 0 {
		result := strconv.FormatInt(value, 10)
		return result
	}
	return "0"
}

func IntToString(value int) string {

	if value != 0 {
		result := strconv.Itoa(value)
		return result
	}
	return "0"
}

func StringToIntNullable(value string) *int {

	if value != "" {
		result, _ := strconv.Atoi(value)
		return &result
	}
	return nil
}

func Int64NullableToInt(value *int64) int {

	if value != nil {
		result := int(*value)
		return result
	}
	return 0
}

func StringToInt(value string) int {

	if value != "" {
		result, _ := strconv.Atoi(value)
		return result
	}
	return 0
}

func StringNullableToInt(value *string) int {

	if value != nil {
		result, _ := strconv.Atoi(*value)
		return result
	}
	return 0
}

func StringNullableToDateTimeNullable(value *string) *time.Time {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "2006-01-02 15:04:05"
		date, _ = time.Parse(layoutFormat, *value)
		return &date
	}

	return nil
}

func DateTimeNullableToStringNullable(value *time.Time) *string {
	if value != nil {
		layoutFormat := "2006-01-02 15:04:05"
		date := value.Format(layoutFormat)
		return &date
	}

	return nil
}

func DateTimeToStringNullable(value time.Time) *string {
	layoutFormat := "2006-01-02 15:04:05"
	date := value.Format(layoutFormat)
	return &date
}

func DateTimeToStringWithFormat(value time.Time, format string) string {
	if !value.IsZero() {
		layoutFormat := format
		date := value.Format(layoutFormat)
		return date
	}

	return ""
}

func DateTimeNullableToStringNullableWithFormat(value *time.Time, format string) *string {
	if value != nil {
		layoutFormat := format
		date := value.Format(layoutFormat)
		return &date
	}

	return nil
}

func StringNullableToStringDefaultFormatDate(value *string) *string {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "2006-01-02T15:04:05Z"
		date, _ = time.Parse(layoutFormat, *value)
		dateString := date.Format(DateTimeFormatDefault)
		return &dateString
	}

	return nil
}

func StringNullableToDateTime(value *string) time.Time {
	if value != nil {
		var layoutFormat string
		var date time.Time
		layoutFormat = "2006-01-02T15:04:05Z"
		date, err := time.Parse(layoutFormat, *value)
		if err != nil {
			return time.Time{}
		}
		return date
	}

	return time.Time{}
}

func StringToDateTimeNullable(value string) *time.Time {
	if value != "" {
		var layoutFormat string
		var date time.Time
		layoutFormat = "2006-01-02T15:04:05.999999999Z07:00"
		date, err := time.Parse(layoutFormat, value)
		if err != nil {
			return &time.Time{}
		}
		return &date
	}

	return &time.Time{}
}

func StringToDateWithFormat(value string, format string) time.Time {
	if value != "" {
		var layoutFormat string
		var date time.Time

		layoutFormat = format
		date, _ = time.Parse(layoutFormat, value)
		return date
	}

	return time.Time{}
}

func StringToDate(value string, layout string) time.Time {
	if value != "" {
		var date time.Time

		date, _ = time.Parse(layout, value)
		return date
	}

	return time.Time{}
}

func StringNullableToDateNullable(value *string) *string {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "20060102"
		date, _ = time.Parse(layoutFormat, *value)
		dateString := date.Format("20060102")
		return &dateString
	}

	return nil
}

func ConvertIntBool(value *int) bool {
	if value != nil {
		if *value == 1 {
			return true
		}
	}
	return false
}

func RandomString(length int) string {
	bytes := make([]byte, length)

	for i := 0; i < length; i++ {
		bytes[i] = byte(RandomInt(65, 90))
	}

	return string(bytes)
}

func RandomInt(min int, max int) int {
	return min + mathRand.Intn(max-min)
}

func JSONEncode(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}

func JSONDecode(c echo.Context, request interface{}) (interface{}, error) {
	dec := json.NewDecoder(c.Request().Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ErrId() string {
	errId := time.Now().Format(DateTimeFormatDefault) + " " + RandomString(20)

	return errId
}

func Pagination(qpage, qperPage string) (limit, page, offset int) {
	limit = 20
	page = 1
	offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	if page == 0 && limit == 0 {
		page = 1
		limit = 10
	}
	offset = (page - 1) * limit

	return
}

func CheckInArray(val interface{}, arrays interface{}) bool {
	kind := reflect.TypeOf(arrays).Kind()
	values := reflect.ValueOf(arrays)

	if kind == reflect.Slice || values.Len() > 0 {
		for i := 0; i < values.Len(); i++ {
			if fmt.Sprint(val) == fmt.Sprint(values.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}

func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func SetHttpCookie(c echo.Context, name, value string, age int) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   age,
	})
}
