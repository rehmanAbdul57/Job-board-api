package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := "root:Dbmysql57?@tcp(127.0.0.1:3306)/job_board?parseTime=true"
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to DB")

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'employer', 'job_seeker'),
    created_at DateTIME DEFAULT CURRENT_TIMESTAMP
)`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		panic(fmt.Sprintf("could not create users table: %v", err))
	}

	createJobsTable := `CREATE TABLE IF NOT EXISTS jobs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(100),
    company_name VARCHAR(100),
    employer_id INT,
    created_at DateTime DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employer_id) REFERENCES users(id)
)`

	_, err = DB.Exec(createJobsTable)
	if err != nil {
		panic(fmt.Sprintf("could not create jobs table: %v", err))
	}

	createApplicationsTable := `CREATE TABLE IF NOT EXISTS applications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    job_id INT NOT NULL,
    seeker_id INT NOT NULL,
    cover_letter TEXT,
    resume VARCHAR(255),
    created_at DateTime DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (job_id) REFERENCES jobs(id),
    FOREIGN KEY (seeker_id) REFERENCES users(id)
)`

	_, err = DB.Exec(createApplicationsTable)
	if err != nil {
		panic(fmt.Sprintf("could not create applications table: %v", err))
	}

	fmt.Println("Tables created successfully")
}
