package user

import (
	"time"

	"github.com/example/golang-api/pkg/ctv"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

func getSession(gtx *gin.Context) *xorm.Session {
	session, _ := ctv.Get[*xorm.Session](gtx, ctv.DBSession)
	return session
}

func findByEmail(gtx *gin.Context, email string) (*User, error) {
	var u User
	found, err := getSession(gtx).Where("email = ?", email).Get(&u)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &u, nil
}

func createUser(gtx *gin.Context, u *User) error {
	now := time.Now().UTC()
	u.Created = now
	u.Updated = now
	_, err := getSession(gtx).Insert(u)
	return err
}

func updatePassword(gtx *gin.Context, email, hashedPassword string) error {
	_, err := getSession(gtx).
		Table("user").
		Where("email = ?", email).
		Update(map[string]interface{}{
			"password": hashedPassword,
			"updated":  time.Now().UTC(),
		})
	return err
}
