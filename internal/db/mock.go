package db

import (
	"errors"
	"time"
	"zoneBackend/internal/models"
)

func GetBalance(exchange string, accountType string) (string, error) {
	var balance string
	var ok bool
	
	if accountType == "spot" {
		balance, ok = tableSpotBalance[exchange]
	} else if accountType == "futures" {
		balance, ok = tableFuturesBalance[exchange]
	} else {
		return "", errors.New("not available account type")
	}

	if !ok {
		return "", errors.New("Not found")
	}

	return balance, nil
}

func GetTransactionRecords(exchange string, startTime, endTime, current, size int64) ([]models.TransactionRecord, error) {
	if _, ok := tableSpotTransactionRecords[exchange]; !ok {
		return nil, errors.New("Not found")
	}

	filteredRecords := []models.TransactionRecord{}
	timeLimit := time.Now().UTC().Add(-30 * 24 * time.Hour).Unix()

	// only pick records which timestamp is in [startTime, endTime] and within 30 days
	for _, record := range tableSpotTransactionRecords[exchange] {
		if record.Timestamp >= timeLimit && record.Timestamp >= startTime && record.Timestamp <= endTime {
			filteredRecords = append(filteredRecords, record)
		}
	}

	res := []models.TransactionRecord{}
	var count int64 = 0

	// pick data with used params
	for i := (current * size) - size; i < int64(len(filteredRecords)); i++ {
		res = append(res, filteredRecords[i])
		count++

		if count == size {
			break
		}
	}

	return res, nil
}

// clean up transaction records if timestamp > 6 years
func CleanUpTransactionRecords() {
	// check all exchanges
	for exchange := range tableSpotTransactionRecords {

		if _, ok := tableSpotTransactionRecords[exchange]; !ok {
			return
		}

		// remove records of an exchange which timestamp > 6 years
		for i := len(tableSpotTransactionRecords[exchange]) - 1; i >= 0; i-- {
			record := tableSpotTransactionRecords[exchange][i]
			timeLimit := time.Now().UTC().Add(-6 * 12 * 30 * 24 * time.Hour).Unix()

			if record.Timestamp < timeLimit {
				tableSpotTransactionRecords[exchange] = append(
					tableSpotTransactionRecords[exchange][:i], tableSpotTransactionRecords[exchange][i+1:]...,
				)
			}
		}
	}
}
