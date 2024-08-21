package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"zoneBackend/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestSpotBalance(t *testing.T) {
	useRateLimit = false
	router := setupServer()

	tests := []struct {
		subject        string
		exchange       string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			subject:        "Valid exchange with free: 1000.0",
			exchange:       "binance",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"free": "1000.0"},
		},
		{
			subject:        "Valid exchange with free: 1234.5678",
			exchange:       "okx",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"free": "1234.5678"},
		},
		{
			subject:        "Missing exchange",
			exchange:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Unauthorized"},
		},
		{
			subject:        "Non-exist exchange",
			exchange:       "non-exist",
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "Not found"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.subject, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/spot/balance", nil)
			token := base64.StdEncoding.EncodeToString([]byte(testCase.exchange))
			req.Header.Set("Authorization", "Basic "+token)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedBody, response)
		})
	}
}

func TestSpotTransactionRecords(t *testing.T) {
	useRateLimit = false
	router := setupServer()

	tests := []struct {
		subject        string
		exchange       string
		startTime      string
		endTime        string
		current        string
		size           string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			subject:        "Valid search with default params",
			exchange:       "binance",
			startTime:      "",
			endTime:        "",
			current:        "",
			size:           "",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"rows": []models.TransactionRecord{
					{
						Amount:    "1.50000000",
						Asset:     "ETH",
						Status:    "CONFIRMED",
						Timestamp: int64(1723111100),
						TxId:      int64(6339808703),
						Type:      "IN",
					},
					{
						Amount:    "0.52000000",
						Asset:     "BTC",
						Status:    "CONFIRMED",
						Timestamp: int64(1723112132),
						TxId:      int64(6339808721),
						Type:      "IN",
					},
					{
						Amount:    "0.52000000",
						Asset:     "BTC",
						Status:    "CONFIRMED",
						Timestamp: int64(1723115155),
						TxId:      int64(6339818721),
						Type:      "OUT",
					},
					{
						Amount:    "104.50000000",
						Asset:     "MATIC",
						Status:    "CONFIRMED",
						Timestamp: int64(1723117190),
						TxId:      int64(6339819955),
						Type:      "IN",
					},
					{
						Amount:    "104.50000000",
						Asset:     "LINK",
						Status:    "CONFIRMED",
						Timestamp: int64(1723127890),
						TxId:      int64(6339834487),
						Type:      "IN",
					},
					{
						Amount:    "88.50000000",
						Asset:     "CRO",
						Status:    "CONFIRMED",
						Timestamp: int64(1723128894),
						TxId:      int64(6339844441),
						Type:      "OUT",
					},
					{
						Amount:    "666.12340000",
						Asset:     "XMR",
						Status:    "CONFIRMED",
						Timestamp: int64(1723546282),
						TxId:      int64(6359866376),
						Type:      "OUT",
					},
					{
						Amount:    "49871.00000000",
						Asset:     "XRP",
						Status:    "CONFIRMED",
						Timestamp: int64(1723599334),
						TxId:      int64(6359867371),
						Type:      "IN",
					},
					{
						Amount:    "888.00000000",
						Asset:     "OKB",
						Status:    "CONFIRMED",
						Timestamp: int64(1723749738),
						TxId:      int64(6359852774),
						Type:      "OUT",
					},
					{
						Amount:    "32.00000000",
						Asset:     "SOL",
						Status:    "FAILED",
						Timestamp: int64(1723759331),
						TxId:      int64(6359875338),
						Type:      "IN",
					},
				},
				"total": 10,
			},
		},
		{
			subject:        "Valid search with set params",
			exchange:       "binance",
			startTime:      "1723117160",
			endTime:        "1723546285",
			current:        "2",
			size:           "3",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"rows": []models.TransactionRecord{
					{
						Amount:    "666.12340000",
						Asset:     "XMR",
						Status:    "CONFIRMED",
						Timestamp: int64(1723546282),
						TxId:      int64(6359866376),
						Type:      "OUT",
					},
				},
				"total": 1,
			},
		},
		{
			subject:        "non-exist exchange account",
			exchange:       "max",
			startTime:      "",
			endTime:        "",
			current:        "",
			size:           "",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"rows":  nil,
				"total": 0,
			},
		},
		{
			subject:        "invalid search",
			exchange:       "binance",
			startTime:      "",
			endTime:        "",
			current:        "",
			size:           "101",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			subject:        "invalid search",
			exchange:       "binance",
			startTime:      "2",
			endTime:        "1",
			current:        "",
			size:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.subject, func(t *testing.T) {
			req, _ := http.NewRequest(
				http.MethodGet,
				"/spot/transfer/records"+
					"?startTime="+testCase.startTime+
					"&endTime="+testCase.endTime+
					"&current="+testCase.current+
					"&size="+testCase.size+"",
				nil,
			)
			token := base64.StdEncoding.EncodeToString([]byte(testCase.exchange))
			req.Header.Set("Authorization", "Basic "+token)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)

			// invalid cases
			if testCase.expectedStatus != http.StatusOK {
				return
			}

			// valid cases
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// convert response to expected format
			convertedResponse := map[string]interface{}{
				"total": int(response["total"].(float64)),
				"rows":  nil,
			}

			if convertedResponse["total"].(int) != 0 {
				convertedResponse["rows"] = make([]models.TransactionRecord, len(response["rows"].([]interface{})))

				for i, row := range response["rows"].([]interface{}) {
					rowMap := row.(map[string]interface{})
					convertedResponse["rows"].([]models.TransactionRecord)[i] = models.TransactionRecord{
						Amount:    rowMap["amount"].(string),
						Asset:     rowMap["asset"].(string),
						Status:    rowMap["status"].(string),
						Timestamp: int64(rowMap["timestamp"].(float64)),
						TxId:      int64(rowMap["txId"].(float64)),
						Type:      rowMap["type"].(string),
					}
				}
			}

			assert.Equal(t, testCase.expectedBody, convertedResponse)
		})
	}
}

func TestFutureBalance(t *testing.T) {
	useRateLimit = false
	router := setupServer()

	tests := []struct {
		subject        string
		exchange       string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			subject:        "Valid exchange with free: 33333.0",
			exchange:       "bingx",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"free": "33333.0"},
		},
		{
			subject:        "Valid exchange with free: 7654.5678",
			exchange:       "okx",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"free": "7654.5678"},
		},
		{
			subject:        "Missing exchange",
			exchange:       "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Unauthorized"},
		},
		{
			subject:        "Non-exist exchange",
			exchange:       "non-exist",
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "Not found"},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.subject, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/futures/balance", nil)
			token := base64.StdEncoding.EncodeToString([]byte(testCase.exchange))
			req.Header.Set("Authorization", "Basic "+token)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedBody, response)
		})
	}
}

func TestExchangeInfo(t *testing.T) {
	useRateLimit = false
	router := setupServer()

	tests := []struct {
		subject        string
		exchange       string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			subject:        "Valid exchange info",
			exchange:       "okx",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"rateLimits": []interface{}{
					map[string]interface{}{
						"rateLimitType": "REQUEST_WEIGHT",
						"interval":      "MINUTE",
						"intervalNum":   float64(1),
						"limit":         float64(20),
					},
					map[string]interface{}{
						"rateLimitType": "RAW_REQUESTS",
						"interval":      "MINUTE",
						"intervalNum":   float64(1),
						"limit":         float64(10),
					},
				}, 
				"serverTime": 1724213841840, 
				"timezone": "UTC",
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.subject, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/exchangeInfo", nil)
			token := base64.StdEncoding.EncodeToString([]byte(testCase.exchange))
			req.Header.Set("Authorization", "Basic "+token)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			
			expectedRateLimits := testCase.expectedBody["rateLimits"].([]interface{})
			actualRateLimits, ok := response["rateLimits"].([]interface{})
			assert.True(t, ok)
			assert.Equal(t, len(expectedRateLimits), len(actualRateLimits))

			for i, expectedLimit := range expectedRateLimits {
				expectedMap := expectedLimit.(map[string]interface{})
				actualMap, ok := actualRateLimits[i].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, expectedMap["rateLimitType"], actualMap["rateLimitType"])
				assert.Equal(t, expectedMap["interval"], actualMap["interval"])
				assert.Equal(t, expectedMap["intervalNum"], actualMap["intervalNum"])
				assert.Equal(t, expectedMap["limit"], actualMap["limit"])
			}

			_, ok = response["serverTime"].(float64)
			assert.True(t, ok, "serverTime should be present and should be a number")
		})
	}
}

func TestRateLimit(t *testing.T) {
	useRateLimit = true
	router := setupServer()

	tests := []struct {
		subject        string
		exchange       string
		expectedStatus int
		callTimes      int
	}{
		{
			subject:        "Valid request",
			exchange:       "okx",
			expectedStatus: http.StatusOK,
			callTimes:      1,
		},
		{
			subject:        "Too many requests that exceeds set limit",
			exchange:       "okx",
			expectedStatus: http.StatusTooManyRequests,
			callTimes:      4,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.subject, func(t *testing.T) {
			token := base64.StdEncoding.EncodeToString([]byte(testCase.exchange))
			
			var lastRecorder *httptest.ResponseRecorder

			for i := 0; i < testCase.callTimes; i++ {
				req, _ := http.NewRequest(http.MethodGet, "/spot/balance", nil)
				req.Header.Set("Authorization", "Basic "+token)

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				lastRecorder = w
			}
			
			// check the status code after all requests done
			assert.Equal(t, testCase.expectedStatus, lastRecorder.Code)
		})
	}
}
