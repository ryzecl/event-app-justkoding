package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
	"github.com/imagekit-developer/imagekit-go/v2/option"
)

func initImageKit() *imagekit.Client {
	client := imagekit.NewClient(
		option.WithPrivateKey(os.Getenv("IMAGEKIT_PRIVATE_KEY")),
	)
	return &client
}

func CreateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Menerima File form data
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image is required",
		})
		return
	}
	defer file.Close()

	// 1. Upload file to ImageKit
	fileName := header.Filename
	ik := initImageKit()
	uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File: file,
		FileName: fileName,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Image upload failed imageKit",
		})
		return
	}

	parsedTime, _ := time.Parse(time.RFC3339, c.PostForm("datetime"))

	// Simpan ke database
	event := models.Event {
		Name: c.PostForm("name"),
		Description: c.PostForm("description"),
		Location: c.PostForm("location"),
		Datetime: parsedTime,
		Image: uploadRes.URL,
		ImageID: uploadRes.FileID,
		UserID: userID.(uint),
	}

	config.DB.Create(&event)
	c.JSON(http.StatusCreated, gin.H{
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

func UpdateEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error;
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event data not found",
		})
		return
	}

	if event.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this event",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close();

		ik := initImageKit()

		// Upload file new Image
		fileName := header.Filename
		uploadRes, errUpload := ik.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File: file,
			FileName: fileName,
		})

		if errUpload == nil {
			// Hapus old image
			if event.ImageID != "" {
				ik.Files.Delete(context.Background(), event.ImageID)
			}

			// Upload new image
			event.Image = uploadRes.URL
			event.ImageID = uploadRes.FileID
		}
	}

	if name := c.PostForm("name"); name != "" {
		event.Name = name
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if dateTimeStr := c.PostForm("datetime"); dateTimeStr != "" {
		parseTime, errParse := time.Parse(time.RFC3339, dateTimeStr)
		if errParse == nil {
			event.Datetime = parseTime
		}
	}

	config.DB.Save(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
		"event": event,
	})
}

func DeleteEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var event models.Event
	paramsId := c.Param("id")

	var eventData = config.DB.First(&event, paramsId).Error;
	if eventData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event data not found",
		})
		return
	}


	if event.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not the owner of this event",
		})
		return
	}

	if event.ImageID != "" {
		ik:= initImageKit()
		ik.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data deleted successfully",
	})
}