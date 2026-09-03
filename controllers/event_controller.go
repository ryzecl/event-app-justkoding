package controllers

import (
	"net/http"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
)

func CreateEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event);
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	event.UserID = 1
	
	config.DB.Create(&event)
	context.JSON(http.StatusCreated, gin.H{
		"message": "Data created successfully",
		"event": event,
	})
}

func GetEvents(context *gin.Context) {
	var events []models.Event

	config.DB.Find(&events)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved successfully",
		"events": events,
	})
}

func GetEventbyId(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error;
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event data not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved successfully",
		"event": event,
	})
}

func UpdateEvent(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error;
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event data not found",
		})
		return
	}

	var input models.Event
	err := context.ShouldBindJSON(&input);
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	config.DB.Model(&event).Updates(input)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
		"event": event,
	})
}

func DeleteEvent(context *gin.Context) {
	var event models.Event
	paramsId := context.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error;
	if eventData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Event data not found",
		})
		return
	}

	config.DB.Unscoped().Delete(&event)
	context.JSON(http.StatusOK, gin.H{
		"message": "Data deleted successfully",
	})
}