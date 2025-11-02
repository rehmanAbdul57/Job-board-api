package models

import (
	"example.com/job-board/config"
	"fmt"
	"time"
)

type Application struct {
	ID          int       `json:"id"`
	JobID       int       `json:"job_id"`
	SeekerID    int       `json:"seeker_id"`
	CoverLetter string    `json:"cover_letter"`
	Resume      string    `json:"resume"`
	CreatedAt   time.Time `json:"created_at"`
}

func (app Application) Save() error {
	query := `INSERT INTO applications (job_id, seeker_id, cover_letter, resume)
VALUES (?, ?, ?, ?)`

	stmt, err := config.DB.Prepare(query)
	if err != nil {

		fmt.Println("save error:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.JobID, app.SeekerID, app.CoverLetter, app.Resume)
	if err != nil {
		fmt.Println("save 2 error:", err)
		return err
	}
	return nil
}

func GetApplication(id int64) ([]Application, error) {
	var ownerID int
	err := config.DB.QueryRow("SELECT * (employer_id) FROM jobs WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		return nil, err
	}

	rows, err := config.DB.Query(`SELECT id, job_id, seeker_id, cover_letter, resume, u.email, u.name, created_at
FROM applications AS a
JOIN users AS u ON a.seeker_id = u.id
WHERE a.job_id = ?
ORDER BY a.created_at DESC 
`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application
	var u User

	for rows.Next() {
		var application Application
		err := rows.Scan(&application.ID, &application.JobID, &application.SeekerID, &application.CoverLetter, &application.Resume, &u.Email, &u.Name, &application.CreatedAt)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}
	return applications, nil

}
