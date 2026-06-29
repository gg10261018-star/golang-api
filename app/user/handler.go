package user

import (
	"github.com/example/golang-api/internal/model"
	"github.com/example/golang-api/pkg/ctv"
	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	r.POST("/register", handleRegister)
	r.POST("/login", handleLogin)
}

func RouterProtected(r *gin.RouterGroup) {
	r.PUT("/password", handleChangePassword)
}

func handleRegister(gtx *gin.Context) {
	var req RegisterRequest
	if err := gtx.ShouldBindJSON(&req); err != nil {
		model.Fail(gtx, errorx.New(errorx.CodeBadRequest, err.Error()))
		return
	}

	if err := register(gtx, req.Email, req.Password); err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Created(gtx, nil)
}

func handleLogin(gtx *gin.Context) {
	var req LoginRequest
	if err := gtx.ShouldBindJSON(&req); err != nil {
		model.Fail(gtx, errorx.New(errorx.CodeBadRequest, err.Error()))
		return
	}

	token, err := login(gtx, req.Email, req.Password)
	if err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Success(gtx, LoginResponse{Token: token})
}

func handleChangePassword(gtx *gin.Context) {
	var req ChangePasswordRequest
	if err := gtx.ShouldBindJSON(&req); err != nil {
		model.Fail(gtx, errorx.New(errorx.CodeBadRequest, err.Error()))
		return
	}

	email, _ := ctv.Get[string](gtx, ctv.EmailKey)

	if err := changePassword(gtx, email, req.OldPassword, req.NewPassword); err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Success(gtx, nil)
}
