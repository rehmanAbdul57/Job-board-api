package models

import (
	"example.com/job-board/config"
	"fmt"
	"time"
)

type Job struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	CompanyName string    `json:"company_name" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	EmployerID  int       `json:"employer_id"`
}

func (j *Job) Save() error {

	query := `INSERT INTO jobs (title, description, location, company_name, employer_id)
VALUES (?, ?, ?, ?, ?)`

	stmt, err := config.DB.Prepare(query)
	if err != nil {
		fmt.Println("DB ERROR:", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(j.Title, j.Description, j.Location, j.CompanyName, j.EmployerID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	j.ID = id
	return err
}

func GetAllJobs() ([]Job, error) {
	query := `SELECT id, title, description, location, company_name, employer_id, created_at FROM jobs`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Location, &job.CompanyName, &job.EmployerID, &job.CreatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
