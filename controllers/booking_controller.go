package controllers

import (
	"fmt"
	"net/http"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookingInput struct {
	Phone   string `json:"phone" binding:"required"`
	EventID uint   `json:"eventId" binding:"required"`
}

func CreateBookinEvent(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input BookingInput
	var booking models.Booking

	errValidation := c.ShouldBindJSON(&input)
	if errValidation != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errValidation.Error(),
		})
		return
	}

	// Check if event exists
	var event models.Event
	if err := config.DB.First(&event, input.EventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	// If user already booking event id
	errCheckBooking := config.DB.Where("user_id = ? AND event_id = ?", userID.(uint), input.EventID).First(&booking).Error
	if errCheckBooking == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You already booked this event",
		})
		return
	}

	// Generate Booking Code
	codeBooking := fmt.Sprintf("BK-%sE%dU%d", time.Now().Format("20060102"), input.EventID, userID.(uint))

	// Save to database
	bookingData := models.Booking{
		Phone:       input.Phone,
		EventID:     input.EventID,
		BookingCode: codeBooking,
		UserID:      userID.(uint),
	}
	
	errCreateBooking := config.DB.Create(&bookingData).Error
	if errCreateBooking != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to booking the event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registered event successfully",
	})
}

func GetBookingbyUser(c *gin.Context) {
	var booking []models.Booking
	userID, _ := c.Get("userID")

	errBookingData := config.DB.Preload("Event").Preload("Event.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).Where("user_id = ?", userID).Find(&booking).Error
	if errBookingData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
	})
}

func DeleteBooking (c *gin.Context) {
	userID, _ := c.Get("userID")
	var booking models.Booking

	paramsId := c.Param("id")

	bookingData := config.DB.First(&booking, paramsId).Error

	if bookingData != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Event not found",
		})
		return
	}

	if booking.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You're not allowed delete others booked",
		})
		return
	}

	config.DB.Unscoped().Delete(&booking)
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking event deleted successfully",
	})
}