package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"example.com/event-booking/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context)  {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing Auth header!"})
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		err_msg := "Unauthorized user!"
		fmt.Printf("%v %v\n", err_msg, err)
		if strings.Contains(err.Error(), "token is expired") {
			err_msg = "Auth Token Expired!"
		}
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err_msg})
		return
	}

	context.Set("userId", userId)
	context.Next()
}