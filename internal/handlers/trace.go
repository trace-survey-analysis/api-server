package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"api-server/internal/database"
	"api-server/internal/kafka"
	"api-server/internal/middleware"
	"api-server/internal/models"
	"api-server/internal/repositories"
	"api-server/internal/services"
	"api-server/internal/utils"
	"api-server/internal/validators"

	"github.com/google/uuid"
)

// Extracts traceId from /v1/course/{courseId}/trace/{traceId}
func extractTraceID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) < 6 || parts[5] == "" {
		return ""
	}
	return parts[5]
}

// for endpoint: /v1/course/{courseId}/trace
func TraceHandler(w http.ResponseWriter, r *http.Request) {
	courseID := extractCourseID(r.URL.Path)
	if courseID == "" {
		http.NotFound(w, r)
		return
	}

	if _, err := uuid.Parse(courseID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}
	switch r.Method {
	case http.MethodPost:
		createTraceHandler(w, r, courseID)
	case http.MethodGet: //get ALL
		getAllTraceHandler(w, r, courseID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// for endpoint: /v1/course/{courseId}/trace/{traceId}
func TraceEntityHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: extract courseID and traceID from the URL
	courseID := extractCourseID(r.URL.Path)
	traceID := extractTraceID(r.URL.Path)
	if courseID == "" || traceID == "" {
		http.NotFound(w, r)
		return
	}
	if _, err := uuid.Parse(courseID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}
	if _, err := uuid.Parse(traceID); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}
	switch r.Method {
	case http.MethodGet:
		getTraceByIDHandler(w, r, courseID, traceID)
	case http.MethodDelete:
		deleteTraceHandler(w, r, courseID, traceID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createTraceHandler(w http.ResponseWriter, r *http.Request, courseID string) {
	//checks
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}
	// Parse multipart form to handle file upload
	err := r.ParseMultipartForm(10 << 20) // 10MB limit
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}
	traceReq := models.TraceRequest{
		InstructorID: r.FormValue("instructor_id"),
		SemesterTerm: r.FormValue("semester_term"),
		Section:      r.FormValue("section"),
	}

	log.Printf("Trace Request: %v", traceReq)

	// Retrieve file
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		respondWithError(w, http.StatusBadRequest, "failed to get file from request")
		return
	}
	defer file.Close()
	log.Printf("Received file: %s, size: %d bytes", handler.Filename, handler.Size)

	if err := validators.ValidateFileName(handler.Filename); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	//extract other fields

	// add trace request validation from validators
	if err := validators.ValidateTraceRequest(traceReq); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validators.ValidateCourseID(courseID); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	//get user id
	user := middleware.GetUserFromContext(r)
	if user == nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := user.UserID

	// check for instructorid, courseid, semesterterm existence
	if _, err := repositories.GetInstructorByID(database.GetDB(), traceReq.InstructorID); err != nil {
		log.Printf("Error fetching instructor: %v", err)
		http.Error(w, "failed to get instructor", http.StatusBadRequest)
		return
	}

	if _, err := repositories.GetCourseByID(database.GetDB(), courseID); err != nil {
		log.Printf("Error fetching course: %v", err)
		http.Error(w, "failed to get course", http.StatusBadRequest)
		return
	}

	if _, err := repositories.GetSemesterTerm(database.GetDB(), traceReq.SemesterTerm); err != nil {
		log.Printf("Error fetching semester term: %v", err)
		http.Error(w, "failed to get semester term", http.StatusBadRequest)
		return
	}

	//upload file to GCS
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("Bucket name is not set in environment variables!")
	}

	uploadedFilePath, err := utils.UploadFileToGCS(file, handler.Filename, bucketName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to upload file")
		return
	}

	trace := models.Trace{
		TraceID:      uuid.New().String(),
		UserID:       userID,
		FileName:     handler.Filename,
		DateCreated:  time.Now().UTC(),
		BucketPath:   uploadedFilePath,
		CourseID:     courseID,
		InstructorID: traceReq.InstructorID,
		SemesterTerm: traceReq.SemesterTerm,
		Section:      traceReq.Section,
	}
	log.Printf("Trace: %v", trace)

	//Create trace
	newTrace, err := repositories.CreateTrace(database.GetDB(), trace)

	// Publish to Kafka if a producer is available
	kafkaProducer := services.GetKafkaProducer()
	if kafkaProducer != nil {
		// Extract bucket name and path from GCS URL
		bucketName := utils.ExtractBucketNameFromGCS(trace.BucketPath)
		filePath := utils.ExtractFilePathFromGCS(trace.BucketPath)

		uploadMessage := kafka.TraceUploadMessage{
			TraceID:      trace.TraceID,
			CourseID:     trace.CourseID,
			FileName:     trace.FileName,
			GCSBucket:    bucketName,
			GCSPath:      filePath,
			InstructorID: trace.InstructorID,
			SemesterTerm: trace.SemesterTerm,
			Section:      trace.Section,
			UploadedBy:   trace.UserID,
			UploadedAt:   trace.DateCreated,
		}

		err := kafkaProducer.PublishTraceUpload(r.Context(), uploadMessage)
		if err != nil {
			// Log error but don't fail the response
			log.Printf("Error publishing to Kafka: %v", err)
		} else {
			log.Printf("Successfully published trace %s to Kafka", trace.TraceID)
		}
	}

	if err != nil {
		log.Printf("Error creating trace: %v", err)
		http.Error(w, "failed to create trace", http.StatusInternalServerError)
		return
	}
	// return 201 status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrace)

}

func getAllTraceHandler(w http.ResponseWriter, r *http.Request, courseID string) {
	//checks
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}

	//get all traces by courseID
	traces, err := repositories.GetTraceByCourseID(database.GetDB(), courseID)
	if err != nil {
		log.Printf("Error fetching traces: %v", err)
		http.Error(w, "failed to get traces", http.StatusInternalServerError)
		return
	}

	// return 200 status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(traces)
}

// get trace by traceid
func getTraceByIDHandler(w http.ResponseWriter, r *http.Request, courseID string, traceID string) {
	//checks
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}
	//check if course id is valid
	if _, err := repositories.GetCourseByID(database.GetDB(), courseID); err != nil {
		log.Printf("Error fetching course: %v", err)
		http.Error(w, "failed to get course", http.StatusBadRequest)
		return
	}

	//get trace by traceID
	trace, err := repositories.GetTraceByID(database.GetDB(), traceID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "trace not found", http.StatusNotFound)
			return
		}
		log.Printf("Error fetching trace: %v", err)
		http.Error(w, "failed to get trace", http.StatusInternalServerError)
		return
	}

	// return 200 status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trace)
}

// delete trace by traceid
func deleteTraceHandler(w http.ResponseWriter, r *http.Request, courseID string, traceID string) {
	//checks
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if len(r.URL.Query()) > 0 {
		respondWithError(w, http.StatusBadRequest, "query parameters are not allowed")
		return
	}
	//check if course id is valid
	if _, err := repositories.GetCourseByID(database.GetDB(), courseID); err != nil {
		log.Printf("Error fetching course: %v", err)
		http.Error(w, "failed to get course", http.StatusBadRequest)
		return
	}
	//get filepath from trace id
	filePath, err := repositories.GetFilePath(database.GetDB(), traceID)
	if err != nil {
		log.Printf("Error fetching file path: %v", err)
		http.Error(w, "failed to get file path", http.StatusInternalServerError)
		return
	}
	//delete file from GCS
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("Bucket name is not set in environment variables!")
	}
	err = utils.DeleteFileFromGCS(filePath, bucketName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete file")
		return
	}
	//delete trace by traceID
	errDelete := repositories.DeleteTrace(database.GetDB(), traceID)
	if errDelete != nil {
		if errors.Is(errDelete, sql.ErrNoRows) {
			http.Error(w, "trace not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting trace from database: %v", errDelete)
		http.Error(w, "failed to delete trace from database", http.StatusInternalServerError)
		return
	}

	// return 204 status code
	w.WriteHeader(http.StatusNoContent)
}
