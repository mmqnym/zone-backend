package general

import (
	"net/http"
	"time"
	"zoneBackend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetExchangeInfo(c *gin.Context) {
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

	rateLimits, err := services.GetRateLimits(exchange)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timezone":   "UTC",
		"serverTime": time.Now().UTC().UnixMilli(),
		"rateLimits": rateLimits,
	})
}
