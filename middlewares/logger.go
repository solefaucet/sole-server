package middlewares

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"github.com/solefaucet/sole-server/models"
)

// Logger returns a middleware that logs all request
func Logger(geo *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		// mark start time
		start := time.Now()

		// process request
		c.Next()

		// log after processing
		status := c.Writer.Status()
		ip := c.ClientIP()
		loc := getLocationFromIP(geo, ip)

		fields := logrus.Fields{
			"event":            models.EventHTTPRequest,
			"method":           c.Request.Method,
			"path":             c.Request.URL.Path,
			"query":            c.Request.URL.Query().Encode(),
			"http_status_code": status,
			"response_time":    float64(time.Since(start).Nanoseconds()) / 1e6,
			"ip":               ip,
			"location":         loc.String(),
			"continent":        loc.continent,
			"country":          loc.country,
			"region":           loc.region,
			"city":             loc.city,
		}
		if err := c.Errors.ByType(gin.ErrorTypeAny).Last(); err != nil {
			fields["error"] = err.Error()
		}

		entry := logrus.WithFields(fields)
		if status == http.StatusInternalServerError {
			entry.Error("internal server error")
			return
		}

		entry.Info("succeed to process http request")
	}
}

func getLocationFromIP(geo *geoip2.Reader, ip string) location {
	record, err := geo.City(net.ParseIP(ip))
	loc := location{
		"unknown",
		"unknown",
		"unknown",
		"unknown",
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": models.EventGetGeoFromIP,
			"error": err.Error(),
		}).Error("fail to get city information")
		return loc
	}

	loc.continent = record.Continent.Names["en"]
	loc.country = record.Country.Names["en"]
	loc.city = record.City.Names["en"]
	for _, v := range record.Subdivisions {
		loc.region += v.Names["en"]
	}
	return loc
}

type location struct {
	continent, country, region, city string
}

func (l location) String() string {
	return strings.Join([]string{l.continent, l.country, l.region, l.city}, "")
}
