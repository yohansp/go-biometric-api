package handlers

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	datastore "github.com/yohansp/go-biometric-api/internal/datastore/sql"
	"github.com/yohansp/go-biometric-api/internal/utils"
)

var secretKey = []byte("indonesia")

type LoginRequestDto struct {
	PhoneNumber string `json:"phone_number"`
	Pin         string `json:"pin"`
}

type BioLoginRequestDto struct {
	PhoneNumber string `json:"phone_number"`
	Data        string `json:"data"`
}

type LoginResponseDto struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type BioAuthClaim struct {
	UserName string `json:"string"`
	jwt.RegisteredClaims
}

func HandlerAuthorizationRoute(app *fiber.App) {
	app.Post("/auth/token", createToken)
	app.Post("/auth/biometric", createTokenFromBio)
}

func createToken(c *fiber.Ctx) error {
	var loginRequest LoginRequestDto
	c.BodyParser(&loginRequest)

	var userCredential datastore.UserCredential
	errDb := datastore.Db.Where("phone_number=? and pin=?", loginRequest.PhoneNumber, loginRequest.Pin).First(&userCredential).Error
	if errDb != nil {
		panic(fiber.NewError(fiber.StatusUnauthorized, "User not exist."))
	}

	// create jwt token
	bioClaims := BioAuthClaim{
		userCredential.PhoneNumber,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bioauth",
		},
	}
	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, bioClaims)

	accessToken, errJwt := tokenGenerator.SignedString(secretKey)
	if errJwt != nil {
		panic(errJwt)
	}
	var response = LoginResponseDto{AccessToken: accessToken, RefreshToken: "00011122233"}
	return c.JSON(response)
}

func createTokenFromBio(c *fiber.Ctx) error {
	var iv = "MjAyNC0wMi0wOCAw"
	var bioLoginRequest BioLoginRequestDto
	c.BodyParser(&bioLoginRequest)

	var userCredential datastore.UserCredential
	errDb := datastore.Db.Where("phone_number=?", bioLoginRequest.PhoneNumber).First(&userCredential).Error
	if errDb != nil {
		panic(fiber.NewError(fiber.StatusUnauthorized, "User not exist"))
	}

	key, err := base64.StdEncoding.DecodeString(userCredential.SharedKey)
	if err != nil {
		panic(fiber.NewError(fiber.StatusUnauthorized, "Invalidate data"))
	}

	data, err := base64.StdEncoding.DecodeString(bioLoginRequest.Data)
	if err != nil {
		panic(fiber.NewError(fiber.StatusUnauthorized, "Invalidate data"))
	}

	result, err := utils.DecryptCBC(key, []byte(iv), data)
	if err != nil {
		panic(fiber.NewError(fiber.StatusUnauthorized, "Invalidate data"))
	}

	resultString, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resultString))

	// create jwt token
	bioClaims := BioAuthClaim{
		userCredential.PhoneNumber,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bioauth",
		},
	}
	tokenGenerator := jwt.NewWithClaims(jwt.SigningMethodHS256, bioClaims)

	accessToken, errJwt := tokenGenerator.SignedString(secretKey)
	if errJwt != nil {
		panic(errJwt)
	}
	var response = LoginResponseDto{AccessToken: accessToken, RefreshToken: "00011122233"}
	return c.JSON(response)
}

func Test(c *fiber.Ctx) error {
	return c.SendString("OK")
}
