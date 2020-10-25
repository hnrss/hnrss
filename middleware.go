package main

import (
	"github.com/gin-gonic/gin"
)

func SetFormat(fmt string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("format", fmt)
		c.Next()
	}
}
