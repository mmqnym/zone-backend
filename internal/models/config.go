package models

type RateLimit struct {
	Interval    int `json:"interval"`
	IntervalNum int `json:"intervalNum"`
	Limit       int `json:"limit"`
}

type RateLimitsConfig struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}
