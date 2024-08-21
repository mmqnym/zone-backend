package services

import (
	"errors"
	"time"
	"zoneBackend/internal/db"
	"zoneBackend/internal/models"
	"zoneBackend/utils"

	"zoneBackend/config"
)

func init() {
	go autoCleanupTransactionRecords()
}

func GetRateLimits(exchange string) ([]models.RateLimitsConfig, error) {
	if rateLimit, ok := config.RateLimits[exchange]; ok {
		rateLimitsConfig := []models.RateLimitsConfig{}

		for rateType := range rateLimit {
			rateLimitsConfig = append(rateLimitsConfig, models.RateLimitsConfig{
				RateLimitType: rateType,
				Interval:      utils.TimeFormat(rateLimit[rateType].Interval),
				IntervalNum:   rateLimit[rateType].IntervalNum,
				Limit:         rateLimit[rateType].Limit,
			})
		}

		return rateLimitsConfig, nil
	}

	return nil, errors.New("not found")
}

func GetBalance(exchange string, accountType string) (string, error) {
	return db.GetBalance(exchange, accountType)
}

func GetTransactionRecords(exchange string, startTime, endTime, current, size int64) ([]models.TransactionRecord, error) {
	res, err := db.GetTransactionRecords(exchange, startTime, endTime, current, size)

	return res, err
}

// auto cleanup transaction records which timestamp > 6 years
func autoCleanupTransactionRecords() {
	for {
		db.CleanUpTransactionRecords()
		time.Sleep(1 * time.Hour)
	}
}
