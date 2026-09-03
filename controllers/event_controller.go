package controllers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/imagekit-developer/imagekit-go/v2"
)

func validateImage(headerFilename string, size int64) error {
	// Maksimal 2MB
	const maxSize = 2 * 1024 * 1024
	if size > maxSize {
		return fmt.Errorf("file size exceeds 2MB limit")
	}

	// Cek ekstensi file
	ext := strings.ToLower(filepath.Ext(headerFilename))
	allowedExts := map[string]bool {
		".jpg": true,
		".jpeg": true,
		".png": true,
		".webp": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("only .jpg, .jpeg, .png, and .webp files are allowed")
	}

	return nil
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

	// validasi sebelum upload
	if errVal := validateImage(header.Filename, header.Size); errVal != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errVal.Error(),
		})
		return
	}

	// Upload file to ImageKit
	uploadRes, errUpload := config.IK.Files.Upload(context.Background(), imagekit.FileUploadParams{
		File: file,
		FileName: header.Filename,
	})

	if errUpload != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Image upload failed imageKit",
		})
		return
	}

	parsedTime, errTime := time.Parse(time.RFC3339, c.PostForm("datetime"))
	if errTime != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid datetime format (use RFC3339, e.g. 2026-09-04T10:00:00Z)",
		})
		return
	}

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

	if errDB := config.DB.Create(&event).Error; errDB != nil {
		config.IK.Files.Delete(context.Background(), uploadRes.FileID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save event",
		})
		return
	}

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

		// Validasi
		if errVal := validateImage(header.Filename, header.Size); errVal != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errVal.Error(),
			})
			return
		}

		// Upload file new Image
		uploadRes, errUpload := config.IK.Files.Upload(context.Background(), imagekit.FileUploadParams{
			File: file,
			FileName: header.Filename,
		})

		if errUpload == nil {
			// Hapus old image
			if event.ImageID != "" {
				config.IK.Files.Delete(context.Background(), event.ImageID)
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
		config.IK.Files.Delete(context.Background(), event.ImageID)
	}

	config.DB.Unscoped().Delete(&event)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data deleted successfully",
	})
}