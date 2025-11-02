# Job Board API (Go + Fiber + MySQL)

A simple Job Board REST API built with **Go (Fiber)** and **MySQL**.  
It allows users to sign up, log in, post jobs, apply to jobs, and view applications — all with JWT authentication.

---

## 🚀 Features

- 👤 **User Roles:** `admin`, `employer`, `job_seeker`
- 🔑 **JWT Authentication**
- 💼 **Employers** can:
  - Create and manage job postings
  - View applications for their jobs
- 🧑‍💻 **Job Seekers** can:
  - View all public jobs
  - Apply to jobs with resume upload (PDF only)
- 📄 **Resume Uploads:** Stored locally in `./uploads/resume/`
