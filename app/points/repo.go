package points

import (
	"errors"
	"time"

	"github.com/example/golang-api/pkg/ctv"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func getSession(gtx *gin.Context) *xorm.Session {
	session, _ := ctv.Get[*xorm.Session](gtx, ctv.DBSession)
	return session
}

func findTransactionByRequestID(gtx *gin.Context, requestID string) (*PointTransaction, error) {
	var t PointTransaction
	found, err := getSession(gtx).Where("request_id = ?", requestID).Get(&t)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &t, nil
}

func findPointsByEmail(gtx *gin.Context, email string) (*UserPoints, error) {
	var up UserPoints
	found, err := getSession(gtx).Where("email = ?", email).Get(&up)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &up, nil
}

// lockPointsByEmail 用 SELECT ... FOR UPDATE 鎖住該筆餘額，避免同一使用者的併發請求同時讀到舊餘額。
func lockPointsByEmail(gtx *gin.Context, email string) (*UserPoints, error) {
	var up UserPoints
	found, err := getSession(gtx).Where("email = ?", email).ForUpdate().Get(&up)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &up, nil
}

func createUserPoints(gtx *gin.Context, email string, balance int) error {
	_, err := getSession(gtx).Insert(&UserPoints{
		Email:   email,
		Balance: balance,
		Updated: time.Now().UTC(),
	})
	return err
}

func updateBalance(gtx *gin.Context, email string, balance int) error {
	_, err := getSession(gtx).
		Table("user_points").
		Where("email = ?", email).
		Update(map[string]interface{}{
			"balance": balance,
			"updated": time.Now().UTC(),
		})
	return err
}

func insertTransaction(gtx *gin.Context, t *PointTransaction) error {
	t.Created = time.Now().UTC()
	_, err := getSession(gtx).Insert(t)
	return err
}

// isDuplicateRequestID 判斷是不是因為 uq_request_id 撞到了 —— 代表另一個併發請求用同一組
// request_id 搶先寫入成功，這裡應該回頭讀那筆結果，而不是當成一般錯誤。
func isDuplicateRequestID(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
