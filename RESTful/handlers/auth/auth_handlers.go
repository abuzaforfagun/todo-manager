package auth_handlers

import (
	"log"
	"net/http"
	"restful-service/models"
	auth_repository "restful-service/repositories/auth"
	utils_encryption "restful-service/utils/ecryption"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var EncryptionKey string
var JwtKey string

func Register(c *gin.Context) {
	var credential models.Credentials

	err := c.BindJSON(&credential)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	isExistingUser, err := auth_repository.HasUser(credential.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if isExistingUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already registed"})
		return
	}

	encryptedPassword, err := utils_encryption.Encrypt(credential.Password, EncryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	err = auth_repository.Register(credential.Username, encryptedPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
}

func Login(c *gin.Context) {
	var credential models.Credentials
	err := c.BindJSON(&credential)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	user, err := auth_repository.GetUser(credential.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	passwordFromDb, err := utils_encryption.Decrypt(user.Password, EncryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if passwordFromDb != credential.Password {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &models.Claims{
		Username: credential.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		log.Println("ERROR: Unable to sing", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
