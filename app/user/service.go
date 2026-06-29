package user

import (
	"time"

	"github.com/example/golang-api/pkg/env"
	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func register(gtx *gin.Context, email, password string) error {
	existing, err := findByEmail(gtx, email)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	if existing != nil {
		return errorx.New(errorx.CodeConflict, "email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	return createUser(gtx, &User{
		Email:    email,
		Password: string(hashed),
	})
}

func login(gtx *gin.Context, email, password string) (string, error) {
	u, err := findByEmail(gtx, email)
	if err != nil {
		return "", errorx.New(errorx.CodeInternalServerError, err)
	}
	// 不揭露是帳號還是密碼錯誤
	if u == nil {
		return "", errorx.New(errorx.CodeUnauthorized, "account or password is incorrect")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errorx.New(errorx.CodeUnauthorized, "account or password is incorrect")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   u.Email,
		"updated": u.Updated.Format(time.RFC3339),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(env.Getenv("JWT_SECRET")))
}

func changePassword(gtx *gin.Context, email, oldPassword, newPassword string) error {
	u, err := findByEmail(gtx, email)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}
	if u == nil {
		return errorx.New(errorx.CodeUnauthorized, "account or password is incorrect")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword)); err != nil {
		return errorx.New(errorx.CodeBadRequest, "old password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errorx.New(errorx.CodeInternalServerError, err)
	}

	return updatePassword(gtx, email, string(hashed))
}
