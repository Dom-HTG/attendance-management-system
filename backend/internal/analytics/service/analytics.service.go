package service

import (
	"fmt"
	"time"

	"github.com/Dom-HTG/attendance-management-system/internal/analytics/domain"
	"github.com/Dom-HTG/attendance-management-system/internal/analytics/repository"
)

// AnalyticsServiceInterface defines analytics service operations
type AnalyticsServiceInterface interface {
	// Student analytics
	GetStudentMetrics(studentID int) (*domain.StudentMetricsResponse, error)
	GetStudentInsights(studentID int) (*domain.InsightResponse, error)

	// Lecturer analytics
	GetLecturerCourseMetrics(lecturerID int) (*domain.LecturerCourseMetricsResponse, error)
	GetLecturerCoursePerformance(lecturerID int, courseCode string) (*domain.CoursePerformanceResponse, error)
	GetLecturerInsights(lecturerID int) (*domain.InsightResponse, error)

	// Admin analytics
	GetAdminOverview() (*domain.AdminOverviewResponse, error)
	GetDepartmentMetrics(department string) (*domain.DepartmentDeepDiveResponse, error)
	GetRealTimeDashboard() (*domain.RealTimeDashboardResponse, error)

	// Temporal analytics
	GetTemporalAnalytics(startDate, endDate time.Time, granularity string) (*domain.TemporalAnalyticsResponse, error)

	// Anomalies
	DetectAnomalies() (*domain.AnomalyResponse, error)

	// Predictions
	PredictStudentAttendance(studentID int) (*domain.PredictionResponse, error)
	PredictCourseAttendance(courseCode string) (*domain.PredictionResponse, error)

	// Benchmarking
	GetBenchmarkComparison(entityType string, entityID int) (*domain.BenchmarkResponse, error)

	// Chart data
	GetChartData(chartType string, entityType string, entityID int) (*domain.ChartDataResponse, error)
}

// AnalyticsService implements AnalyticsServiceInterface
type AnalyticsService struct {
	repo repository.AnalyticsRepoInterface
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(repo repository.AnalyticsRepoInterface) AnalyticsServiceInterface {
	return &AnalyticsService{repo: repo}
}

// ===== Student Analytics =====

// GetStudentMetrics returns comprehensive metrics for a student
func (as *AnalyticsService) GetStudentMetrics(studentID int) (*domain.StudentMetricsResponse, error) {
	return as.repo.GetStudentMetrics(studentID)
}

// GetStudentInsights generates natural language insights for a student
func (as *AnalyticsService) GetStudentInsights(studentID int) (*domain.InsightResponse, error) {
	metrics, err := as.repo.GetStudentMetrics(studentID)
	if err != nil {
		return nil, err
	}

	insights := &domain.InsightResponse{
		EntityType:      "student",
		EntityID:        studentID,
		EntityName:      metrics.StudentName,
		Trends:          []domain.TrendExplanation{},
		Recommendations: []domain.Recommendation{},
		GeneratedAt:     time.Now(),
	}

	// Generate summary
	if metrics.OverallAttendanceRate >= 80 {
		insights.Summary = fmt.Sprintf("%s has excellent attendance at %.1f%%. Keep up the consistency and engagement with courses.",
			metrics.StudentName, metrics.OverallAttendanceRate)
		insights.KeyTakeaways = append(insights.KeyTakeaways, "Excellent attendance record")
		insights.KeyTakeaways = append(insights.KeyTakeaways, "High engagement score")
	} else if metrics.OverallAttendanceRate >= 75 {
		insights.Summary = fmt.Sprintf("%s has good attendance at %.1f%%. Continue to maintain consistency throughout the semester.",
			metrics.StudentName, metrics.OverallAttendanceRate)
		insights.KeyTakeaways = append(insights.KeyTakeaways, "Attendance is acceptable")
	} else {
		insights.Summary = fmt.Sprintf("%s has attendance at %.1f%%, which is below the 75%% threshold. Immediate action is recommended.",
			metrics.StudentName, metrics.OverallAttendanceRate)
		insights.KeyTakeaways = append(insights.KeyTakeaways, "Attendance is at-risk")
		insights.Recommendations = append(insights.Recommendations, domain.Recommendation{
			Action:         "Attend more classes regularly",
			Priority:       "high",
			ExpectedImpact: "Improve overall attendance and engagement",
			Timeframe:      "immediate",
		})
	}

	// Add trend analysis
	if len(metrics.AttendanceTrend) > 1 {
		lastTrend := metrics.AttendanceTrend[len(metrics.AttendanceTrend)-1]
		prevTrend := metrics.AttendanceTrend[len(metrics.AttendanceTrend)-2]

		if lastTrend.AttendanceRate > prevTrend.AttendanceRate {
			insights.Trends = append(insights.Trends, domain.TrendExplanation{
				Trend:       "Attendance improving",
				Explanation: "Recent weeks show improvement in attendance rates",
				Timeframe:   "past 4 weeks",
			})
		} else if lastTrend.AttendanceRate < prevTrend.AttendanceRate {
			insights.Trends = append(insights.Trends, domain.TrendExplanation{
				Trend:       "Attendance declining",
				Explanation: "Recent weeks show a decline in attendance",
				Timeframe:   "past 4 weeks",
			})
			insights.Recommendations = append(insights.Recommendations, domain.Recommendation{
				Action:         "Investigate reasons for declining attendance",
				Priority:       "high",
				ExpectedImpact: "Reverse downward trend",
				Timeframe:      "this week",
			})
		}
	}

	// Add late arrival analysis
	if metrics.LateCheckInFrequency > 3 {
		insights.Trends = append(insights.Trends, domain.TrendExplanation{
			Trend:       "Frequent late arrivals",
			Explanation: fmt.Sprintf("Student has %d late check-ins, indicating time management issues", metrics.LateCheckInFrequency),
			Timeframe:   "this semester",
		})
		insights.Recommendations = append(insights.Recommendations, domain.Recommendation{
			Action:         "Arrive on time for classes",
			Priority:       "medium",
			ExpectedImpact: "Improve punctuality and class engagement",
			Timeframe:      "ongoing",
		})
	}

	return insights, nil
}

// ===== Lecturer Analytics =====

// GetLecturerCourseMetrics returns course metrics for a lecturer
func (as *AnalyticsService) GetLecturerCourseMetrics(lecturerID int) (*domain.LecturerCourseMetricsResponse, error) {
	return as.repo.GetLecturerCourseMetrics(lecturerID)
}

// GetLecturerCoursePerformance returns detailed performance for a lecturer's course
func (as *AnalyticsService) GetLecturerCoursePerformance(lecturerID int, courseCode string) (*domain.CoursePerformanceResponse, error) {
	return as.repo.GetLecturerCoursePerformance(lecturerID, courseCode)
}

// GetLecturerInsights generates insights for a lecturer
func (as *AnalyticsService) GetLecturerInsights(lecturerID int) (*domain.InsightResponse, error) {
	metrics, err := as.repo.GetLecturerCourseMetrics(lecturerID)
	if err != nil {
		return nil, err
	}

	insights := &domain.InsightResponse{
		EntityType:      "lecturer",
		EntityID:        lecturerID,
		EntityName:      metrics.LecturerName,
		Trends:          []domain.TrendExplanation{},
		Recommendations: []domain.Recommendation{},
		GeneratedAt:     time.Now(),
	}

	// Generate summary
	if metrics.AverageAttendance >= 80 {
		insights.Summary = fmt.Sprintf("%s maintains excellent class attendance at %.1f%% across all courses. Classes are well-attended and engaging.",
			metrics.LecturerName, metrics.AverageAttendance)
		insights.KeyTakeaways = append(insights.KeyTakeaways, "High average class attendance")
		insights.KeyTakeaways = append(insights.KeyTakeaways, fmt.Sprintf("Generated %d QR codes", metrics.QRGeneratedCount))
	} else if metrics.AverageAttendance >= 70 {
		insights.Summary = fmt.Sprintf("%s has acceptable class attendance at %.1f%%. There may be opportunities to improve student engagement.",
			metrics.LecturerName, metrics.AverageAttendance)
		insights.Recommendations = append(insights.Recommendations, domain.Recommendation{
			Action:         "Consider strategies to improve student engagement",
			Priority:       "medium",
			ExpectedImpact: "Increase class attendance rates",
			Timeframe:      "this month",
		})
	} else {
		insights.Summary = fmt.Sprintf("%s has lower than desired class attendance at %.1f%%. Action to improve engagement is recommended.",
			metrics.LecturerName, metrics.AverageAttendance)
		insights.Recommendations = append(insights.Recommendations, domain.Recommendation{
			Action:         "Review teaching strategies and student engagement methods",
			Priority:       "high",
			ExpectedImpact: "Improve overall course attendance",
			Timeframe:      "immediate",
		})
	}

	insights.Trends = append(insights.Trends, domain.TrendExplanation{
		Trend:       fmt.Sprintf("Managing %d courses", metrics.TotalCourses),
		Explanation: fmt.Sprintf("Average attendance across all courses: %.1f%%", metrics.AverageAttendance),
		Timeframe:   "current semester",
	})

	return insights, nil
}

// ===== Admin Analytics =====

// GetAdminOverview returns university-wide overview
func (as *AnalyticsService) GetAdminOverview() (*domain.AdminOverviewResponse, error) {
	return as.repo.GetAdminOverview()
}

// GetDepartmentMetrics returns department-level metrics
func (as *AnalyticsService) GetDepartmentMetrics(department string) (*domain.DepartmentDeepDiveResponse, error) {
	return as.repo.GetDepartmentMetrics(department)
}

// GetRealTimeDashboard returns real-time dashboard data
func (as *AnalyticsService) GetRealTimeDashboard() (*domain.RealTimeDashboardResponse, error) {
	return as.repo.GetRealTimeDashboard()
}

// ===== Temporal Analytics =====

// GetTemporalAnalytics returns time-based analytics
func (as *AnalyticsService) GetTemporalAnalytics(startDate, endDate time.Time, granularity string) (*domain.TemporalAnalyticsResponse, error) {
	return as.repo.GetTemporalAnalytics(startDate, endDate, granularity)
}

// ===== Anomalies =====

// DetectAnomalies identifies unusual patterns
func (as *AnalyticsService) DetectAnomalies() (*domain.AnomalyResponse, error) {
	return as.repo.DetectAnomalies()
}

// ===== Predictions =====

// PredictStudentAttendance predicts future student attendance
func (as *AnalyticsService) PredictStudentAttendance(studentID int) (*domain.PredictionResponse, error) {
	pred, err := as.repo.PredictStudentAttendance(studentID)
	if err != nil {
		return nil, err
	}

	// Add risk factors and recommendations
	if pred.CurrentAttendance < 75 {
		pred.RiskFactors = append(pred.RiskFactors, "Current attendance below 75%")
		pred.RecommendedActions = append(pred.RecommendedActions, "Student should increase class attendance immediately")
	}

	if pred.CurrentAttendance < 50 {
		pred.RiskFactors = append(pred.RiskFactors, "Critical: Attendance below 50%")
		pred.RecommendedActions = append(pred.RecommendedActions, "Intervention required: Contact student and provide support")
	}

	return pred, nil
}

// PredictCourseAttendance predicts future course attendance
func (as *AnalyticsService) PredictCourseAttendance(courseCode string) (*domain.PredictionResponse, error) {
	return as.repo.PredictCourseAttendance(courseCode)
}

// ===== Benchmarking =====

// GetBenchmarkComparison returns peer comparison data
func (as *AnalyticsService) GetBenchmarkComparison(entityType string, entityID int) (*domain.BenchmarkResponse, error) {
	return as.repo.GetBenchmarkComparison(entityType, entityID)
}

// ===== Chart Data =====

// GetChartData returns data formatted for charts
func (as *AnalyticsService) GetChartData(chartType string, entityType string, entityID int) (*domain.ChartDataResponse, error) {
	response := &domain.ChartDataResponse{
		ChartType:   chartType,
		GeneratedAt: time.Now(),
	}

	switch chartType {
	case "line_trend":
		response.Title = "Attendance Trend"
		response.Description = "Attendance rate over time"

		// Get trend data
		startDate := time.Now().AddDate(0, -3, 0)
		endDate := time.Now()
		trend, err := as.repo.GetStudentAttendanceTrend(entityID, startDate, endDate)
		if err != nil {
			return nil, err
		}

		chartData := domain.LineChartData{
			Datasets: []domain.LineDataset{
				{
					Label: "Attendance Rate",
					Color: "#2ecc71",
				},
			},
		}

		for _, t := range trend {
			chartData.Labels = append(chartData.Labels, t.Period)
			chartData.Datasets[0].Data = append(chartData.Datasets[0].Data, t.AttendanceRate)
		}

		response.DataPoints = chartData

	case "bar_comparison":
		response.Title = "Course Comparison"
		response.Description = "Attendance rates across courses"

		rates, err := as.repo.GetStudentPerCourseRates(entityID)
		if err != nil {
			return nil, err
		}

		chartData := domain.BarChartData{
			Datasets: []domain.BarDataset{
				{
					Label: "Attendance %",
					Color: "#3498db",
				},
			},
		}

		for _, rate := range rates {
			chartData.Labels = append(chartData.Labels, rate.CourseName)
			chartData.Datasets[0].Data = append(chartData.Datasets[0].Data, rate.AttendanceRate)
		}

		response.DataPoints = chartData
	}

	return response, nil
}
