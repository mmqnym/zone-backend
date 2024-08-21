package auth

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Resolve(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Basic ")

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	exchange, err := base64.StdEncoding.DecodeString(token)

	if err != nil {
		log.Printf("Decoding Base64 error, encoded data: %v", err)
		
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}

	c.Set("exchange", string(exchange))
	c.Next()
}
