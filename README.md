# API Server for askTRACE

![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache_Kafka-231F20.svg?style=for-the-badge&logo=apache-kafka&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED.svg?style=for-the-badge&logo=docker&logoColor=white)
![Semantic Release](https://img.shields.io/badge/Semantic_Release-494949.svg?style=for-the-badge&logo=semantic-release&logoColor=white)
![Make](https://img.shields.io/badge/Make-427819.svg?style=for-the-badge&logo=gnu&logoColor=white)
![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-000000.svg?style=for-the-badge&logo=opentelemetry&logoColor=white)
![Jaeger](https://img.shields.io/badge/Jaeger-66CFE3.svg?style=for-the-badge&logo=jaeger&logoColor=white)
![Jenkins](https://img.shields.io/badge/Jenkins-D24939.svg?style=for-the-badge&logo=jenkins&logoColor=white)

The API Server is a backend application designed to handle API requests, manage data, and serve business logic for the project `askTRACE`. It follows a structured architecture with modular components for maintainability and scalability.

The server provides RESTful endpoints for managing users, instructors, courses, and trace data. It integrates with Kafka for event streaming and implements OpenTelemetry for distributed tracing.

## Folder Structure

```
api-server/
│── cmd/                # Application entry points
│── internal/           # Internal application logic
│   ├── config/         # Configuration settings
│   ├── database/       # Database connection and migrations
│   ├── handlers/       # API request handlers
│   ├── kafka/          # Kafka producer/consumer logic
│   ├── middleware/     # Middleware functions
│   ├── models/         # Data models
│   ├── observability/  # Tracing and monitoring setup
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
│── Jenkinsfile         # CI/CD pipeline configuration for Jenkins
│── Jenkinsfile.commitlint # Linting pipeline for commit messages
│── Jenkinsfile.prcheck # Pipeline for PR validation
│── Makefile            # Build and run automation
│── package.json        # Dependencies
│── README.md           # Project documentation
```

## API Routes

The API Server provides the following endpoints:

### Public Routes
- `GET /healthz` - Health check endpoint
- `POST /v1/user` - Create a new user
- `GET /v1/instructor/{instructor_id}` - Get instructor details
- `GET /v1/course/{course_id}` - Get course details

### Private Routes (Require Authentication)

**User Management:**
- `GET/PUT /v1/user/{user_id}` - Get or update user details

**Instructor Management:**
- `POST /v1/instructor` - Create a new instructor
- `PUT/PATCH/DELETE /v1/instructor/{instructor_id}` - Update or delete instructor
- `GET /v1/instructors` - Get all instructors

**Course Management:**
- `POST /v1/course` - Create a new course
- `PUT/PATCH/DELETE /v1/course/{course_id}` - Update or delete course
- `GET /v1/courses` - Get all courses

**Trace Management:**
- `POST/GET /v1/course/{course_id}/trace` - Create or get traces for a course
- `GET/DELETE /v1/course/{course_id}/trace/{trace_id}` - Get or delete specific trace
- `GET /v1/traces` - Get all traces
- `GET /v1/course/{course_id}/trace/{trace_id}/pdf` - Download trace as PDF

**Reference Data:**
- `GET /v1/departments` - Get all departments
- `GET /v1/semesters` - Get all semester terms

## Environment Variables

Before running the application, ensure you have the required environment variables set:

| Variable         | Description                         | Default                                             |
|------------------|-------------------------------------|-----------------------------------------------------|
| `DB_HOST`        | Database host address               | `""`                                                |
| `DB_PORT`        | Database port                       | `5432`                                              |
| `DB_NAME`        | Database name                       | `""`                                                |
| `DB_USER`        | Database username                   | `""`                                                |
| `DB_PASSWORD`    | Database password                   | `""`                                                |
| `SERVER_PORT`    | Port for API server                 | `8080`                                              |
| `KAFKA_BROKERS`  | Comma-separated Kafka broker list   | `"kafka-controller-0.kafka-controller-headless.kafka.svc.cluster.local:9092,kafka-controller-1.kafka-controller-headless.kafka.svc.cluster.local:9092,kafka-controller-2.kafka-controller-headless.kafka.svc.cluster.local:9092"` |
| `KAFKA_TOPIC`    | Kafka topic for trace events        | `trace-survey-uploaded`                             |
| `KAFKA_USERNAME` | Kafka authentication username       | `""`                                                |
| `KAFKA_PASSWORD` | Kafka authentication password       | `""`                                                |
| `SERVICE_NAME`   | OpenTelemetry service name          | `api-server`                                        |
| `OTLP_ENDPOINT`  | OpenTelemetry collector endpoint    | `localhost:4317`                                    |

You can configure these by exporting them in your terminal before running the application or by using environment files with Docker.

## Instructions To Run

### Prerequisites

Ensure you have the following installed on your system:

- Go (latest version recommended)
- Docker (optional for containerized deployment)
- Make (for simplified build and run commands)
- Kafka (for local development with event streaming)
- PostgreSQL (for local database)

### Running Locally

```sh
git clone https://github.com/cyse7125-sp25-team03/api-server
cd api-server
```

Run `make` on your terminal in the root directory to build the binary of the application. And then run the binary which will be created in `/bin` folder.

Alternatively, run `make run` to run the application without building the binary.

### Running using Docker

To build and run using Docker:

```sh
docker build -t api-server .
docker run -p 8080:8080 --env-file .env api-server
```

## Deployment with Helm

The API Server can be deployed to Kubernetes using the Helm chart available at [helm-charts repository](https://github.com/cyse7125-sp25-team03/helm-charts.git) in the `api-server` folder.

To deploy with Helm:

```sh
helm install api-server ./api-server -n api-server

cd manifests

#clusterissuer.yaml
EMAIL="username@gmail.com" envsubst < clusterissuer.yaml | kubectl apply -f -

#ingress.yaml
HOST="api-server.prd.gcp.csyeteam03.xyz" envsubst < ingress.yaml | kubectl apply -f -

```

## CI/CD and Releases

### CI/CD

This project uses Jenkins for continuous integration. The repository includes:

- `Jenkinsfile`: Main CI pipeline.
- `Jenkinsfile.commitlint`: Linting pipeline for commit messages.
- `Jenkinsfile.prcheck`: Pipeline for PR validation.

The pipeline performs the following steps:
1. Code linting and validation
2. Build and test the application
3. Security scanning
4. Docker image building and tagging
5. Semantic versioning

### Releases

- When a pull request is successfully merged, a Docker image is built.
- The Semantic Versioning bot creates a release on GitHub with a tag.
- The tagged release is used for the Docker image, which is then pushed to Docker Hub.

## Observability

The API Server includes OpenTelemetry integration for distributed tracing. Traces are collected and can be visualized using Jaeger or other compatible tools. This provides insights into request flows, performance bottlenecks, and system behavior.

## Contributing

1. Fork the repository
2. Create a new feature branch (`git checkout -b feature-branch`)
3. Commit your changes (`git commit -m "Add new feature"`)
4. Push to your branch (`git push origin feature-branch`)
5. Open a pull request

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.