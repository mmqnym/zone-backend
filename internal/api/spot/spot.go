package spot

import (
	"net/http"
	"strconv"
	"time"
	"zoneBackend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	exchangeParam, _ := c.Get("exchange")

	exchange, ok := exchangeParam.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if exchange == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	balance, err := services.GetBalance(exchange, "spot")

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"free": balance,
	})
}

func GetTransactionRecords(c *gin.Context) {
	exchangeParam, _ := c.Get("exchange")

	exchange, ok := exchangeParam.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	if exchange == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	startTimeParam := c.Query("startTime")
	endTimeParam := c.Query("endTime")
	currentParam := c.Query("current")
	sizeParam := c.Query("size")

	if startTimeParam == "" {
		// default: last 30 days
		startTimeParam = strconv.FormatInt(time.Now().UTC().Add(-30*24*time.Hour).Unix(), 10)
	}

	if endTimeParam == "" {
		// default: now
		endTimeParam = strconv.FormatInt(time.Now().UTC().Unix(), 10)
	}

	if currentParam == "" {
		currentParam = "1"
	}

	if sizeParam == "" {
		sizeParam = "10"
	}

	startTime, err := strconv.ParseInt(startTimeParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "startTime is invalid",
		})
		return
	}

	endTime, err := strconv.ParseInt(endTimeParam, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "endTime is invalid",
		})
		return
	}

	current, err := strconv.ParseInt(currentParam, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "current is invalid",
		})
		return
	}

	size, err := strconv.ParseInt(sizeParam, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "size is invalid",
		})
		return
	}

	// Users can't query records for more than 30 days ago
	if (time.Now().Unix() - startTime) > int64(24*time.Hour*30) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "startTime is invalid",
		})
		return
	}

	// endTime must be greater than startTime
	if endTime < startTime {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "endTime is invalid",
		})
		return
	}

	// current must be greater than 0
	if current < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "current is invalid",
		})
		return
	}

	// size must be greater than 0 and less than or equal to 100
	if size < 1 || size > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "size is invalid",
		})
		return
	}

	result, err := services.GetTransactionRecords(exchange, startTime, endTime, current, size)

	c.JSON(http.StatusOK, gin.H{
		"rows":  result,
		"total": len(result),
	})
}
