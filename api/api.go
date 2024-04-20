package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mdtosif/openletter/model"
)

var jwksURL = os.Getenv("JWKS_URL")

type MyCustomClaims struct {
	Sub []string `json:"sub"`
	jwt.StandardClaims
}

func GetLetters(c *gin.Context) {
	token := c.GetHeader("Authorization")
	username, err := GetUserName(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	print(username, token)
	leters := model.GetUserMessage(username)

	// Here you can perform validation, database operations, etc.
	// For simplicity, let's just return the received user data
	c.JSON(http.StatusCreated, leters)
}

func AddLetter(c *gin.Context) {
	var letter model.Letters
	token := c.GetHeader("Authorization")
	username, err := GetUserName(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	print(username, token)
	if err := c.BindJSON(&letter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	letter.From = username

	model.AddLetter(&letter)

	// Here you can perform validation, database operations, etc.
	// For simplicity, let's just return the received user data
	c.JSON(http.StatusCreated, letter)
}

func GetUserName(userToken string) (string, error) {

	// Create a context that, when cancelled, ends the JWKS background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())

	// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
	// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
	// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
	options := keyfunc.Options{
		Ctx: ctx,
		RefreshErrorHandler: func(err error) {
			log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	}

	// Create the JWKS from the resource at the given URL.
	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		cancel()
		return "", errors.New("failed to create jwks from resource")
	}

	// Get a JWT to parse.
	jwtB64 := userToken

	// Parse the JWT.
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtB64, &claims, jwks.Keyfunc)
	if err != nil {
		cancel()
		return "", errors.New("failed to parse the JWT")
	}

	// Check if the token is valid.
	if !token.Valid {
		cancel()
		return "", errors.New("the token is not valid")
	}
	log.Println("The token is valid.")

	// End the background refresh goroutine when it's no longer needed.
	cancel()

	// This will be ineffectual because the line above this canceled the parent context.Context.
	// This method call is idempotent similar to context.CancelFunc.
	jwks.EndBackground()
	// do something with decoded claims
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
		if key == "sub" {
			s := val.(string)
			return s, nil
		}
	}
	return "", errors.New("invalid token")
}
