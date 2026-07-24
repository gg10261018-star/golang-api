package points

import (
	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
)

const (
	statusSuccess = "SUCCESS"
	statusFailed  = "FAILED"
)

// DeductPoints 扣點。requestID 是冪等鍵：同一組 requestID 不論重送幾次，
// 只有第一次會真的扣款，之後的重送都直接回傳第一次的處理結果（Duplicated=true）。
func DeductPoints(gtx *gin.Context, email, requestID string, amount int) (*ChangeResponse, error) {
	if amount <= 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "amount must be positive")
	}
	return applyPointsChange(gtx, email, requestID, -amount)
}

// AddPoints 加點，邏輯與 DeductPoints 共用同一套冪等機制，主要提供測試用來準備餘額。
func AddPoints(gtx *gin.Context, email, requestID string, amount int) (*ChangeResponse, error) {
	if amount <= 0 {
		return nil, errorx.New(errorx.CodeBadRequest, "amount must be positive")
	}
	return applyPointsChange(gtx, email, requestID, amount)
}

func GetBalance(gtx *gin.Context, email string) (int, error) {
	up, err := findPointsByEmail(gtx, email)
	if err != nil {
		return 0, errorx.New(errorx.CodeInternalServerError, err)
	}
	if up == nil {
		return 0, nil
	}
	return up.Balance, nil
}

func applyPointsChange(gtx *gin.Context, email, requestID string, delta int) (*ChangeResponse, error) {
	// 先查一次，同一個 requestID 已經處理過就直接回放結果，不用再走一次交易。
	if existing, err := findTransactionByRequestID(gtx, requestID); err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	} else if existing != nil {
		return replay(existing)
	}

	session := getSession(gtx)
	if err := session.Begin(); err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	up, err := lockPointsByEmail(gtx, email)
	if err != nil {
		session.Rollback()
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}
	if up == nil {
		if err := createUserPoints(gtx, email, 0); err != nil {
			session.Rollback()
			return nil, errorx.New(errorx.CodeInternalServerError, err)
		}
		up = &UserPoints{Email: email, Balance: 0}
	}

	newBalance := up.Balance + delta
	if newBalance < 0 {
		record := &PointTransaction{
			Email: email, RequestID: requestID, Amount: delta,
			Status: statusFailed, BalanceAfter: up.Balance, Message: "insufficient balance",
		}
		if err := insertTransaction(gtx, record); err != nil {
			session.Rollback()
			if isDuplicateRequestID(err) {
				return replayByRequestID(gtx, requestID)
			}
			return nil, errorx.New(errorx.CodeInternalServerError, err)
		}
		if err := session.Commit(); err != nil {
			return nil, errorx.New(errorx.CodeInternalServerError, err)
		}
		return nil, errorx.New(errorx.CodeConflict, "insufficient balance")
	}

	if err := updateBalance(gtx, email, newBalance); err != nil {
		session.Rollback()
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	record := &PointTransaction{
		Email: email, RequestID: requestID, Amount: delta,
		Status: statusSuccess, BalanceAfter: newBalance,
	}
	if err := insertTransaction(gtx, record); err != nil {
		session.Rollback()
		if isDuplicateRequestID(err) {
			return replayByRequestID(gtx, requestID)
		}
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	if err := session.Commit(); err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	return &ChangeResponse{Balance: newBalance, Duplicated: false}, nil
}

func replayByRequestID(gtx *gin.Context, requestID string) (*ChangeResponse, error) {
	existing, err := findTransactionByRequestID(gtx, requestID)
	if err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}
	if existing == nil {
		return nil, errorx.New(errorx.CodeInternalServerError, "transaction disappeared after duplicate key conflict")
	}
	return replay(existing)
}

func replay(t *PointTransaction) (*ChangeResponse, error) {
	if t.Status == statusFailed {
		return nil, errorx.New(errorx.CodeConflict, t.Message)
	}
	return &ChangeResponse{Balance: t.BalanceAfter, Duplicated: true}, nil
}
