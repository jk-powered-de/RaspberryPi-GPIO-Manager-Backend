package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func printError(err error, code int, message string, c *gin.Context) {
	fmt.Println(err)
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
