package rateLimit

import (
	"net/http"
	"time"
	"zoneBackend/config"

	"github.com/gin-gonic/gin"
)

// map: exchange => rateType => time
var rateLimitsReg map[string]map[string]int

// map: path => weight
var endpointsWeight map[string]int

var startWeightTime time.Time
var startRateLimitTime time.Time

func init() {
	endpointsWeight = map[string]int{
		"/spot/balance":          5,
		"/spot/transfer/records": 5,
		"/futures/balance":       5,
	}

	rateLimitsReg = map[string]map[string]int{}

	for exchange := range config.RateLimits {
		rateLimitsReg[exchange] = map[string]int{
			"REQUEST_WEIGHT": 0,
			"RAW_REQUESTS":   0,
		}
	}

	go resetReg()
}

// reset reg at the end of each interval
func resetReg() {
	startWeightTime = time.Now()
	startRateLimitTime = time.Now()

	for {
		time.Sleep(10 * time.Second)
		currentTime := time.Now()

		for exchange := range config.RateLimits {
			weightTimeLimit := time.Duration(config.RateLimits[exchange]["REQUEST_WEIGHT"].IntervalNum*config.RateLimits[exchange]["REQUEST_WEIGHT"].Interval) * time.Second
			if currentTime.Sub(startWeightTime) >= weightTimeLimit {
				rateLimitsReg[exchange]["REQUEST_WEIGHT"] = 0
				startWeightTime = time.Now()
			}

			rateTimeLimit := time.Duration(config.RateLimits[exchange]["RAW_REQUESTS"].IntervalNum*config.RateLimits[exchange]["RAW_REQUESTS"].Interval) * time.Second
			if currentTime.Sub(startRateLimitTime) >= rateTimeLimit {
				rateLimitsReg[exchange]["RAW_REQUESTS"] = 0
				startRateLimitTime = time.Now()
			}
		}
	}
}

func Resolve(c *gin.Context) {
	// only check pathes in endpointsWeight
	if weight, ok := endpointsWeight[c.Request.URL.Path]; ok {
		exchangeParam, _ := c.Get("exchange") // pre-checked at auth middleware
		exchange := exchangeParam.(string)

		if _, ok := rateLimitsReg[exchange]; !ok {
			// not exist exchange
			c.Next()
		}

		rateLimitsReg[exchange]["REQUEST_WEIGHT"] += weight
		rateLimitsReg[exchange]["RAW_REQUESTS"]++

		if rateLimitsReg[exchange]["REQUEST_WEIGHT"] > config.RateLimits[exchange]["REQUEST_WEIGHT"].Limit ||
			rateLimitsReg[exchange]["RAW_REQUESTS"] > config.RateLimits[exchange]["RAW_REQUESTS"].Limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			c.Abort()
			return
		}
	}
	c.Next()
}
