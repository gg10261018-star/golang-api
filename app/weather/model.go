package weather

import "time"

type Weather struct {
	ID      int64     `xorm:"pk autoincr bigint 'id'"`
	City    string    `xorm:"varchar(50) 'city'"`
	MinT    float64   `xorm:"decimal(4,1) 'min_t'"`
	MaxT    float64   `xorm:"decimal(4,1) 'max_t'"`
	Period  string    `xorm:"enum('AM','PM') 'period'"`
	Created time.Time `xorm:"datetime 'created'"`
	Updated time.Time `xorm:"datetime 'updated'"`
}

type WeatherResult struct {
	AM *periodData `json:"AM"`
	PM *periodData `json:"PM"`
}

type periodData struct {
	MinT float64 `json:"min_t"`
	MaxT float64 `json:"max_t"`
}

// CWA API response
type cwaResponse struct {
	Records struct {
		Location []struct {
			LocationName   string `json:"locationName"`
			WeatherElement []struct {
				ElementName string `json:"elementName"`
				Time        []struct {
					StartTime string `json:"startTime"`
					Parameter struct {
						ParameterValue string `json:"parameterValue"`
					} `json:"parameter"`
				} `json:"time"`
			} `json:"weatherElement"`
		} `json:"location"`
	} `json:"records"`
}
