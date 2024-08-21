package db

import "zoneBackend/internal/models"

var tableSpotBalance map[string]string
var tableFuturesBalance map[string]string

var tableSpotTransactionRecords map[string][]models.TransactionRecord

func init() {
	tableSpotBalance = map[string]string{
		"binance": "1000.0",
		"okx":     "1234.5678",
	}

	tableFuturesBalance = map[string]string{
		"bingx": "33333.0",
		"okx":   "7654.5678",
	}

	tableSpotTransactionRecords = map[string][]models.TransactionRecord{
		"binance": {
			{
				Amount:    "0.10000000",
				Asset:     "BNB",
				Status:    "CONFIRMED",
				Timestamp: 1566898617,
				TxId:      5240372201,
				Type:      "IN",
			},
			{
				Amount:    "5.00000000",
				Asset:     "USDT",
				Status:    "CONFIRMED",
				Timestamp: 1566888436,
				TxId:      5239810406,
				Type:      "OUT",
			},
			{
				Amount:    "1.00000000",
				Asset:     "EOS",
				Status:    "CONFIRMED",
				Timestamp: 1566888403,
				TxId:      5239808703,
				Type:      "IN",
			},
			{
				Amount:    "1.50000000",
				Asset:     "ETH",
				Status:    "CONFIRMED",
				Timestamp: 1723111100,
				TxId:      6339808703,
				Type:      "IN",
			},
			{
				Amount:    "0.52000000",
				Asset:     "BTC",
				Status:    "CONFIRMED",
				Timestamp: 1723112132,
				TxId:      6339808721,
				Type:      "IN",
			},
			{
				Amount:    "0.52000000",
				Asset:     "BTC",
				Status:    "CONFIRMED",
				Timestamp: 1723115155,
				TxId:      6339818721,
				Type:      "OUT",
			},
			{
				Amount:    "104.50000000",
				Asset:     "MATIC",
				Status:    "CONFIRMED",
				Timestamp: 1723117190,
				TxId:      6339819955,
				Type:      "IN",
			},
			{
				Amount:    "104.50000000",
				Asset:     "LINK",
				Status:    "CONFIRMED",
				Timestamp: 1723127890,
				TxId:      6339834487,
				Type:      "IN",
			},
			{
				Amount:    "88.50000000",
				Asset:     "CRO",
				Status:    "CONFIRMED",
				Timestamp: 1723128894,
				TxId:      6339844441,
				Type:      "OUT",
			},
			{
				Amount:    "666.12340000",
				Asset:     "XMR",
				Status:    "CONFIRMED",
				Timestamp: 1723546282,
				TxId:      6359866376,
				Type:      "OUT",
			},
			{
				Amount:    "49871.00000000",
				Asset:     "XRP",
				Status:    "CONFIRMED",
				Timestamp: 1723599334,
				TxId:      6359867371,
				Type:      "IN",
			},
			{
				Amount:    "888.00000000",
				Asset:     "OKB",
				Status:    "CONFIRMED",
				Timestamp: 1723749738,
				TxId:      6359852774,
				Type:      "OUT",
			},
			{
				Amount:    "32.00000000",
				Asset:     "SOL",
				Status:    "FAILED",
				Timestamp: 1723759331,
				TxId:      6359875338,
				Type:      "IN",
			},
			{
				Amount:    "32.00000000",
				Asset:     "SOL",
				Status:    "CONFIRMED",
				Timestamp: 1723767620,
				TxId:      6359847571,
				Type:      "IN",
			},
			{
				Amount:    "777.77777777",
				Asset:     "TON",
				Status:    "PENDING",
				Timestamp: 1723888888,
				TxId:      6359862383,
				Type:      "OUT",
			},
		},
	}
}
