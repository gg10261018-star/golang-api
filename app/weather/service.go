package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/golang-api/pkg/env"
	"github.com/example/golang-api/pkg/errorx"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func getWeather(gtx *gin.Context) (*WeatherResult, error) {
	am, err := findByPeriod(gtx, "AM")
	if err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}
	pm, err := findByPeriod(gtx, "PM")
	if err != nil {
		return nil, errorx.New(errorx.CodeInternalServerError, err)
	}

	result := &WeatherResult{}
	if am != nil {
		result.AM = &periodData{MinT: am.MinT, MaxT: am.MaxT}
	}
	if pm != nil {
		result.PM = &periodData{MinT: pm.MinT, MaxT: pm.MaxT}
	}
	return result, nil
}

func StartScheduler() {
	go fetchAndStore()

	c := cron.New()
	c.AddFunc("@every 24h", fetchAndStore)
	c.Start()
}

func fetchAndStore() {
	apiKey := env.Getenv("CWA_API_KEY")
	if apiKey == "" {
		log.Println("[weather] CWA_API_KEY not set, skipping")
		return
	}

	url := fmt.Sprintf(
		"https://opendata.cwa.gov.tw/api/v1/rest/datastore/F-C0032-001?Authorization=%s&locationName=新北市",
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[weather] fetch error: %v", err)
		return
	}
	defer resp.Body.Close()

	var data cwaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("[weather] decode error: %v", err)
		return
	}

	for _, loc := range data.Records.Location {
		if loc.LocationName != "新北市" {
			continue
		}

		minT := map[string]float64{}
		maxT := map[string]float64{}

		for _, elem := range loc.WeatherElement {
			for _, t := range elem.Time {
				period := classifyPeriod(t.StartTime)
				if period == "" {
					continue
				}
				var val float64
				fmt.Sscanf(t.Parameter.ParameterValue, "%f", &val)
				switch elem.ElementName {
				case "MinT":
					minT[period] = val
				case "MaxT":
					maxT[period] = val
				}
			}
		}

		for _, period := range []string{"AM", "PM"} {
			if err := upsertWeather(&Weather{
				City:   "新北市",
				MinT:   minT[period],
				MaxT:   maxT[period],
				Period: period,
			}); err != nil {
				log.Printf("[weather] upsert error: %v", err)
			}
		}
	}

	log.Println("[weather] data updated")
}

func classifyPeriod(startTime string) string {
	t, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return ""
	}
	h := t.Hour()
	if h >= 6 && h < 12 {
		return "AM"
	}
	if h >= 12 && h < 18 {
		return "PM"
	}
	return ""
}
