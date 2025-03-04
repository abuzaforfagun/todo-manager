package auth_handlers

import (
	"log"
	"net/http"
	"restful-service/models"
	auth_repository "restful-service/repositories/auth"
	utils_encryption "restful-service/utils/encryption"
	"sync"
	"time"

	_ "restful-service/docs"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	repository auth_repository.AuthRepository
}

func NewHandler(repository auth_repository.AuthRepository) AuthHandler {
	return &authHandler{
		repository: repository,
	}
}

var EncryptionKey string
var JwtKey string

// @Summary Register user
// @Description Register new user
// @Tags user
// @Param user body models.UserLoginDto true "Registration payload"
// @Produce json
// @Success 200
// @Router /user/register [post]
func (h authHandler) Register(c *gin.Context) {
	var credential models.UserLoginDto

	err := c.BindJSON(&credential)
	if err != nil {
		log.Fatalln(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	var isExistingUserChan = make(chan bool)
	var encryptedPasswordChan = make(chan string)
	var errChan = make(chan error)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer func() {
			close(isExistingUserChan)
			wg.Done()
		}()

		isExistingUser, err := h.repository.HasUser(credential.Username)
		if err != nil {
			errChan <- err
			return
		}
		isExistingUserChan <- isExistingUser

	}()

	wg.Add(1)
	go func() {
		defer func() {
			close(encryptedPasswordChan)
			wg.Done()
		}()
		encryptedPassword, err := utils_encryption.Encrypt(credential.Password, EncryptionKey)
		if err != nil {
			errChan <- err
		}
		encryptedPasswordChan <- encryptedPassword

	}()

	var isExistingUser bool
	var encryptedPassword string

	select {
	case err := <-errChan:
		if err != nil {
			log.Println(err)
			if err.Error() == "context canceled" {
				c.JSON(http.StatusRequestTimeout, gin.H{})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{})
			}
			return
		}
	case isExistingUser = <-isExistingUserChan:
	}

	select {
	case err := <-errChan:
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	case encryptedPassword = <-encryptedPasswordChan:
	}

	if isExistingUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already registered"})
		return
	}

	err = h.repository.Register(credential.Username, encryptedPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary User login
// @Description User login
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.UserLoginDto true "Login payload"
// @Success 201 {object} models.LoginResponse
// @Router /login [post]
func (h authHandler) Login(c *gin.Context) {
	var userDto models.UserLoginDto
	err := c.BindJSON(&userDto)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	user, err := h.repository.GetUser(userDto.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	passwordFromDb, err := utils_encryption.Decrypt(user.Password, EncryptionKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if passwordFromDb != userDto.Password {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &models.Claims{
		UserId:   user.UserId,
		Username: userDto.Username,
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

	result := models.LoginResponse{Token: tokenString}
	c.JSON(http.StatusOK, result)
}
