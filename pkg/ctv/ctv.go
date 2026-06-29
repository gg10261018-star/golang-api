package ctv

import "github.com/gin-gonic/gin"

type contextKey string

const (
	DBSession contextKey = "db_session"
	EmailKey  contextKey = "email"
)

func Set(gtx *gin.Context, key contextKey, val any) {
	gtx.Set(string(key), val)
}

func Get[T any](gtx *gin.Context, key contextKey) (T, bool) {
	v, ok := gtx.Get(string(key))
	if !ok {
		var zero T
		return zero, false
	}
	val, ok := v.(T)
	return val, ok
}
