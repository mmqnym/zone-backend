package config

import (
	"zoneBackend/internal/models"
)

var RateLimits map[string]map[string]models.RateLimit

func init() {
	// map: exchange => rateType => rateLimitConfig
	RateLimits = map[string]map[string]models.RateLimit{
		"binance": {
			"REQUEST_WEIGHT": {
				Interval:    60,
				IntervalNum: 1,
				Limit:       1200,
			},
			"RAW_REQUESTS": {
				Interval:    60,
				IntervalNum: 5,
				Limit:       6100,
			},
		},
		"okx": {
			"REQUEST_WEIGHT": {
				Interval:    60,
				IntervalNum: 1,
				Limit:       20,
			},
			"RAW_REQUESTS": {
				Interval:    60,
				IntervalNum: 1,
				Limit:       10,
			},
		},
	}
}
