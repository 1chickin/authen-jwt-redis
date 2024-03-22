package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/1chickin/authen-jwt-redis/config"
	"github.com/1chickin/authen-jwt-redis/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	// get username & password from request
	var requestBody struct {
		Username string
		Password string
	}

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load request body!",
		})
		return
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password!",
		})
		return
	}

	// create user
	user := &model.User{Username: requestBody.Username, Password: string(hashedPassword)}
	result := config.DB.Create(&user) // pass pointer of data to Create
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user!",
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"password": user.Password,
	})
}

func Login(c *gin.Context) {
	// get username & password from request
	var requestBody struct {
		Username string
		Password string
	}

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load request body!",
		})
		return
	}

	// check username exist
	var user model.User
	config.DB.Where("username = ?", requestBody.Username).First(&user)
	// result := config.DB.First(&user, "username = ?", requestBody.Username)
	// ref sql injection: https://gorm.io/docs/security.html

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username or password was wrong!",
		})
		return
	}

	// check password to user hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username or password was wrong!",
		})
		return
	}
	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     user.ID,
		"expire": time.Now().Add(time.Second * 30).Unix(),
	})

	// sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to generate token!",
		})
		return
	}

	// set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60, "", "", false, true)

	// response
	c.JSON(http.StatusOK, gin.H{})
}

func ValidateToken(c *gin.Context) {
	// response
	c.JSON(http.StatusOK, gin.H{})
}
