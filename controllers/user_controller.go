package controllers

import (
	"example.com/job-board/models"
	"example.com/job-board/utils"
	"github.com/gofiber/fiber/v2"
)

func Signup(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"couldn't parse json": err.Error(),
		})
	}

	err := user.Save()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"couldn't save user": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"couldn't parse json": err.Error(),
		})
	}

	err := user.ValidateCredentials()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"couldn't authenticate user": err.Error(),
		})
	}

	token, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"couldn't authenticate user": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Login successfully": token,
	})
}
