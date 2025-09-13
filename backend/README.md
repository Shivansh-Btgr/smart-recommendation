# Internship Recommendation System Backend

This backend provides user authentication, profile management, internship management and recommendation using external ML Service

Built with Go, Gin, GORM, and JWT.

## Project Structure

- `main.go` — Entry point, sets up routes and database
- `models/` — Data models (User, Profile)
- `handlers/` — HTTP handlers for auth and profile
- `utils/` — Utility functions (password hashing, JWT)

## Data Models

### User
```go
{
  "id": 1,
  "email": "user@example.com",
  "passwordHash": "...",
  "isProfileComplete": false,
  "profile": { ...Profile fields... }
}
```

### Profile
```go
{
  "id": 1,
  "userId": 1,
  "name": "John Doe",
  "phone": "1234567890",
  "education": "B.Tech Computer Science",
  "cgpa": 8.5,
  "skills": ["Go", "Python", "Machine Learning"],
  "experience": 2,
  "socialLinks": ["https://linkedin.com/in/johndoe"],
  "location": "Delhi, India",
  "interest": "Backend Development",
  "resumeLink": "https://...",
  "preferredJobType": "Remote",
  "availability": "Immediate",
  "languages": ["English", "Hindi"]
}
```

### Internship
```json
{
  "id": 1,
  "title": "Backend Intern",
  "organization": "Tech Corp",
  "location": "Remote",
  "stipend_inr": 15000,
  "duration": 6,
  "skillsreq": ["Go", "SQL"],
  "langsreq": ["English"],
  "active": true,
  "posted_at": "2025-09-01T12:00:00Z",
  "deadline": "2025-09-30T23:59:59Z",
  "apply_url": "https://apply.here",
  "description": "Work on backend APIs.",
  "min_cgpa": 7.0,
  "experience": 0
}
```

## API Endpoints

### Register
- `POST /register`
- Request JSON:
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```
- Response: `{ "message": "Registration successful" }`

### Login
- `POST /login`
- Request JSON:
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```
- Response: `{ "token": "<JWT token>" }`

### Get Profile
- `GET /profile` (Requires Bearer JWT)
- Response: Profile JSON (see above)

### Update Profile
- `PUT /profile` (Requires Bearer JWT)
- Request JSON:
```json
{
  "name": "John Doe",
  "phone": "1234567890",
  "education": "B.Tech Computer Science",
  "cgpa": 8.5,
  "skills": ["Go", "Python"],
  "experience": 2,
  "social_links": ["https://linkedin.com/in/johndoe"],
  "location": "Delhi, India",
  "interest": "Backend Development",
  "resume_link": "https://...",
  "preferred_job_type": "Remote",
  "availability": "Immediate",
  "languages": ["English", "Hindi"]
}
```
- Response: `{ "message": "Profile updated" }`

### Create Internship
- `POST /internships` (Requires Bearer JWT)
- Request JSON: (see Internship model above, omit id/posted_at)
- Response: Internship JSON

### Example: Create Internship Request
```json
{
  "title": "Backend Intern",
  "organization": "Tech Corp",
  "location": "Remote",
  "stipend_inr": 15000,
  "duration": 6,
  "skillsreq": ["Go", "SQL"],
  "langsreq": ["English"],
  "active": true,
  "deadline": "2025-09-30T23:59:59Z",
  "apply_url": "https://apply.here",
  "description": "Work on backend APIs.",
  "min_cgpa": 7.0,
  "experience": 0
}
```

### Edit Internship
- `PUT /internships/:id` (Requires Bearer JWT)
- Request JSON: (same as create)
- Response: Internship JSON

### Get All Internships
- `GET /internships`
- Query params:
  - `location` (string, filter by location)
  - `min_stipend` (int, filter by minimum stipend)
  - `min_duration` (int, filter by minimum duration in months)
  - `active` (true/false, filter by active status)
  - `sort` ("stipend" or "deadline", default: posted_at desc)
- Example:
```
GET /internships?location=Remote&min_stipend=10000&sort=stipend
```
- Response: Array of internships

### Get Internship by ID
- `GET /internships/:id`
- Response: Internship JSON

### Get Active Internships
- `GET /internships/active`
- Response: Array of active internships

### Get Recommendations
- `GET /recommendations` (Requires Bearer JWT, profile must be complete)
- Response:
```json
{
  "recommendation_ids": [2, 5, 7, 9, 12],
  "recommendations": [
    { "id": 2, "title": "Backend Intern", ... },
    { "id": 5, "title": "ML Intern", ... }
  ]
}
```

## Docker & Docker Compose

### Build and Run with Docker Compose

1. Build and start the backend and Postgres database:
   ```sh
   docker compose up --build
   ```
   This will start both the backend (on port 8080) and Postgres (on port 5432).

2. The backend will auto-migrate the database tables on startup.

3. To stop and remove containers:
   ```sh
   docker compose down
   ```

### Environment Variables
- The backend reads DB and ML service config from the `.env` file or from environment variables set in `docker-compose.yml`.
- Do not commit `.env` to version control (see `.gitignore`).

### Notes
- You can connect to the Postgres database using any client at `localhost:5432` with the credentials in `docker-compose.yml`.
- The backend service will wait for the database to be ready before starting.
- For production, update passwords and secrets in `.env` and `docker-compose.yml`.

## Notes
- All protected endpoints require an `Authorization: Bearer <token>` header.
- Passwords are securely hashed.
- Profile must be completed before accessing recommendations.
- Internships can be filtered and sorted using query parameters.

---
