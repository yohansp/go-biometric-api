package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	datastore "github.com/yohansp/go-biometric-api/internal/datastore/sql"
	"github.com/yohansp/go-biometric-api/internal/utils"
)

type SharedKeyRequest struct {
	Pin string
}

type SharedKeyResponse struct {
	SharedKey string `json:"shared_key"`
}

func HandlerSettingRoute(app *fiber.App) {
	app.Patch("auth/biometric/sharedkey", func(c *fiber.Ctx) error {
		var accessToken = c.GetReqHeaders()["Authorization"][0]
		accessToken = accessToken[7:]
		var bioClaim BioAuthClaim
		token, err := jwt.ParseWithClaims(accessToken, &bioClaim, func(t *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			panic(fiber.NewError(401, "Invalid token."))
		}

		if !token.Valid {
			panic(fiber.NewError(401, "Invalid token."))
		}

		var userCredential datastore.UserCredential
		dbErr := datastore.Db.Where("phone_number=?", bioClaim.UserName).First(&userCredential).Error
		if dbErr != nil {
			panic(fiber.NewError(401, "Invalid token, user not found."))
		}

		c.Locals("x-user-id", userCredential.Id)
		c.Locals("x-user-phone-number", userCredential.PhoneNumber)
		c.Locals("x-user-pin", userCredential.Pin)
		return c.Next()
	}, updateSharedKey)
}

func updateSharedKey(c *fiber.Ctx) error {
	var requestData SharedKeyRequest
	c.BodyParser(&requestData)

	// check & validate the user Pin
	xUserId := c.Locals("x-user-id").(int)
	xUserPin := c.Locals("x-user-pin").(string)
	fmt.Println("pin: ", requestData.Pin, " : ", xUserPin)
	if requestData.Pin != xUserPin {
		panic(fiber.NewError(401, "Pin is not valid."))
	}

	var responseData SharedKeyResponse
	responseData.SharedKey = utils.GenerateAesKey()
	//userId, _ := strconv.Atoi(xUserId)
	//var userCredential = datastore.UserCredential{Id: xUserId, SharedKey: responseData.SharedKey}
	//datastore.Db.Save(&userCredential)
	datastore.Db.Model(&datastore.UserCredential{}).Where("id=?", xUserId).Update("shared_key", responseData.SharedKey)

	return c.JSON(responseData)
}
