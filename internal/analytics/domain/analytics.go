package domain

import "time"

// ===== Student Analytics =====

// StudentMetricsResponse represents overall student attendance metrics
type StudentMetricsResponse struct {
	StudentID              int                    `json:"student_id"`
	StudentName            string                 `json:"student_name"`
	MatricNumber           string                 `json:"matric_number"`
	OverallAttendanceRate  float64                `json:"overall_attendance_rate"`
	TotalSessions          int                    `json:"total_sessions"`
	TotalPresent           int                    `json:"total_present"`
	TotalAbsent            int                    `json:"total_absent"`
	TotalLate              int                    `json:"total_late"`
	AttendanceStreak       int                    `json:"attendance_streak"`
	LateCheckInFrequency   int                    `json:"late_checkin_frequency"`
	ClassAverageComparison float64                `json:"class_average_comparison"`
	AtRiskStatus           bool                   `json:"at_risk_status"`
	EngagementScore        float64                `json:"engagement_score"` // 0-100
	PerCourseRates         []CourseAttendanceRate `json:"per_course_rates"`
	AttendanceTrend        []TrendDataPoint       `json:"attendance_trend"`
	GeneratedAt            time.Time              `json:"generated_at"`
}

// CourseAttendanceRate represents per-course attendance
type CourseAttendanceRate struct {
	CourseCode       string  `json:"course_code"`
	CourseName       string  `json:"course_name"`
	AttendanceRate   float64 `json:"attendance_rate"`
	SessionsAttended int     `json:"sessions_attended"`
	TotalSessions    int     `json:"total_sessions"`
	Department       string  `json:"department"`
}

// TrendDataPoint represents attendance over time
type TrendDataPoint struct {
	Period             string  `json:"period"` // YYYY-WW for weekly, YYYY-MM for monthly
	AttendanceRate     float64 `json:"attendance_rate"`
	SessionsAttended   int     `json:"sessions_attended"`
	TotalSessions      int     `json:"total_sessions"`
	AverageCheckInTime int     `json:"average_checkin_time_minutes"` // avg minutes from session start
}

// ===== Lecturer Analytics =====

// LecturerCourseMetricsResponse represents per-course metrics for a lecturer
type LecturerCourseMetricsResponse struct {
	LecturerID        int             `json:"lecturer_id"`
	LecturerName      string          `json:"lecturer_name"`
	Department        string          `json:"department"`
	TotalCourses      int             `json:"total_courses"`
	AverageAttendance float64         `json:"average_attendance_rate"`
	CourseMetrics     []CourseMetrics `json:"course_metrics"`
	QRGeneratedCount  int             `json:"qr_generated_count"`
	GeneratedAt       time.Time       `json:"generated_at"`
}

// CourseMetrics represents detailed metrics for a single course
type CourseMetrics struct {
	CourseCode           string           `json:"course_code"`
	CourseName           string           `json:"course_name"`
	AttendanceAverage    float64          `json:"attendance_average"`
	SessionCount         int              `json:"session_count"`
	StudentCount         int              `json:"student_count"`
	MostAttendedSession  SessionSummary   `json:"most_attended_session"`
	LeastAttendedSession SessionSummary   `json:"least_attended_session"`
	DropOffTrend         []TrendDataPoint `json:"dropoff_trend"`
	AverageCheckInTime   int              `json:"average_checkin_time_minutes"`
}

// SessionSummary represents a single session
type SessionSummary struct {
	EventID         int       `json:"event_id"`
	EventName       string    `json:"event_name"`
	StartTime       time.Time `json:"start_time"`
	AttendanceRate  float64   `json:"attendance_rate"`
	StudentsPresent int       `json:"students_present"`
	TotalEnrolled   int       `json:"total_enrolled"`
}

// CoursePerformanceResponse for lecturer course analytics
type CoursePerformanceResponse struct {
	CourseCode                  string                 `json:"course_code"`
	CourseName                  string                 `json:"course_name"`
	LecturerName                string                 `json:"lecturer_name"`
	Department                  string                 `json:"department"`
	StudentCount                int                    `json:"student_count"`
	OverallAttendanceRate       float64                `json:"overall_attendance_rate"`
	AttendanceDistribution      AttendanceDistribution `json:"attendance_distribution"`
	StudentsAtRisk              int                    `json:"students_at_risk"`
	AverageCheckInTime          int                    `json:"average_checkin_time_minutes"`
	LateArrivalsCount           int                    `json:"late_arrivals_count"`
	SessionDurationVsAttendance []DurationCorrelation  `json:"session_duration_vs_attendance"`
	GeneratedAt                 time.Time              `json:"generated_at"`
}

// AttendanceDistribution represents histogram of attendance rates
type AttendanceDistribution struct {
	Range0To20   int `json:"range_0_to_20"`
	Range20To40  int `json:"range_20_to_40"`
	Range40To60  int `json:"range_40_to_60"`
	Range60To80  int `json:"range_60_to_80"`
	Range80To100 int `json:"range_80_to_100"`
}

// DurationCorrelation represents correlation between session duration and attendance
type DurationCorrelation struct {
	DurationMinutes int     `json:"duration_minutes"`
	AttendanceRate  float64 `json:"attendance_rate"`
	SessionCount    int     `json:"session_count"`
}

// ===== Admin Analytics =====

// AdminOverviewResponse represents university-wide metrics
type AdminOverviewResponse struct {
	OverallAttendanceRate   float64               `json:"overall_attendance_rate"`
	TotalActiveSessions     int                   `json:"total_active_sessions"`
	TotalStudents           int                   `json:"total_students"`
	TotalLecturers          int                   `json:"total_lecturers"`
	DepartmentComparison    []DepartmentMetrics   `json:"department_comparison"`
	TopPerformingCourses    []CourseMetrics       `json:"top_performing_courses"`
	LowestPerformingCourses []CourseMetrics       `json:"lowest_performing_courses"`
	LecturerPerformance     []LecturerPerformance `json:"lecturer_performance"`
	GeneratedAt             time.Time             `json:"generated_at"`
}

// DepartmentMetrics represents department-level metrics
type DepartmentMetrics struct {
	DepartmentName     string           `json:"department_name"`
	AttendanceRate     float64          `json:"attendance_rate"`
	StudentCount       int              `json:"student_count"`
	LecturerCount      int              `json:"lecturer_count"`
	CourseCount        int              `json:"course_count"`
	AverageCheckInTime int              `json:"average_checkin_time_minutes"`
	TrendOverSemester  []TrendDataPoint `json:"trend_over_semester"`
}

// LecturerPerformance represents lecturer-level performance
type LecturerPerformance struct {
	LecturerID             int     `json:"lecturer_id"`
	LecturerName           string  `json:"lecturer_name"`
	Department             string  `json:"department"`
	AverageClassAttendance float64 `json:"average_class_attendance"`
	CoursesManaged         int     `json:"courses_managed"`
	QRSessionsCreated      int     `json:"qr_sessions_created"`
	EfficiencyScore        float64 `json:"efficiency_score"` // 0-100
}

// DepartmentDeepDiveResponse for department-specific analytics
type DepartmentDeepDiveResponse struct {
	DepartmentName               string                  `json:"department_name"`
	OverallAttendanceRate        float64                 `json:"overall_attendance_rate"`
	StudentCount                 int                     `json:"student_count"`
	LecturerCount                int                     `json:"lecturer_count"`
	CourseCount                  int                     `json:"course_count"`
	AttendanceTrend              []TrendDataPoint        `json:"attendance_trend"`
	CourseEnrollmentVsAttendance []CourseEnrollmentData  `json:"course_enrollment_vs_attendance"`
	LecturerEfficiency           []LecturerEfficiency    `json:"lecturer_efficiency"`
	StudentEngagementByYear      []StudentEngagementYear `json:"student_engagement_by_year"`
	VenueUtilization             []VenueUtilization      `json:"venue_utilization"`
	GeneratedAt                  time.Time               `json:"generated_at"`
}

// CourseEnrollmentData for enrollment vs attendance analysis
type CourseEnrollmentData struct {
	CourseCode     string  `json:"course_code"`
	CourseName     string  `json:"course_name"`
	Enrolled       int     `json:"enrolled"`
	ActualAttended int     `json:"actual_attended"`
	AttendanceRate float64 `json:"attendance_rate"`
}

// LecturerEfficiency for lecturer metrics
type LecturerEfficiency struct {
	LecturerID        int     `json:"lecturer_id"`
	LecturerName      string  `json:"lecturer_name"`
	QRSessionsCreated int     `json:"qr_sessions_created"`
	SessionsConducted int     `json:"sessions_conducted"`
	EfficiencyRate    float64 `json:"efficiency_rate"`
}

// StudentEngagementYear for student engagement by year
type StudentEngagementYear struct {
	Year              string  `json:"year"` // e.g., "100L", "200L"
	StudentCount      int     `json:"student_count"`
	AverageAttendance float64 `json:"average_attendance"`
	EngagementScore   float64 `json:"engagement_score"`
}

// VenueUtilization for venue analysis
type VenueUtilization struct {
	Venue             string  `json:"venue"`
	SessionsHeld      int     `json:"sessions_held"`
	AverageAttendance int     `json:"average_attendance"`
	Capacity          int     `json:"capacity"`
	UtilizationRate   float64 `json:"utilization_rate"`
}

// RealTimeDashboardResponse for real-time admin dashboard
type RealTimeDashboardResponse struct {
	ActiveSessionsNow      int              `json:"active_sessions_now"`
	TotalCheckInsToday     int              `json:"total_checkins_today"`
	AverageAttendanceToday float64          `json:"average_attendance_today"`
	OngoingSessions        []OngoingSession `json:"ongoing_sessions"`
	SystemUsageStats       SystemUsageStats `json:"system_usage_stats"`
	GeneratedAt            time.Time        `json:"generated_at"`
}

// OngoingSession represents a currently active session
type OngoingSession struct {
	EventID          int       `json:"event_id"`
	CourseName       string    `json:"course_name"`
	Lecturer         string    `json:"lecturer"`
	Venue            string    `json:"venue"`
	StartTime        time.Time `json:"start_time"`
	CheckInsCount    int       `json:"checkins_count"`
	StudentsEnrolled int       `json:"students_enrolled"`
	AttendanceRate   float64   `json:"attendance_rate"`
}

// SystemUsageStats for system-wide usage
type SystemUsageStats struct {
	TotalAPICallsToday     int `json:"total_api_calls_today"`
	QRCodesGeneratedToday  int `json:"qr_codes_generated_today"`
	CheckInsProcessedToday int `json:"checkins_processed_today"`
	ActiveUsersToday       int `json:"active_users_today"`
}

// ===== Benchmarking =====

// BenchmarkResponse for comparative analytics
type BenchmarkResponse struct {
	EntityType           string         `json:"entity_type"` // student, course, department
	EntityID             int            `json:"entity_id"`
	EntityName           string         `json:"entity_name"`
	PerformanceValue     float64        `json:"performance_value"`
	PeerAverage          float64        `json:"peer_average"`
	PeerStdDev           float64        `json:"peer_std_dev"`
	PercentileRank       float64        `json:"percentile_rank"`      // 0-100
	PerformanceVsPeers   string         `json:"performance_vs_peers"` // "above", "average", "below"
	HistoricalComparison HistoricalData `json:"historical_comparison"`
	GoalTracking         GoalTracking   `json:"goal_tracking"`
	GeneratedAt          time.Time      `json:"generated_at"`
}

// HistoricalData for semester comparison
type HistoricalData struct {
	CurrentSemester  float64 `json:"current_semester"`
	PreviousSemester float64 `json:"previous_semester"`
	ChangePercent    float64 `json:"change_percent"`
	TrendDirection   string  `json:"trend_direction"` // "up", "down", "stable"
}

// GoalTracking for goal vs actual
type GoalTracking struct {
	TargetAttendance float64 `json:"target_attendance"`
	ActualAttendance float64 `json:"actual_attendance"`
	GoalMetStatus    string  `json:"goal_met_status"` // "on_track", "at_risk", "failed"
	DaysRemaining    int     `json:"days_remaining_in_period"`
}

// ===== Temporal Analytics =====

// TemporalAnalyticsResponse for time-based patterns
type TemporalAnalyticsResponse struct {
	Granularity       string              `json:"granularity"` // daily, weekly, monthly
	StartDate         time.Time           `json:"start_date"`
	EndDate           time.Time           `json:"end_date"`
	AttendanceHeatmap []HeatmapCell       `json:"attendance_heatmap"`
	SeasonalTrends    []TrendDataPoint    `json:"seasonal_trends"`
	DayOfWeekAnalysis []DayOfWeekMetrics  `json:"day_of_week_analysis"`
	HolidayImpact     []HolidayImpactData `json:"holiday_impact"`
	GeneratedAt       time.Time           `json:"generated_at"`
}

// HeatmapCell represents attendance at specific day/time
type HeatmapCell struct {
	DayOfWeek      string  `json:"day_of_week"` // Monday, Tuesday, etc.
	TimeSlot       string  `json:"time_slot"`   // HH:MM (e.g., "10:00")
	AttendanceRate float64 `json:"attendance_rate"`
	SessionCount   int     `json:"session_count"`
	AvgCheckInTime int     `json:"avg_checkin_time_minutes"`
}

// DayOfWeekMetrics for day-of-week performance
type DayOfWeekMetrics struct {
	DayOfWeek      string            `json:"day_of_week"`
	AttendanceRate float64           `json:"attendance_rate"`
	SessionCount   int               `json:"session_count"`
	AveragePresent int               `json:"average_present"`
	TimeSlots      []TimeSlotMetrics `json:"time_slots"`
}

// TimeSlotMetrics for intra-day analysis
type TimeSlotMetrics struct {
	TimeSlot       string  `json:"time_slot"`
	AttendanceRate float64 `json:"attendance_rate"`
	SessionCount   int     `json:"session_count"`
}

// HolidayImpactData for holiday period analysis
type HolidayImpactData struct {
	PeriodName       string    `json:"period_name"` // "Christmas", "Exam Week", etc.
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	AttendanceRate   float64   `json:"attendance_rate"`
	BeforePeriodRate float64   `json:"before_period_rate"`
	ImpactPercent    float64   `json:"impact_percent"`
}

// ===== Anomalies & Predictions =====

// AnomalyResponse for anomaly detection
type AnomalyResponse struct {
	AnomalyCount      int       `json:"anomaly_count"`
	CriticalAnomalies int       `json:"critical_anomalies"`
	Anomalies         []Anomaly `json:"anomalies"`
	GeneratedAt       time.Time `json:"generated_at"`
}

// Anomaly represents a detected anomaly
type Anomaly struct {
	ID                int       `json:"id"`
	Type              string    `json:"type"`     // "unusual_pattern", "fraud_suspected", "duplicate_checkin", "timing_anomaly"
	Severity          string    `json:"severity"` // "low", "medium", "high", "critical"
	Description       string    `json:"description"`
	StudentID         int       `json:"student_id,omitempty"`
	StudentName       string    `json:"student_name,omitempty"`
	EventID           int       `json:"event_id,omitempty"`
	CourseName        string    `json:"course_name,omitempty"`
	DetectionTime     time.Time `json:"detection_time"`
	RecommendedAction string    `json:"recommended_action"`
}

// PredictionResponse for predictive analytics
type PredictionResponse struct {
	EntityType           string              `json:"entity_type"` // student, course, department
	EntityID             int                 `json:"entity_id"`
	EntityName           string              `json:"entity_name"`
	ForecastedAttendance float64             `json:"forecasted_attendance"`
	CurrentAttendance    float64             `json:"current_attendance"`
	ConfidenceLevel      float64             `json:"confidence_level"` // 0-100
	StudentsPredicted    []StudentPrediction `json:"students_predicted,omitempty"`
	RiskFactors          []string            `json:"risk_factors"`
	RecommendedActions   []string            `json:"recommended_actions"`
	GeneratedAt          time.Time           `json:"generated_at"`
}

// StudentPrediction represents prediction for individual student
type StudentPrediction struct {
	StudentID               int     `json:"student_id"`
	StudentName             string  `json:"student_name"`
	PredictedAttendanceRate float64 `json:"predicted_attendance_rate"`
	RiskLevel               string  `json:"risk_level"`            // "low", "medium", "high"
	LikelihoodDropOut       float64 `json:"likelihood_dropout"`    // 0-1
	DaysBeforeThreshold     int     `json:"days_before_threshold"` // -1 if already below
}

// InsightResponse for natural language insights
type InsightResponse struct {
	EntityType      string             `json:"entity_type"`
	EntityID        int                `json:"entity_id"`
	EntityName      string             `json:"entity_name"`
	Summary         string             `json:"summary"` // 2-3 sentence plain English summary
	KeyTakeaways    []string           `json:"key_takeaways"`
	Trends          []TrendExplanation `json:"trends"`
	Recommendations []Recommendation   `json:"recommendations"`
	GeneratedAt     time.Time          `json:"generated_at"`
}

// TrendExplanation represents explanation of a trend
type TrendExplanation struct {
	Trend       string `json:"trend"` // e.g., "Attendance increasing", "Late arrivals rising"
	Explanation string `json:"explanation"`
	Timeframe   string `json:"timeframe"` // e.g., "past 4 weeks"
}

// Recommendation represents an actionable recommendation
type Recommendation struct {
	Action         string `json:"action"`
	Priority       string `json:"priority"` // "high", "medium", "low"
	ExpectedImpact string `json:"expected_impact"`
	Timeframe      string `json:"timeframe"` // e.g., "immediate", "this week"
}

// ===== Chart Data =====

// ChartDataResponse for visualization
type ChartDataResponse struct {
	ChartType   string      `json:"chart_type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	DataPoints  interface{} `json:"data_points"` // Type varies by chart
	Labels      []string    `json:"labels"`
	GeneratedAt time.Time   `json:"generated_at"`
}

// LineChartData for trend visualization
type LineChartData struct {
	Datasets []LineDataset `json:"datasets"`
	Labels   []string      `json:"labels"`
}

// LineDataset for line chart
type LineDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Color string    `json:"color,omitempty"`
}

// BarChartData for comparison
type BarChartData struct {
	Labels   []string     `json:"labels"`
	Datasets []BarDataset `json:"datasets"`
}

// BarDataset for bar chart
type BarDataset struct {
	Label string    `json:"label"`
	Data  []float64 `json:"data"`
	Color string    `json:"color,omitempty"`
}

// PieChartData for distribution
type PieChartData struct {
	Labels []string  `json:"labels"`
	Data   []float64 `json:"data"`
	Colors []string  `json:"colors,omitempty"`
}

// ===== Export/Report Request =====

// GenerateReportRequest for report generation
type GenerateReportRequest struct {
	ReportType  string    `json:"report_type"` // student_summary, course_summary, department_summary
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Departments []string  `json:"departments,omitempty"`
	Courses     []string  `json:"courses,omitempty"`
	Students    []int     `json:"students,omitempty"`
	Format      string    `json:"format"` // pdf, csv, json
}

// ReportGenerationResponse for report result
type ReportGenerationResponse struct {
	ReportID    string     `json:"report_id"`
	ReportType  string     `json:"report_type"`
	Status      string     `json:"status"` // processing, completed, failed
	DownloadURL string     `json:"download_url,omitempty"`
	GeneratedAt time.Time  `json:"generated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	FileSize    int        `json:"file_size,omitempty"`
}

// ScheduleReportRequest for scheduled reports
type ScheduleReportRequest struct {
	ReportConfiguration GenerateReportRequest `json:"report_configuration"`
	Recipients          []string              `json:"recipients"`
	Frequency           string                `json:"frequency"`       // daily, weekly, monthly
	DeliveryMethod      string                `json:"delivery_method"` // email, dashboard
	StartDate           time.Time             `json:"start_date"`
	EndDate             *time.Time            `json:"end_date,omitempty"`
}

// ScheduledReportResponse for scheduled report status
type ScheduledReportResponse struct {
	ScheduleID  string     `json:"schedule_id"`
	ReportType  string     `json:"report_type"`
	Frequency   string     `json:"frequency"`
	Recipients  []string   `json:"recipients"`
	NextRunTime time.Time  `json:"next_run_time"`
	LastRunTime *time.Time `json:"last_run_time,omitempty"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ===== Alert Configuration =====

// AlertConfiguration for threshold-based alerts
type AlertConfiguration struct {
	ID                 int       `json:"id"`
	EntityType         string    `json:"entity_type"` // student, course, department
	Condition          string    `json:"condition"`   // attendance_below, late_arrivals_exceeding
	Threshold          float64   `json:"threshold"`
	AlertRecipients    []string  `json:"alert_recipients"`
	NotificationMethod string    `json:"notification_method"` // email, in_app, sms
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
}
