package models

import (
	"errors"
	"example.com/job-board/config"
	"example.com/job-board/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"-" binding:"required"`
	Role     string `json:"role"`
}

func (u User) Save() error {
	query := `INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)`
	stmt, err := config.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPswd, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(u.Name, u.Email, hashedPswd, u.Role)
	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, name, email, password, role FROM users WHERE email = ?`
	row := config.DB.QueryRow(query, u.Email)

	var retrievedPswd string
	err := row.Scan(&u.ID, &u.Name, &u.Email, &retrievedPswd, &u.Role)
	if err != nil {
		return err
	}

	pswdIsValid := utils.CheckHashPassword(u.Password, retrievedPswd)
	if !pswdIsValid {
		return errors.New("Invalid Credentials")
	}
	return nil
}
