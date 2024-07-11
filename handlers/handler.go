package handlers

import (
	"example/users/database"
	"example/users/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.IndentedJSON(http.StatusOK, users)
}

func AddUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error binding data"})
		return
	}

	database.DB.Create(&newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error binding data"})
		return
	}

	var foundUser models.User
	if err := database.DB.Where("name = ?", user.Name).First(&foundUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": foundUser.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Read the secret from the environment variable
	secret := os.Getenv("SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "secret not found"})
		return
	}

	// Sign the token with the secret
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}

	// Set cookie with the token
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenStr, 3600*24, "/", "", false, true)

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}
