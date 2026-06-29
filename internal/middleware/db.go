package middleware

import (
	"github.com/example/golang-api/pkg/connect"
	"github.com/example/golang-api/pkg/ctv"
	"github.com/gin-gonic/gin"
)

const dbname = "APP"

func AssignDBSession(gtx *gin.Context) {
	eng, err := connect.NewDB(dbname)
	if err != nil {
		gtx.AbortWithStatusJSON(500, gin.H{"message": "db connection failed"})
		return
	}

	session := eng.NewSession().Context(gtx.Request.Context())
	ctv.Set(gtx, ctv.DBSession, session)

	gtx.Next()
}

func CloseDBSession(gtx *gin.Context) {
	gtx.Next()

	session, ok := ctv.Get[interface{ Close() }](gtx, ctv.DBSession)
	if ok {
		session.Close()
	}
}
