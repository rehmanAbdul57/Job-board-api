package controllers

import (
	"example.com/job-board/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func CreateJob(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	employerID := c.Locals("user_id")

	if role != "employer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only employers can post jobs",
		})
	}

	var job models.Job
	if err := c.BodyParser(&job); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "couldn't parse request body",
		})
	}

	job.EmployerID = employerID.(int)

	err := job.Save()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "couldn't save job",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "job created successfully",
		"job": fiber.Map{
			"id":           job.ID,
			"title":        job.Title,
			"description":  job.Description,
			"location":     job.Location,
			"company_name": job.CompanyName,
			"employer_id":  job.EmployerID,
		},
	})
}

func GetJobs(c *fiber.Ctx) error {
	jobs, err := models.GetAllJobs()
	if err != nil {
		fmt.Println("DEBUG SQL ERROR:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to fetch jobs ",
		})
	}
	return c.Status(fiber.StatusOK).JSON(jobs)
}
