package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	domain "github.com/Dom-HTG/attendance-management-system/internal/analytics/domain"
	"gorm.io/gorm"
)

// AnalyticsRepoInterface defines analytics repository operations
type AnalyticsRepoInterface interface {
	// Student analytics
	GetStudentMetrics(studentID int) (*domain.StudentMetricsResponse, error)
	GetStudentPerCourseRates(studentID int) ([]domain.CourseAttendanceRate, error)
	GetStudentAttendanceTrend(studentID int, startDate, endDate time.Time) ([]domain.TrendDataPoint, error)
	GetStudentEngagementScore(studentID int) (float64, error)
	IsStudentAtRisk(studentID int, threshold float64) (bool, error)

	// Lecturer analytics
	GetLecturerCourseMetrics(lecturerID int) (*domain.LecturerCourseMetricsResponse, error)
	GetLecturerCoursePerformance(lecturerID int, courseCode string) (*domain.CoursePerformanceResponse, error)

	// Admin analytics
	GetAdminOverview() (*domain.AdminOverviewResponse, error)
	GetDepartmentMetrics(department string) (*domain.DepartmentDeepDiveResponse, error)
	GetRealTimeDashboard() (*domain.RealTimeDashboardResponse, error)

	// Temporal analytics
	GetTemporalAnalytics(startDate, endDate time.Time, granularity string) (*domain.TemporalAnalyticsResponse, error)

	// Anomalies
	DetectAnomalies() (*domain.AnomalyResponse, error)
	GetAnomaliesByStudent(studentID int) ([]domain.Anomaly, error)

	// Predictions
	PredictStudentAttendance(studentID int) (*domain.PredictionResponse, error)
	PredictCourseAttendance(courseCode string) (*domain.PredictionResponse, error)

	// Benchmarking
	GetBenchmarkComparison(entityType string, entityID int) (*domain.BenchmarkResponse, error)

	// Utility methods
	GetAttendanceRateForEntity(entityType string, entityID int, startDate, endDate time.Time) (float64, error)
	GetLateCheckInCount(studentID int, startDate, endDate time.Time) (int, error)
	GetAttendanceStreak(studentID int) (int, error)
}

// AnalyticsRepo implements AnalyticsRepoInterface
type AnalyticsRepo struct {
	db *gorm.DB
}

// NewAnalyticsRepo creates a new analytics repository
func NewAnalyticsRepo(db *gorm.DB) AnalyticsRepoInterface {
	return &AnalyticsRepo{db: db}
}

// ===== Student Analytics =====

// GetStudentMetrics returns comprehensive metrics for a student
func (ar *AnalyticsRepo) GetStudentMetrics(studentID int) (*domain.StudentMetricsResponse, error) {
	var response domain.StudentMetricsResponse

	// Get student info
	var student entities.Student
	if err := ar.db.First(&student, studentID).Error; err != nil {
		return nil, errors.New("student not found")
	}

	// Get overall attendance rate
	var result struct {
		TotalSessions int
		TotalPresent  int
	}

	query := `
		SELECT COUNT(ua.id) as total_sessions,
		       SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) as total_present
		FROM user_attendances ua
		WHERE ua.student_id = ?
	`
	if err := ar.db.Raw(query, studentID).Scan(&result).Error; err != nil {
		return nil, err
	}

	response.StudentID = studentID
	response.StudentName = student.FirstName + " " + student.LastName
	response.MatricNumber = student.MatricNumber
	response.TotalSessions = result.TotalSessions
	response.TotalPresent = result.TotalPresent
	response.TotalAbsent = result.TotalSessions - result.TotalPresent

	if result.TotalSessions > 0 {
		response.OverallAttendanceRate = float64(result.TotalPresent) / float64(result.TotalSessions) * 100
	}

	// Get per-course rates
	perCourseRates, err := ar.GetStudentPerCourseRates(studentID)
	if err == nil {
		response.PerCourseRates = perCourseRates
	}

	// Get attendance trend
	trend, err := ar.GetStudentAttendanceTrend(studentID, time.Now().AddDate(0, -3, 0), time.Now())
	if err == nil {
		response.AttendanceTrend = trend
	}

	// Get engagement score
	engScore, err := ar.GetStudentEngagementScore(studentID)
	if err == nil {
		response.EngagementScore = engScore
	}

	// Get late check-ins
	lateCount, err := ar.GetLateCheckInCount(studentID, time.Time{}, time.Now())
	if err == nil {
		response.TotalLate = lateCount
		response.LateCheckInFrequency = lateCount
	}

	// Get attendance streak
	streak, err := ar.GetAttendanceStreak(studentID)
	if err == nil {
		response.AttendanceStreak = streak
	}

	// Check if at risk
	atRisk, err := ar.IsStudentAtRisk(studentID, 75)
	if err == nil {
		response.AtRiskStatus = atRisk
	}

	response.GeneratedAt = time.Now()
	return &response, nil
}

// GetStudentPerCourseRates returns attendance rate per course
func (ar *AnalyticsRepo) GetStudentPerCourseRates(studentID int) ([]domain.CourseAttendanceRate, error) {
	var rates []domain.CourseAttendanceRate

	query := `
		SELECT DISTINCT
			COALESCE(e.event_name, 'Unknown') as course_code,
			COALESCE(e.event_name, 'Unknown') as course_name,
			COUNT(ua.id) as total_sessions,
			SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) as sessions_attended,
			ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / COUNT(ua.id), 2) as attendance_rate
		FROM user_attendances ua
		JOIN events e ON ua.attendance_id = e.id
		WHERE ua.student_id = ?
		GROUP BY e.id, e.event_name
		ORDER BY attendance_rate DESC
	`

	if err := ar.db.Raw(query, studentID).Scan(&rates).Error; err != nil {
		return nil, err
	}

	return rates, nil
}

// GetStudentAttendanceTrend returns attendance trend over time
func (ar *AnalyticsRepo) GetStudentAttendanceTrend(studentID int, startDate, endDate time.Time) ([]domain.TrendDataPoint, error) {
	var trends []domain.TrendDataPoint

	query := `
		SELECT 
			to_char(ua.marked_time, 'YYYY-IW') as period,
			COUNT(ua.id) as total_sessions,
			SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) as sessions_attended,
			ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / COUNT(ua.id), 2) as attendance_rate,
			ROUND(CAST(EXTRACT(EPOCH FROM AVG(ua.marked_time - e.start_time)) AS INT) / 60) as average_checkin_time
		FROM user_attendances ua
		JOIN events e ON ua.attendance_id = e.id
		WHERE ua.student_id = ? AND ua.marked_time >= ? AND ua.marked_time <= ?
		GROUP BY period
		ORDER BY period
	`

	if err := ar.db.Raw(query, studentID, startDate, endDate).Scan(&trends).Error; err != nil {
		return nil, err
	}

	return trends, nil
}

// GetStudentEngagementScore calculates engagement score (0-100)
func (ar *AnalyticsRepo) GetStudentEngagementScore(studentID int) (float64, error) {
	var result struct {
		Score float64
	}

	// Engagement is based on consistency + punctuality
	query := `
		SELECT 
			ROUND((attendance_rate * 0.7) + (punctuality_score * 0.3), 2) as score
		FROM (
			SELECT
				COALESCE(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 0) as attendance_rate,
				COALESCE(100 - (COUNT(CASE WHEN (ua.marked_time - e.start_time) > INTERVAL '5 minutes' THEN 1 END) * 100 / NULLIF(COUNT(ua.id), 0)), 100) as punctuality_score
			FROM user_attendances ua
			LEFT JOIN events e ON ua.attendance_id = e.id
			WHERE ua.student_id = ?
		) sub
	`

	if err := ar.db.Raw(query, studentID).Scan(&result).Error; err != nil {
		return 0, err
	}

	return result.Score, nil
}

// IsStudentAtRisk checks if student is below attendance threshold
func (ar *AnalyticsRepo) IsStudentAtRisk(studentID int, threshold float64) (bool, error) {
	var result struct {
		AttendanceRate float64
	}

	query := `
		SELECT COALESCE(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 0) as attendance_rate
		FROM user_attendances ua
		WHERE ua.student_id = ?
	`

	if err := ar.db.Raw(query, studentID).Scan(&result).Error; err != nil {
		return false, err
	}

	return result.AttendanceRate < threshold, nil
}

// ===== Lecturer Analytics =====

// GetLecturerCourseMetrics returns all course metrics for a lecturer
func (ar *AnalyticsRepo) GetLecturerCourseMetrics(lecturerID int) (*domain.LecturerCourseMetricsResponse, error) {
	var response domain.LecturerCourseMetricsResponse

	var lecturer entities.Lecturer
	if err := ar.db.First(&lecturer, lecturerID).Error; err != nil {
		return nil, errors.New("lecturer not found")
	}

	response.LecturerID = lecturerID
	response.LecturerName = lecturer.FirstName + " " + lecturer.LastName
	response.Department = lecturer.Department

	// Get all courses for lecturer (inferred from events they created)
	query := `
		SELECT COUNT(DISTINCT e.id) as total_courses
		FROM events e
		WHERE e.created_by = ? OR e.id IN (
			SELECT DISTINCT e2.id FROM events e2
		)
	`

	var courseCount int
	ar.db.Raw(query, lecturerID).Scan(&courseCount)
	response.TotalCourses = courseCount

	response.GeneratedAt = time.Now()
	return &response, nil
}

// GetLecturerCoursePerformance returns detailed performance for a course
func (ar *AnalyticsRepo) GetLecturerCoursePerformance(lecturerID int, courseCode string) (*domain.CoursePerformanceResponse, error) {
	var response domain.CoursePerformanceResponse

	query := `
		SELECT 
			COUNT(DISTINCT e.id) as session_count,
			COUNT(DISTINCT ua.student_id) as student_count,
			ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 2) as overall_rate
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.attendance_id
		WHERE e.event_name LIKE ?
	`

	if err := ar.db.Raw(query, "%"+courseCode+"%").Scan(&response).Error; err != nil {
		return nil, err
	}

	response.CourseCode = courseCode
	response.GeneratedAt = time.Now()
	return &response, nil
}

// ===== Admin Analytics =====

// GetAdminOverview returns university-wide metrics
func (ar *AnalyticsRepo) GetAdminOverview() (*domain.AdminOverviewResponse, error) {
	var response domain.AdminOverviewResponse

	// Overall attendance rate
	query := `
		SELECT ROUND(CAST(SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(*), 0), 2)
		FROM user_attendances
	`
	ar.db.Raw(query).Scan(&response.OverallAttendanceRate)

	// Active sessions
	query = `
		SELECT COUNT(*) FROM events WHERE start_time <= NOW() AND end_time >= NOW()
	`
	ar.db.Raw(query).Scan(&response.TotalActiveSessions)

	// Total students and lecturers
	var studentCount, lecturerCount int64
	ar.db.Model(&entities.Student{}).Count(&studentCount)
	ar.db.Model(&entities.Lecturer{}).Count(&lecturerCount)
	response.TotalStudents = int(studentCount)
	response.TotalLecturers = int(lecturerCount)

	response.GeneratedAt = time.Now()
	return &response, nil
}

// GetDepartmentMetrics returns metrics for a specific department
func (ar *AnalyticsRepo) GetDepartmentMetrics(department string) (*domain.DepartmentDeepDiveResponse, error) {
	var response domain.DepartmentDeepDiveResponse

	response.DepartmentName = department
	response.GeneratedAt = time.Now()

	// Get department-level attendance rate
	query := `
		SELECT ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 2)
		FROM user_attendances ua
		JOIN students s ON ua.student_id = s.id
		WHERE s.department = ?
	`
	ar.db.Raw(query, department).Scan(&response.OverallAttendanceRate)

	// Count students, lecturers, courses
	var studentCount, lecturerCount int64
	ar.db.Model(&entities.Student{}).Where("department = ?", department).Count(&studentCount)
	ar.db.Model(&entities.Lecturer{}).Where("department = ?", department).Count(&lecturerCount)
	response.StudentCount = int(studentCount)
	response.LecturerCount = int(lecturerCount)

	return &response, nil
}

// GetRealTimeDashboard returns live dashboard data
func (ar *AnalyticsRepo) GetRealTimeDashboard() (*domain.RealTimeDashboardResponse, error) {
	var response domain.RealTimeDashboardResponse

	// Active sessions right now
	query := `
		SELECT COUNT(*) FROM events WHERE start_time <= NOW() AND end_time >= NOW()
	`
	ar.db.Raw(query).Scan(&response.ActiveSessionsNow)

	// Total check-ins today
	query = `
		SELECT COUNT(*) FROM user_attendances WHERE DATE(marked_time) = CURRENT_DATE
	`
	ar.db.Raw(query).Scan(&response.TotalCheckInsToday)

	// Average attendance today
	query = `
		SELECT ROUND(CAST(SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(*), 0), 2)
		FROM user_attendances
		WHERE DATE(marked_time) = CURRENT_DATE
	`
	ar.db.Raw(query).Scan(&response.AverageAttendanceToday)

	response.GeneratedAt = time.Now()
	return &response, nil
}

// ===== Temporal Analytics =====

// GetTemporalAnalytics returns time-based patterns
func (ar *AnalyticsRepo) GetTemporalAnalytics(startDate, endDate time.Time, granularity string) (*domain.TemporalAnalyticsResponse, error) {
	var response domain.TemporalAnalyticsResponse

	response.Granularity = granularity
	response.StartDate = startDate
	response.EndDate = endDate
	response.GeneratedAt = time.Now()

	// Get day-of-week analysis
	query := `
		SELECT 
			to_char(ua.marked_time, 'Day') as day_of_week,
			ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 2) as attendance_rate,
			COUNT(ua.id) as session_count
		FROM user_attendances ua
		WHERE ua.marked_time >= ? AND ua.marked_time <= ?
		GROUP BY to_char(ua.marked_time, 'Day')
	`

	var dayMetrics []domain.DayOfWeekMetrics
	if err := ar.db.Raw(query, startDate, endDate).Scan(&dayMetrics).Error; err == nil {
		response.DayOfWeekAnalysis = dayMetrics
	}

	return &response, nil
}

// ===== Anomalies =====

// DetectAnomalies identifies unusual attendance patterns
func (ar *AnalyticsRepo) DetectAnomalies() (*domain.AnomalyResponse, error) {
	var response domain.AnomalyResponse

	// Query for duplicate check-ins (same student, same event, within 1 minute)
	query := `
		SELECT 
			ua1.student_id,
			ua1.attendance_id,
			COUNT(*) as duplicate_count
		FROM user_attendances ua1
		WHERE EXISTS (
			SELECT 1 FROM user_attendances ua2
			WHERE ua2.student_id = ua1.student_id
			AND ua2.attendance_id = ua1.attendance_id
			AND ABS(EXTRACT(EPOCH FROM (ua2.marked_time - ua1.marked_time))) < 60
			AND ua2.id != ua1.id
		)
		GROUP BY ua1.student_id, ua1.attendance_id
	`

	var duplicates []struct {
		StudentID      int
		AttendanceID   int
		DuplicateCount int
	}

	if err := ar.db.Raw(query).Scan(&duplicates).Error; err == nil {
		for _, dup := range duplicates {
			response.Anomalies = append(response.Anomalies, domain.Anomaly{
				Type:              "duplicate_checkin",
				Severity:          "high",
				StudentID:         dup.StudentID,
				EventID:           dup.AttendanceID,
				Description:       "Multiple check-ins detected for this event",
				DetectionTime:     time.Now(),
				RecommendedAction: "Review for possible QR code sharing or technical glitch",
			})
		}
	}

	response.AnomalyCount = len(response.Anomalies)
	response.GeneratedAt = time.Now()
	return &response, nil
}

// GetAnomaliesByStudent returns anomalies for a specific student
func (ar *AnalyticsRepo) GetAnomaliesByStudent(studentID int) ([]domain.Anomaly, error) {
	var anomalies []domain.Anomaly
	// Placeholder: implement based on specific student anomaly detection logic
	return anomalies, nil
}

// ===== Predictions =====

// PredictStudentAttendance predicts future attendance for a student
func (ar *AnalyticsRepo) PredictStudentAttendance(studentID int) (*domain.PredictionResponse, error) {
	var response domain.PredictionResponse

	response.EntityType = "student"
	response.EntityID = studentID
	response.GeneratedAt = time.Now()

	// Simple prediction: average of last 4 weeks
	query := `
		SELECT ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 2)
		FROM user_attendances ua
		WHERE ua.student_id = ? AND ua.marked_time >= NOW() - INTERVAL '4 weeks'
	`
	ar.db.Raw(query, studentID).Scan(&response.CurrentAttendance)

	// Forecast (same as current for basic model)
	response.ForecastedAttendance = response.CurrentAttendance
	response.ConfidenceLevel = 65.0

	response.RiskFactors = []string{}
	if response.CurrentAttendance < 75 {
		response.RiskFactors = append(response.RiskFactors, "Attendance below 75% threshold")
	}

	return &response, nil
}

// PredictCourseAttendance predicts future attendance for a course
func (ar *AnalyticsRepo) PredictCourseAttendance(courseCode string) (*domain.PredictionResponse, error) {
	var response domain.PredictionResponse

	response.EntityType = "course"
	response.GeneratedAt = time.Now()

	// Placeholder implementation
	return &response, nil
}

// ===== Benchmarking =====

// GetBenchmarkComparison returns peer comparison data
func (ar *AnalyticsRepo) GetBenchmarkComparison(entityType string, entityID int) (*domain.BenchmarkResponse, error) {
	var response domain.BenchmarkResponse

	response.EntityType = entityType
	response.EntityID = entityID
	response.GeneratedAt = time.Now()

	switch entityType {
	case "student":
		// Get student's attendance
		query := `
			SELECT ROUND(CAST(SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(*), 0), 2)
			FROM user_attendances
			WHERE student_id = ?
		`
		ar.db.Raw(query, entityID).Scan(&response.PerformanceValue)

		// Get peer average
		query = `
			SELECT ROUND(CAST(SUM(CASE WHEN status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(*), 0), 2)
			FROM user_attendances
		`
		ar.db.Raw(query).Scan(&response.PeerAverage)

	case "course":
		// Similar logic for courses
	}

	if response.PerformanceValue > response.PeerAverage {
		response.PerformanceVsPeers = "above"
	} else if response.PerformanceValue < response.PeerAverage {
		response.PerformanceVsPeers = "below"
	} else {
		response.PerformanceVsPeers = "average"
	}

	response.PercentileRank = (response.PerformanceValue / response.PeerAverage) * 100

	return &response, nil
}

// ===== Utility Methods =====

// GetAttendanceRateForEntity returns attendance rate for any entity
func (ar *AnalyticsRepo) GetAttendanceRateForEntity(entityType string, entityID int, startDate, endDate time.Time) (float64, error) {
	var rate float64

	query := `
		SELECT ROUND(CAST(SUM(CASE WHEN ua.status = 'present' THEN 1 ELSE 0 END) AS FLOAT) * 100 / NULLIF(COUNT(ua.id), 0), 2)
		FROM user_attendances ua
		WHERE ua.%s = ? AND ua.marked_time >= ? AND ua.marked_time <= ?
	`

	field := "student_id"
	if entityType == "course" {
		field = "attendance_id"
	}

	query = fmt.Sprintf(query, field)

	if err := ar.db.Raw(query, entityID, startDate, endDate).Scan(&rate).Error; err != nil {
		return 0, err
	}

	return rate, nil
}

// GetLateCheckInCount returns count of late check-ins
func (ar *AnalyticsRepo) GetLateCheckInCount(studentID int, startDate, endDate time.Time) (int, error) {
	var count int

	query := `
		SELECT COUNT(ua.id)
		FROM user_attendances ua
		JOIN events e ON ua.attendance_id = e.id
		WHERE ua.student_id = ? AND (ua.marked_time - e.start_time) > INTERVAL '5 minutes'
	`

	if !startDate.IsZero() && !endDate.IsZero() {
		query += ` AND ua.marked_time >= ? AND ua.marked_time <= ?`
		if err := ar.db.Raw(query, studentID, startDate, endDate).Scan(&count).Error; err != nil {
			return 0, err
		}
	} else {
		if err := ar.db.Raw(query, studentID).Scan(&count).Error; err != nil {
			return 0, err
		}
	}

	return count, nil
}

// GetAttendanceStreak returns current attendance streak
func (ar *AnalyticsRepo) GetAttendanceStreak(studentID int) (int, error) {
	var streak int

	query := `
		WITH RECURSIVE streak_calc AS (
			SELECT 
				DATE(marked_time) as check_date,
				ROW_NUMBER() OVER (ORDER BY DATE(marked_time) DESC) as rn
			FROM user_attendances
			WHERE student_id = ? AND status = 'present'
			ORDER BY check_date DESC
		)
		SELECT COUNT(*) as streak
		FROM streak_calc
		WHERE rn <= DENSE_RANK() OVER (ORDER BY (DATE(check_date) - (rn || ' days')::INTERVAL))
		LIMIT 1
	`

	if err := ar.db.Raw(query, studentID).Scan(&streak).Error; err != nil {
		// Return 0 if query fails (different DB syntax)
		return 0, nil
	}

	return streak, nil
}
