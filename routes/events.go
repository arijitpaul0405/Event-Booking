package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"example.com/event-booking/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	userId := context.GetInt64("userId")

	events, err := models.GetAllEvents(userId)

	if err != nil {
		err_msg := "Error: Could not retrieve events!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
	}

	fmt.Println("Successfully retrieved all the events!")

	context.JSON(http.StatusOK, gin.H{"message": events})
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindBodyWithJSON(&event)

	if err != nil {
		err_msg := fmt.Sprintf("Error: Could not parse request data! %v\n", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	event.UserID = context.GetInt64("userId")

	err = event.Save()

	if err != nil {
		err_msg := "Error: Could not create event!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	fmt.Println("Successfully created event!")

	context.JSON(http.StatusCreated, gin.H{"message": "Event Created!", "event": event})
}

func getEvent(context *gin.Context) {
	event_id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	userid := context.GetInt64("userId")

	if err != nil {
		err_msg := "Error: Could not parse event id from request! Please check your input."
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	var event *models.Event
	event, err = models.GetEventByID(event_id, userid)

	if err != nil {
		err_msg := fmt.Sprintf("Not Found: Event with id %v not found!", event_id)
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusNotFound, gin.H{"message": err_msg})
		return
	}

	fmt.Printf("Successfully retrieved event with id %v!\n", event_id)

	context.JSON(http.StatusOK, gin.H{"message": event})
}

func updateEvent(context *gin.Context) {
	event_id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		err_msg := "Error: Could not parse event id from request! Please check your input."
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	user_id := context.GetInt64("userId")

	var retrieved_event *models.Event
	retrieved_event, err = models.GetEventByID(user_id, event_id)

	if err != nil {
		err_msg := fmt.Sprintf("Not Found: Event with id %v not found!", event_id)
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusNotFound, gin.H{"message": err_msg})
		return
	}

	if retrieved_event.UserID != context.GetInt64("userId") {
		err_msg := fmt.Sprintf("Not authorized to update the event with id %v", event_id)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err_msg})
		return
	}

	var event *models.Event
	err = context.ShouldBindBodyWithJSON(&event)

	// updated event would have the same old event id
	event.ID = event_id

	if err != nil {
		err_msg := "Error: Could not parse request data!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	err = event.UpdateByID()

	if err != nil {
		err_msg := fmt.Sprintf("Error: Could not update event with id %v!", event_id)
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	success_msg := fmt.Sprintf("Successfully updated event with id %v", event_id)

	fmt.Println(success_msg)

	context.JSON(http.StatusOK, gin.H{"message": success_msg})
}

func deleteEvent(context *gin.Context) {
	event_id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		err_msg := "Error: Could not parse event id from request! Please check your input."
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	user_id := context.GetInt64("userId")

	var retrieved_event *models.Event
	retrieved_event, err = models.GetEventByID(event_id, user_id)

	if err != nil {
		err_msg := fmt.Sprintf("Not Found: Event with id %v not found!", event_id)
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusNotFound, gin.H{"message": err_msg})
		return
	}

	if retrieved_event.UserID != context.GetInt64("userId") {
		err_msg := fmt.Sprintf("Not authorized to delete the event with id %v", event_id)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err_msg})
		return
	}

	err = retrieved_event.DeleteByID()

	if err != nil {
		err_msg := fmt.Sprintf("Error: Could not delete event with id %v!", event_id)
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	success_msg := fmt.Sprintf("Successfully deleted event with id %v", event_id)

	fmt.Println(success_msg)

	context.JSON(http.StatusOK, gin.H{"message": success_msg})
}

func registerEvent(context *gin.Context) {
	event_id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		err_msg := "Error: Could not parse event id from request! Please check your input."
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	user_id := context.GetInt64("userId")

	var retrieved_event *models.Event
	retrieved_event, err = models.GetEventByID(event_id, user_id)

	if err != nil {
		err_msg := "Event does not exists!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusNotFound, gin.H{"error": err_msg})
		return
	}

	if retrieved_event.UserID != user_id {
		err_msg := fmt.Sprintf("User not allowed to register event for event id %v", event_id)
		fmt.Println(err_msg)
		context.JSON(http.StatusUnauthorized, gin.H{"error": err_msg})
		return
	}

	registeration_id, err := retrieved_event.Register(user_id)

	if err != nil {
		err_msg := fmt.Sprintf("Failed to register event with id %v!", event_id)
		if strings.Contains(err.Error(), "already registered") {
			err_msg += " " + err.Error()
		}
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	success_msg := "Successfully registered event!"
	fmt.Println(success_msg)
	context.JSON(http.StatusCreated, gin.H{"message": success_msg, "registeration_id": registeration_id})
}

func getRegistrationByID(context *gin.Context) {
	user_id := context.GetInt64("userId")

	var registered_event *[]models.ResultRegisteredEvent
	registered_event, err := models.GetRegistrationByUserID(user_id)

	if err != nil {
		err_msg := "Could not find registeration for the given user!"
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusNotFound, gin.H{"message": err_msg})
		return
	}

	success_msg := "Successfully retrieved registration for the given user!"
	fmt.Println(success_msg)
	context.JSON(http.StatusOK, gin.H{"message": success_msg, "registered_event": registered_event})
}

func cancelEvent(context *gin.Context) {
	registeration_id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		err_msg := "Error: Could not parse registeration id from request! Please check your input."
		fmt.Printf("%v %v\n", err_msg, err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err_msg})
		return
	}

	user_id := context.GetInt64("userId")

	// var registered_event *models.RegisteredEvent
	// registered_event, err = models.GetRegistrationByID(registeration_id)

	// if err != nil {
	// 	err_msg := fmt.Sprintf("Could not find registeration with id %v!", registeration_id)
	// 	fmt.Printf("%v %v\n", err_msg, err)
	// 	context.JSON(http.StatusNotFound, gin.H{"error": err_msg})
	// 	return
	// }

	err = models.CancelRegisteration(registeration_id, user_id)

	if err != nil {
		err_msg := fmt.Sprintf("Failed to cancel registeration with id %v!", registeration_id)
		fmt.Printf("%v %v\n", err_msg, err)
		if strings.Contains(err.Error(), "No registration ") {
			err_msg = err.Error()
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err_msg})
		return
	}

	success_msg := "Successfully cancelled registeration!"
	fmt.Println(success_msg)
	context.JSON(http.StatusNoContent, gin.H{"message": success_msg})
}
