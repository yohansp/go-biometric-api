package handlers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	datastore "github.com/yohansp/go-biometric-api/internal/datastore/sql"
	"github.com/yohansp/go-biometric-api/internal/utils"
)

func HandlerAdminRoute(app *fiber.App) {

	app.Get("/admin/userlist", func(c *fiber.Ctx) error {
		db := datastore.Db
		var listCredential []datastore.UserCredential
		db.Find(&listCredential)
		return c.JSON(listCredential)
	})

	app.Post("/admin/user/addexample", func(c *fiber.Ctx) error {
		passwordMd5 := md5.Sum([]byte("123456"))
		password := hex.EncodeToString(passwordMd5[:])
		datastore.Db.Create(&datastore.UserCredential{PhoneNumber: "081311137368", Pin: password, SharedKey: "", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		datastore.Db.Create(&datastore.UserCredential{PhoneNumber: "081311137369", Pin: password, SharedKey: "", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		datastore.Db.Create(&datastore.UserCredential{PhoneNumber: "081311137370", Pin: password, SharedKey: "", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		datastore.Db.Create(&datastore.UserCredential{PhoneNumber: "081311137372", Pin: password, SharedKey: "", CreatedAt: time.Now(), UpdatedAt: time.Now()})
		var listCredential []datastore.UserCredential
		datastore.Db.Find(&listCredential)
		return c.JSON(listCredential)
	})

	app.Get("/admin/test/add", func(c *fiber.Ctx) error {

		var key64 = "xMlNA14AlEN+IrZbslsO4hZFd8/rEExx+08X8SVg4rY="
		var data64 = "X631LZTLg5+Qy1wkfRUEz5OtLYh9fbFNFkrqKEFlPAZv6F9Q0WEbs7MJm6TMkomyUFcJo5YdTSqzyd8YmI8jqInPhkuOPRhKElkmrwhSp+s="
		var iv = "MjAyNC0wMi0wOCAw"

		key, err := base64.StdEncoding.DecodeString(key64)
		if err != nil {
			panic(err)
		}

		data, err := base64.StdEncoding.DecodeString(data64)
		if err != nil {
			panic(err)
		}

		result, err := utils.DecryptCBC(key, []byte(iv), data)
		if err != nil {
			panic(err)
		}
		resultString, err := base64.StdEncoding.DecodeString(string(result))
		if err != nil {
			panic(err)
		}
		fmt.Println(string(resultString))

		return c.SendString(string(resultString))
	})
}
