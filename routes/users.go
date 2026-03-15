package routes

import (
	"fmt"
	"net/http"

	"example.com/event-booking/models"
	"example.com/event-booking/utils"
	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User
	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		err_msg := fmt.Sprintf("Error: Could not parse request data! %v\n", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	err = user.New()

	if err != nil {
		err_msg := "Error: Could not create user!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	success_msg := "Successfully created user!"
	fmt.Println(success_msg)
	context.JSON(http.StatusCreated, gin.H{"message": success_msg})
}

func login(context *gin.Context) {
	var user models.User
	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		err_msg := fmt.Sprintf("Error: Could not parse request data! %v\n", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		failure_msg := "Invalid Credentials!"
		fmt.Printf("%v %v\n", failure_msg, err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": failure_msg})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		failure_msg := "Error while login user!"
		fmt.Printf("%v %v\n", failure_msg, err)
		context.JSON(http.StatusUnauthorized, gin.H{"error": failure_msg})
		return
	}

	success_msg := "Successfully Authorized user!"
	fmt.Println(success_msg)
	context.JSON(http.StatusOK, gin.H{"message": success_msg, "token": token})
}