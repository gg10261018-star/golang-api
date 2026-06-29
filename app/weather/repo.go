package weather

import (
	"time"

	"github.com/example/golang-api/pkg/connect"
	"github.com/example/golang-api/pkg/ctv"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

const dbname = "APP"

func getSession(gtx *gin.Context) *xorm.Session {
	session, _ := ctv.Get[*xorm.Session](gtx, ctv.DBSession)
	return session
}

func findByPeriod(gtx *gin.Context, period string) (*Weather, error) {
	var w Weather
	found, err := getSession(gtx).
		Where("city = ? AND period = ?", "新北市", period).
		OrderBy("updated DESC").
		Get(&w)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &w, nil
}

// upsertWeather 排程用，不走 gin.Context，直接拿 session
func upsertWeather(w *Weather) error {
	eng, err := connect.Engine(dbname)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	existing := &Weather{}
	found, err := eng.Where("city = ? AND period = ?", w.City, w.Period).Get(existing)
	if err != nil {
		return err
	}

	if found {
		_, err = eng.Where("city = ? AND period = ?", w.City, w.Period).
			Update(map[string]interface{}{
				"min_t":   w.MinT,
				"max_t":   w.MaxT,
				"updated": now,
			})
	} else {
		w.Created = now
		w.Updated = now
		_, err = eng.Insert(w)
	}
	return err
}
