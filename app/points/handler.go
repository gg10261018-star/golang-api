package points

import (
	"github.com/example/golang-api/internal/model"
	"github.com/example/golang-api/pkg/ctv"
	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
)

func RouterProtected(r *gin.RouterGroup) {
	r.GET("/points", handleGetBalance)
	r.POST("/points/topup", handleTopUp)
	r.POST("/points/deduct", handleDeduct)
}

func handleGetBalance(gtx *gin.Context) {
	email, _ := ctv.Get[string](gtx, ctv.EmailKey)

	balance, err := GetBalance(gtx, email)
	if err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Success(gtx, gin.H{"balance": balance})
}

func handleTopUp(gtx *gin.Context) {
	var req ChangeRequest
	if err := gtx.ShouldBindJSON(&req); err != nil {
		model.Fail(gtx, errorx.New(errorx.CodeBadRequest, err.Error()))
		return
	}

	email, _ := ctv.Get[string](gtx, ctv.EmailKey)

	resp, err := AddPoints(gtx, email, req.RequestID, req.Amount)
	if err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Success(gtx, resp)
}

func handleDeduct(gtx *gin.Context) {
	var req ChangeRequest
	if err := gtx.ShouldBindJSON(&req); err != nil {
		model.Fail(gtx, errorx.New(errorx.CodeBadRequest, err.Error()))
		return
	}

	email, _ := ctv.Get[string](gtx, ctv.EmailKey)

	resp, err := DeductPoints(gtx, email, req.RequestID, req.Amount)
	if err != nil {
		model.Fail(gtx, err)
		return
	}

	model.Success(gtx, resp)
}
