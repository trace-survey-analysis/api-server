## Table of Contents
- [api-server](#api-server)
- [Folder Structure](#folder-structure)
- [Environment Variables](#environment-variables)
- [Instructions To Run](#instructions-to-run)
  - [Prerequisites](#prerequisites)
  - [Running Locally](#running-locally)
  - [Running using Docker](#running-using-docker)
- [CI/CD and Releases](#cicd-and-releases)
  - [CI/CD](#cicd)
  - [Releases](#releases)
- [Contributing](#contributing)

# api-server
The api-server is a backend application designed to handle API requests, manage data, and serve business logic for the project `trace-survey-analysis`. It follows a structured architecture with modular components for maintainability and scalability.

# Folder Structure
```
api-server/
│── cmd/                # Application entry points
│── internal/           # Internal application logic
│   ├── config/         # Configuration settings
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # API request handlers
│   ├── middleware/     # Middleware functions
│   ├── models/         # Data models
│   ├── repositories/   # Data access layer
│   ├── routes/         # API route definitions
│   ├── services/       # Business logic
│   ├── utils/          # Utility functions
│   ├── validators/     # Input validation logic
│── .gitignore          # Files and folders to ignore in Git
│── .releaserc.json     # Configuration for release management
│── Dockerfile          # Docker configuration for containerization
│── go.mod              # Go module dependencies
│── go.sum              # Dependency checksums
│── Jenkinsfile*        # CI/CD pipeline configurations for Jenkins
│── Makefile            # Build and run automation
│── package.json        # Dependencies
│── README.md           # Project documentation

```

# Environment Variables
Environment Variables
Before running the application, ensure you have the required environment variables set:

| Variable      | Description                            | Default  |
|--------------|----------------------------------------|----------|
| `DB_HOST`    | Database host address                 | `""`     |
| `DB_PORT`    | Database port                         | `5432`   |
| `DB_NAME`    | Database name                         | `""`     |
| `DB_USER`    | Database username                     | `""`     |
| `DB_PASSWORD` | Database password                   | `""`     |
| `SERVER_PORT` | Port for API server                  | `8080`   |
| `BUCKET_NAME` | Google Cloud Storage bucket name    | `""`     |

You can configure these by exporting them in your terminal before running the application.

# Instructions To Run
## Prerequisites
  Ensure you have the following installed on your system:

- Go (latest version recommended)
- Docker (optional for containerized deployment)
- Make (for simplified build and run commands)

## Running Locally
```
git clone https://github.com/cyse7125-sp25-team03/api-server
cd api-server
```
Run ```make``` on your terminal in the root directory to build the binary of the application. And then run the binary which will be created in `/bin` folder

Alternatively, run ```make run``` to run the application without building the binary.

## Running using Docker
To build and run using Docker:
```
docker build -t api-server .
docker run -p 8080:8080 api-server
```

# CI/CD and Releases

## CI/CD
This project uses Jenkins for continuous integration. The repository includes:

- Jenkinsfile: Main CI pipeline.
- Jenkinsfile.commitlint: Linting pipeline for commit messages.
- Jenkinsfile.prcheck: Pipeline for PR validation.

## Releases
- When a pull request is successfully merged, a Docker image is built.
- The Semantic Versioning bot creates a release on GitHub with a tag.
- The tagged release is used for the Docker image, which is then pushed to Docker Hub.

# Contributing
1. Fork the repository
2. Create a new feature branch (git checkout -b feature-branch)
3. Commit your changes (git commit -m "Add new feature")
4. Push to your branch (git push origin feature-branch)
5. Open a pull request