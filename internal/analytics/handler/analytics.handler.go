package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Dom-HTG/attendance-management-system/internal/analytics/service"
	"github.com/Dom-HTG/attendance-management-system/pkg/middleware"
	"github.com/Dom-HTG/attendance-management-system/pkg/responses"
	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics endpoints
type AnalyticsHandler struct {
	service service.AnalyticsServiceInterface
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(svc service.AnalyticsServiceInterface) *AnalyticsHandler {
	return &AnalyticsHandler{service: svc}
}

// ===== Student Analytics Endpoints =====

// GetStudentMetrics handles GET /api/analytics/student/{student_id}
func (ah *AnalyticsHandler) GetStudentMetrics(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("student_id"))
	if err != nil {
		responses.ApiFailure(ctx, "Invalid student ID", http.StatusBadRequest, err)
		return
	}

	// Get current user's ID for authorization
	currentUserID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	// Students can only view their own metrics; lecturers/admins can view any
	userRole, _ := middleware.GetUserRoleFromContext(ctx)
	if userRole == "student" && currentUserID != studentID {
		responses.ApiFailure(ctx, "Unauthorized: Students can only view their own metrics", http.StatusForbidden, nil)
		return
	}

	metrics, err := ah.service.GetStudentMetrics(studentID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve student metrics", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Student metrics retrieved successfully", metrics)
}

// GetStudentInsights handles GET /api/analytics/student/{student_id}/insights
func (ah *AnalyticsHandler) GetStudentInsights(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("student_id"))
	if err != nil {
		responses.ApiFailure(ctx, "Invalid student ID", http.StatusBadRequest, err)
		return
	}

	insights, err := ah.service.GetStudentInsights(studentID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate student insights", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Student insights generated successfully", insights)
}

// ===== Lecturer Analytics Endpoints =====

// GetLecturerCourseMetrics handles GET /api/analytics/lecturer/courses
func (ah *AnalyticsHandler) GetLecturerCourseMetrics(ctx *gin.Context) {
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	metrics, err := ah.service.GetLecturerCourseMetrics(lecturerID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve lecturer metrics", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Lecturer course metrics retrieved successfully", metrics)
}

// GetLecturerCoursePerformance handles GET /api/analytics/lecturer/course/{course_code}
func (ah *AnalyticsHandler) GetLecturerCoursePerformance(ctx *gin.Context) {
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	courseCode := ctx.Param("course_code")
	if courseCode == "" {
		responses.ApiFailure(ctx, "Course code is required", http.StatusBadRequest, nil)
		return
	}

	performance, err := ah.service.GetLecturerCoursePerformance(lecturerID, courseCode)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve course performance", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Course performance retrieved successfully", performance)
}

// GetLecturerInsights handles GET /api/analytics/lecturer/insights
func (ah *AnalyticsHandler) GetLecturerInsights(ctx *gin.Context) {
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	insights, err := ah.service.GetLecturerInsights(lecturerID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate insights", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Lecturer insights generated successfully", insights)
}

// ===== Admin Analytics Endpoints =====

// GetAdminOverview handles GET /api/analytics/admin/overview
func (ah *AnalyticsHandler) GetAdminOverview(ctx *gin.Context) {
	overview, err := ah.service.GetAdminOverview()
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve admin overview", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Admin overview retrieved successfully", overview)
}

// GetDepartmentMetrics handles GET /api/analytics/admin/department/{department}
func (ah *AnalyticsHandler) GetDepartmentMetrics(ctx *gin.Context) {
	department := ctx.Param("department")
	if department == "" {
		responses.ApiFailure(ctx, "Department name is required", http.StatusBadRequest, nil)
		return
	}

	metrics, err := ah.service.GetDepartmentMetrics(department)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve department metrics", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Department metrics retrieved successfully", metrics)
}

// GetRealTimeDashboard handles GET /api/analytics/admin/realtime
func (ah *AnalyticsHandler) GetRealTimeDashboard(ctx *gin.Context) {
	dashboard, err := ah.service.GetRealTimeDashboard()
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve real-time dashboard", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Real-time dashboard retrieved successfully", dashboard)
}

// ===== Temporal Analytics Endpoint =====

// GetTemporalAnalytics handles GET /api/analytics/temporal
func (ah *AnalyticsHandler) GetTemporalAnalytics(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")
	granularity := ctx.Query("granularity") // daily, weekly, monthly

	if startDateStr == "" || endDateStr == "" {
		responses.ApiFailure(ctx, "start_date and end_date query parameters are required (RFC3339 format)", http.StatusBadRequest, nil)
		return
	}

	if granularity == "" {
		granularity = "weekly"
	}

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		responses.ApiFailure(ctx, "Invalid start_date format", http.StatusBadRequest, err)
		return
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		responses.ApiFailure(ctx, "Invalid end_date format", http.StatusBadRequest, err)
		return
	}

	temporal, err := ah.service.GetTemporalAnalytics(startDate, endDate, granularity)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve temporal analytics", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Temporal analytics retrieved successfully", temporal)
}

// ===== Anomaly Detection Endpoint =====

// DetectAnomalies handles GET /api/analytics/anomalies
func (ah *AnalyticsHandler) DetectAnomalies(ctx *gin.Context) {
	anomalies, err := ah.service.DetectAnomalies()
	if err != nil {
		responses.ApiFailure(ctx, "Failed to detect anomalies", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Anomaly detection completed", anomalies)
}

// ===== Predictive Analytics Endpoints =====

// PredictStudentAttendance handles GET /api/analytics/predictions/student/{student_id}
func (ah *AnalyticsHandler) PredictStudentAttendance(ctx *gin.Context) {
	studentID, err := strconv.Atoi(ctx.Param("student_id"))
	if err != nil {
		responses.ApiFailure(ctx, "Invalid student ID", http.StatusBadRequest, err)
		return
	}

	prediction, err := ah.service.PredictStudentAttendance(studentID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate prediction", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Attendance prediction generated successfully", prediction)
}

// PredictCourseAttendance handles GET /api/analytics/predictions/course/{course_code}
func (ah *AnalyticsHandler) PredictCourseAttendance(ctx *gin.Context) {
	courseCode := ctx.Param("course_code")
	if courseCode == "" {
		responses.ApiFailure(ctx, "Course code is required", http.StatusBadRequest, nil)
		return
	}

	prediction, err := ah.service.PredictCourseAttendance(courseCode)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate prediction", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Course attendance prediction generated successfully", prediction)
}

// ===== Benchmarking Endpoint =====

// GetBenchmarkComparison handles GET /api/analytics/benchmark
func (ah *AnalyticsHandler) GetBenchmarkComparison(ctx *gin.Context) {
	entityType := ctx.Query("entity_type") // student, course, department
	entityIDStr := ctx.Query("entity_id")

	if entityType == "" || entityIDStr == "" {
		responses.ApiFailure(ctx, "entity_type and entity_id query parameters are required", http.StatusBadRequest, nil)
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		responses.ApiFailure(ctx, "Invalid entity_id", http.StatusBadRequest, err)
		return
	}

	comparison, err := ah.service.GetBenchmarkComparison(entityType, entityID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve benchmark comparison", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Benchmark comparison retrieved successfully", comparison)
}

// ===== Chart Data Endpoints =====

// GetChartData handles GET /api/analytics/charts/{chart_type}
func (ah *AnalyticsHandler) GetChartData(ctx *gin.Context) {
	chartType := ctx.Param("chart_type")
	entityType := ctx.Query("entity_type") // student, course, department
	entityIDStr := ctx.Query("entity_id")

	if chartType == "" || entityType == "" || entityIDStr == "" {
		responses.ApiFailure(ctx, "chart_type, entity_type, and entity_id are required", http.StatusBadRequest, nil)
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		responses.ApiFailure(ctx, "Invalid entity_id", http.StatusBadRequest, err)
		return
	}

	chartData, err := ah.service.GetChartData(chartType, entityType, entityID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve chart data", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Chart data retrieved successfully", chartData)
}
