package controllers

import (
	"github.com/labstack/echo/v4"
)

const (
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type FailureResponse struct {
	Success bool   `json:"success" example:"false"`
	Failure string `json:"failure"`
}

// SuccessResponse Response - Application response success struct
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Result  interface{} `json:"result"`
}

func WriteSuccess(c echo.Context, statusCode int, result interface{}) error {
	return c.JSON(statusCode, SuccessResponse{
		Success: true,
		Result:  result,
	})
}

func WriteError(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, FailureResponse{
		Success: false,
		Failure: err.Error(),
	})
}

func WriteErrorMsg(c echo.Context, statusCode int, errMsg string) error {
	return c.JSON(statusCode, FailureResponse{
		Success: false,
		Failure: errMsg,
	})
}
