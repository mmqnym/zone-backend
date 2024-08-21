package futures

import (
	"net/http"
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

	balance, err := services.GetBalance(exchange, "futures")

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
