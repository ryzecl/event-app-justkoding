package controllers

import (
	"net/http"
	"os"
	"time"

	"example.com/event-app/config"
	"example.com/event-app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInputRegister struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthInputLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func RegisterUser(context *gin.Context) {
	var input AuthInputRegister

	// Validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash Password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if errHash != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create User
	user := models.User{
		Name: input.Name,
		Email: input.Email,
		Password: string(hashedPassword),
	}
	
	userCreated := config.DB.Create(&user).Error
	if userCreated != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already registered",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"events": user.Events,
		},
	})
}

func LoginUser(context *gin.Context) {
	var input AuthInputLogin

	// Validation
	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	userData := config.DB.Where("email = ?", input.Email).First(&user).Error
	if userData != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email not registered",
		})
		return
	}

	errMatchPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if errMatchPassword != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Incorrect password",
		})
		return
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token": tokenString,
		"user": gin.H{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"events": user.Events,
		},
	})
}

func GetCurrentUser(context *gin.Context) {
	// Get user ID
	userID, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var user models.User
	userData := config.DB.Select("id", "name", "email", "events").First(&user, userID).Error
	if userData != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}