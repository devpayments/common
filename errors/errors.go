package errors

import (
	"errors"
	"fmt"
	errorsPkg "github.com/pkg/errors"
)

type ErrorCode string

func (e ErrorCode) ToString() string {
	return fmt.Sprintf("%s", e)
}

const (
	REQUEST_NOT_VALID          ErrorCode = "REQUEST_NOT_VALID"
	SQL_INSERT_ERROR           ErrorCode = "SQL_INSERT_ERROR"
	SQL_UPDATE_ERROR           ErrorCode = "SQL_UPDATE_ERROR"
	SQL_FETCH_ERROR            ErrorCode = "SQL_FETCH_ERROR"
	API_URL_PARSING_ERROR      ErrorCode = "API_URL_PARSING_ERROR"
	API_REQUEST_CREATION_ERROR ErrorCode = "API_REQUEST_CREATION_ERROR"
	API_REQUEST_ERROR          ErrorCode = "API_REQUEST_ERROR"
	API_REQUEST_STATUS_ERROR   ErrorCode = "API_REQUEST_STATUS_ERROR"
	JSON_SERIALIZATION_ERROR   ErrorCode = "JSON_SERIALIZATION_ERROR"
	JSON_DESERIALIZATION_ERROR ErrorCode = "JSON_DESERIALIZATION_ERROR"
	FORM_SERIALIZATION_ERROR   ErrorCode = "FORM_SERIALIZATION_ERROR"
)

var errorInfo map[ErrorCode]string = map[ErrorCode]string{}

type CustomError struct {
	errorCode     ErrorCode
	errorInfo     string
	errorMsg      string
	error         error
	loggingParams map[string]interface{}
}

func NewCustomError(errorCode ErrorCode, error string) CustomError {
	c := CustomError{errorCode: errorCode, errorInfo: errorInfo[errorCode], errorMsg: error}
	e := errors.New(fmt.Sprintf("Code: %s | %s", c.errorCode, c.errorMsg))
	c.error = errorsPkg.WithStack(e)
	c.loggingParams = make(map[string]interface{}, 0)
	return c
}

func (c CustomError) Error() string {
	return c.error.Error()
}

func (c CustomError) ErrorCode() ErrorCode {
	return c.errorCode
}

func (c CustomError) WithParam(key string, val interface{}) CustomError {
	if c.loggingParams == nil {
		c.loggingParams = make(map[string]interface{}, 0)
	}
	c.loggingParams[key] = val
	return c
}

func PanicIfNecessary(err error) {
	if err != nil {
		panic(err)
	}
}
