package weather

import (
	"github.com/example/golang-api/internal/model"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.RouterGroup) {
	r.GET("/weather/new-taipei", handleGetWeather)
}

func RouterProtected(r *gin.RouterGroup) {
	r.GET("/weather/new-taipei/protected", handleGetWeather)
}

func handleGetWeather(gtx *gin.Context) {
	result, err := getWeather(gtx)
	if err != nil {
		model.Fail(gtx, err)
		return
	}
	model.Success(gtx, result)
}
