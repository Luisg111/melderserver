package weather

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type DWDWeatherWarning struct {
	Features []struct {
		Type_      string
		Properties struct {
			EVENT     string
			SEVERITY  string
			EFFECTIVE time.Time
			EXPIRES   time.Time
		}
	}
}

type DWDWeatherClient struct {
	warnings []Warning
}

func NewDWDWeatherClient() *DWDWeatherClient {
	client := &DWDWeatherClient{}
	client.startServer()
	return client
}

func (client *DWDWeatherClient) startServer() {
	go func() {
		first := true
		for {
			if !first {
				time.Sleep(time.Minute * 5)
			}
			first = false
			res, err := http.Get("https://maps.dwd.de/geoserver/dwd/ows?service=WFS&version=1.0.0&request=GetFeature&typeName=dwd%3AWarnungen_Gemeinden&CQL_FILTER=WARNCELLID%20IN%20(%27808125046%27)&outputFormat=application%2Fjson")
			if err != nil {
				log.Println("Could not get warnings from dwd " + err.Error())
				continue
			}
			var DWDWeatherWarning DWDWeatherWarning
			err = json.NewDecoder(res.Body).Decode(&DWDWeatherWarning)
			res.Body.Close()
			if err != nil {
				log.Println("Could not decode json " + err.Error())
				continue
			}
			var warnings = make([]Warning, len(DWDWeatherWarning.Features))
			for _, feature := range DWDWeatherWarning.Features {
				isActive := time.Now().After(feature.Properties.EFFECTIVE) && time.Now().Before(feature.Properties.EXPIRES)
				warnings = append(warnings, Warning{
					Title:    feature.Properties.EVENT,
					Text:     feature.Properties.SEVERITY,
					Begin:    feature.Properties.EFFECTIVE,
					End:      feature.Properties.EXPIRES,
					IsActive: isActive,
				},
				)
				log.Println(feature.Properties.EVENT)
			}
			client.warnings = warnings

		}
	}()
}

func (client *DWDWeatherClient) GetWarnings() []Warning {
	return client.warnings
}
