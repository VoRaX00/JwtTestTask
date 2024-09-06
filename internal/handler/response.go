package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewErrorResponse(c *gin.Context, code int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(code, gin.H{
		message: message,
	})
}
