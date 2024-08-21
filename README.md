# Zone-Backend 前測

### 簡述

1. 完成了以下的 API：
    
    - `/exchangeInfo`
    - `/spot/balance`
    - `/spot/transfer/records`
    - `/futures/balance `

2. 完成並套用以下 middleware 分別用於：
    
    - `auth`：將 `basic token` 視為以 `Base64` 編碼後的交易所名稱，呼叫者可以利用該 `token` 存取後端的路由，並取得相應的資料。
    - `rateLimit`：在完成基本 API 需求開發後，替 API 加上權重以及速率上限，用於限制請求頻率。

### 測試

可以於根目錄下運行 `go test -v` 來進行整合測試，其中測試檔撰寫於 `main_test.go`，模擬資料則是在 `./internal/db/data.go` 以及 `./config/config.go`。


### 運行邏輯

`通用前置檢查`：

```js
if (!token) {
    // head 必須攜帶 token，並且 token 必須是合法 base64 值
    return "402 Unauthorize"
}

if (weight[path]) {
    // 若指定端點存在權重配置，更新暫存器內的值
    // 之後檢查當前權重累積值或是請求數量是否有超過設定值
    reg[weight] += weight[path]
    reg[rawRequests]++

    if (reg[weight] > config[path][weight] || reg[rawRequests] > config[path][rawRequests]) {
        return "429 Too many requests"
    }
}

// 通過檢查，前往 api handler
next()
```

`其他`：

```go
// 伺服器會在背景(./internal/services/services.go)定期清除超過 6 年的現貨交易資料
func init() {
	go autoCleanupTransactionRecords()
}

func autoCleanupTransactionRecords() {
	for {
		db.CleanUpTransactionRecords()
		time.Sleep(1 * time.Hour)
	}
}
```

```go
// 伺服器會在背景(./internal/middleware/rateLimit/resolve.go)根據每間交易所設定的 interval，
// 定期重置速率限制暫存器以及權重限制暫存器
func init() {
	go resetReg()
}

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
```

`/exchangeInfo`：

```js
// 根據 token 決定從 config 回傳哪個交易所的設定

if (!exchange) {
    return "404 not found"
}

return {
  "timezone": "UTC",
  "serverTime": 1565246363776,
  "rateLimits": [
     {
        "rateLimitType": "REQUEST_WEIGHT",
        "interval": "MINUTE",
        "intervalNum": 1,
        "limit": 1200
     },
     {
        "rateLimitType": "RAW_REQUESTS",
        "interval": "MINUTE",
        "intervalNum": 5,
        "limit": 6100
     }
  ]
}
```

`/spot/balance`：
`/futures/balance`：

```js
// 根據 token 決定從模擬資料庫回傳哪個交易所的現貨／合約帳戶 USDT 餘額

if (!exchange) {
    return "404 not found"
}

return {
  "free": "10.12345"
}
```

`/transfer/records`：

```js
// 根據 token 決定從模擬資料庫回傳哪個交易所的現貨帳戶交易紀錄

if (!exchange) {
    return {
        "rows": null,
        "total": 0
    }
}

// 從請求路由端點獲取 query 參數
// startTime, endTime, current, size

// 進行資料前處理：確保 endTime > startTime、startTime 必須在 30 天內等
if (...) {}

// 去模擬資料庫取得相關紀錄
// 查詢表的資料必須介於 [startTime, endTime] 並且 <= 30 天
if (record.Timestamp >= timeLimit && record.Timestamp >= startTime && record.Timestamp <= endTime) {
	filteredRecords = append(filteredRecords, record)
}

// 之後根據 current 以及 size 參數決定回傳的資料
// i 代表起始 index，由 第幾頁 * 每頁大小 決定
for (let i = (current * size) - size; i < len(filteredRecords); i++) {
    res = append(res, filteredRecords[i])
    count++

    if count === size {
        break
    }
}

return {
   "rows": res,
   "total": len(res)
}
```
