package model

import (
	"net/http"
	"time"

	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error"`
	Code      int         `json:"code"`
	Timestamp string      `json:"timestamp"`
}

func Success(gtx *gin.Context, data interface{}) {
	gtx.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   "success",
		Data:      data,
		Error:     nil,
		Code:      http.StatusOK,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func Created(gtx *gin.Context, data interface{}) {
	gtx.JSON(http.StatusCreated, Response{
		Success:   true,
		Message:   "success",
		Data:      data,
		Error:     nil,
		Code:      http.StatusCreated,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func Fail(gtx *gin.Context, err error) {
	code := http.StatusInternalServerError
	if e, ok := err.(*errorx.AppError); ok {
		code = e.Code
	}
	gtx.JSON(code, Response{
		Success:   false,
		Message:   err.Error(),
		Data:      nil,
		Error:     err.Error(),
		Code:      code,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
