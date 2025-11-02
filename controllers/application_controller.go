package controllers

import (
	"example.com/job-board/config"
	"example.com/job-board/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ApplyToJob(c *fiber.Ctx) error {
	role := c.Locals("role")
	seekerID := c.Locals("user_id")

	if role != "job_seeker" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only job seekers can apply to jobs",
		})
	}

	jobID, _ := strconv.Atoi(c.FormValue("job_id"))
	coverLetter := c.FormValue("cover_letter")

	file, err := c.FormFile("resume")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "resume upload failed",
		})
	}

	if !strings.HasSuffix(file.Filename, ".pdf") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only PDF files can be uploaded",
		})
	}

	uploadPath := "./uploads/resume"
	os.MkdirAll(uploadPath, os.ModePerm)

	filename := fmt.Sprintf("%d_%s", seekerID.(int), file.Filename)
	savePath := filepath.Join(uploadPath, filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to save resume",
		})
	}

	app := models.Application{
		JobID:       jobID,
		SeekerID:    seekerID.(int),
		CoverLetter: coverLetter,
		Resume:      filename,
	}

	if err := app.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot save application",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "application successfully saved",
		"application": app,
	})
}

func GetApplicationsForJob(c *fiber.Ctx) error {
	role := c.Locals("role")
	employerID := c.Locals("user_id")
	jobID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "could not parse job id",
		})
	}

	if role != "employer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only employer can view applications",
		})
	}

	var ownerID int

	if ownerID != employerID.(int) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you do not own this job",
		})
	}

	applications, err := models.GetApplication(jobID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot get applications",
		})
	}
	return c.JSON(applications)
}

func GetMyJobs(c *fiber.Ctx) error {
	role := c.Locals("role")
	employerID := c.Locals("user_id")

	if role != "employer" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only employers can view their jobs",
		})

	}

	rows, err := config.DB.Query(`SELECT id, title, description, location, company_name, created_at FROM Jobs WHERE employer_id = ? ORDER BY created_at DESC `, employerID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch employer jobs",
		})
	}

	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job
		err = rows.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.CompanyName, &job.CreatedAt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to parse job data",
			})
		}
		jobs = append(jobs, job)
	}
	return c.JSON(jobs)
}
