package repository

// This file contains additional repository methods for frontend-required analytics endpoints
// Added to support lecturer and admin dashboards

import (
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	domain "github.com/Dom-HTG/attendance-management-system/internal/analytics/domain"
)

// ===== NEW Lecturer Analytics Methods =====

// GetLecturerEvents returns all events created by a lecturer with attendance counts
func (ar *AnalyticsRepo) GetLecturerEvents(lecturerID int) (*domain.LecturerEventsResponse, error) {
	var events []domain.LecturerEventDetail

	query := `
		SELECT 
			e.id as event_id,
			COALESCE(e.course_name, '') as course_name,
			COALESCE(e.course_code, '') as course_code,
			COALESCE(e.department, '') as department,
			COALESCE(e.venue, '') as venue,
			e.start_time,
			e.end_time,
			COALESCE(e.qr_code_token, '') as qr_token,
			CASE 
				WHEN e.end_time < NOW() THEN 'expired'
				ELSE 'active'
			END as status,
			COUNT(DISTINCT ua.student_id) as total_attendance,
			e.created_at
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		WHERE e.lecturer_id = ?
		GROUP BY e.id
		ORDER BY e.created_at DESC
	`

	if err := ar.db.Raw(query, lecturerID).Scan(&events).Error; err != nil {
		return nil, err
	}

	// Get total unique students reached
	var totalStudentsReached int
	studentQuery := `
		SELECT COUNT(DISTINCT ua.student_id)
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		WHERE e.lecturer_id = ?
	`
	ar.db.Raw(studentQuery, lecturerID).Scan(&totalStudentsReached)

	return &domain.LecturerEventsResponse{
		Events:               events,
		TotalEvents:          len(events),
		TotalStudentsReached: totalStudentsReached,
	}, nil
}

// GetLecturerSummary returns aggregated statistics for lecturer dashboard
func (ar *AnalyticsRepo) GetLecturerSummary(lecturerID int) (*domain.LecturerSummaryResponse, error) {
	var response domain.LecturerSummaryResponse

	// Total events created
	var totalEvents int64
	ar.db.Model(&entities.Event{}).Where("lecturer_id = ?", lecturerID).Count(&totalEvents)
	response.TotalEventsCreated = int(totalEvents)

	// Total students reached
	var totalStudents int
	query := `
		SELECT COUNT(DISTINCT ua.student_id)
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		WHERE e.lecturer_id = ?
	`
	ar.db.Raw(query, lecturerID).Scan(&totalStudents)
	response.TotalStudentsReached = totalStudents

	// Average attendance rate
	var avgRate float64
	avgQuery := `
		SELECT 
			ROUND(CAST(AVG(event_rate) AS NUMERIC), 2)
		FROM (
			SELECT 
				e.id,
				CASE 
					WHEN COUNT(ua.id) > 0 
					THEN (COUNT(CASE WHEN ua.status = 'present' THEN 1 END)::float / COUNT(ua.id)) * 100
					ELSE 0 
				END as event_rate
			FROM events e
			LEFT JOIN user_attendances ua ON e.id = ua.event_id
			WHERE e.lecturer_id = ?
			GROUP BY e.id
		) rates
		WHERE event_rate > 0
	`
	ar.db.Raw(avgQuery, lecturerID).Scan(&avgRate)
	response.AverageAttendanceRate = avgRate

	// Sessions this week
	var sessionsThisWeek int
	weekQuery := `
		SELECT COUNT(*)
		FROM events
		WHERE lecturer_id = ?
		  AND start_time >= date_trunc('week', NOW())
		  AND start_time < date_trunc('week', NOW()) + interval '1 week'
	`
	ar.db.Raw(weekQuery, lecturerID).Scan(&sessionsThisWeek)
	response.SessionsThisWeek = sessionsThisWeek

	// Sessions today
	var sessionsToday int
	todayQuery := `
		SELECT COUNT(*)
		FROM events
		WHERE lecturer_id = ?
		  AND DATE(start_time) = CURRENT_DATE
	`
	ar.db.Raw(todayQuery, lecturerID).Scan(&sessionsToday)
	response.SessionsToday = sessionsToday

	// Most attended course
	type CourseAvg struct {
		CourseCode string
		CourseName string
		AvgRate    float64
	}
	var topCourse CourseAvg
	topCourseQuery := `
		SELECT 
			e.course_code,
			e.course_name,
			ROUND(CAST(AVG(
				CASE 
					WHEN COUNT(ua.id) > 0 
					THEN (COUNT(CASE WHEN ua.status = 'present' THEN 1 END)::float / COUNT(ua.id)) * 100
					ELSE 0 
				END
			) AS NUMERIC), 2) as avg_rate
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		WHERE e.lecturer_id = ? AND e.course_code IS NOT NULL AND e.course_code != ''
		GROUP BY e.course_code, e.course_name
		ORDER BY avg_rate DESC
		LIMIT 1
	`
	err := ar.db.Raw(topCourseQuery, lecturerID).Scan(&topCourse).Error
	if err == nil && topCourse.CourseCode != "" {
		response.MostAttendedCourse = &domain.MostAttendedCourse{
			CourseCode:    topCourse.CourseCode,
			CourseName:    topCourse.CourseName,
			AvgAttendance: topCourse.AvgRate,
		}
	}

	// Attendance trend (last 8 weeks)
	var trends []domain.TrendDataPoint
	trendQuery := `
		SELECT 
			to_char(e.start_time, 'YYYY-MM-DD') as period,
			COUNT(DISTINCT e.id) as total_sessions,
			COUNT(CASE WHEN ua.status = 'present' THEN 1 END) as sessions_attended,
			ROUND(CAST(
				CASE 
					WHEN COUNT(ua.id) > 0 
					THEN (COUNT(CASE WHEN ua.status = 'present' THEN 1 END)::float / COUNT(ua.id)) * 100
					ELSE 0 
				END
			AS NUMERIC), 2) as attendance_rate
		FROM events e
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		WHERE e.lecturer_id = ? 
		  AND e.start_time >= NOW() - interval '8 weeks'
		GROUP BY period
		ORDER BY period
	`
	ar.db.Raw(trendQuery, lecturerID).Scan(&trends)
	response.AttendanceTrend = trends

	return &response, nil
}

// ===== NEW Admin Analytics Methods =====

// GetAdminOverviewNew returns university-wide statistics for admin dashboard
func (ar *AnalyticsRepo) GetAdminOverviewNew() (*domain.AdminOverviewResponse, error) {
	var response domain.AdminOverviewResponse

	// Total students
	var studentCount int64
	ar.db.Model(&entities.Student{}).Count(&studentCount)
	response.TotalStudents = int(studentCount)

	// Total lecturers
	var lecturerCount int64
	ar.db.Model(&entities.Lecturer{}).Count(&lecturerCount)
	response.TotalLecturers = int(lecturerCount)

	// Total departments (count distinct departments from students table)
	var deptCount int64
	ar.db.Model(&entities.Student{}).Distinct("department").Count(&deptCount)
	response.TotalDepartments = int(deptCount)

	// Total events
	var eventCount int64
	ar.db.Model(&entities.Event{}).Count(&eventCount)
	response.TotalEvents = int(eventCount)

	// Average attendance rate
	var avgRate float64
	avgQuery := `
		SELECT 
			COALESCE(ROUND(CAST(
				AVG(
					CASE 
						WHEN event_totals.total > 0 
						THEN (event_totals.present::float / event_totals.total) * 100 
						ELSE 0 
					END
				)
			AS NUMERIC), 2), 0)
		FROM (
			SELECT 
				event_id,
				COUNT(*) as total,
				COUNT(CASE WHEN status = 'present' THEN 1 END) as present
			FROM user_attendances
			GROUP BY event_id
		) event_totals
	`
	ar.db.Raw(avgQuery).Scan(&avgRate)
	response.AverageAttendanceRate = avgRate

	// Active sessions now
	var activeSessions int
	activeQuery := `
		SELECT COUNT(*) 
		FROM events 
		WHERE NOW() BETWEEN start_time AND end_time
	`
	ar.db.Raw(activeQuery).Scan(&activeSessions)
	response.ActiveSessionsNow = activeSessions

	// QR codes generated today
	var qrToday int
	qrQuery := `
		SELECT COUNT(*) 
		FROM events 
		WHERE DATE(created_at) = CURRENT_DATE
	`
	ar.db.Raw(qrQuery).Scan(&qrToday)
	response.QRCodesGeneratedToday = qrToday

	// Total check-ins today
	var checkInsToday int
	checkInQuery := `
		SELECT COUNT(*) 
		FROM user_attendances 
		WHERE DATE(marked_time) = CURRENT_DATE
	`
	ar.db.Raw(checkInQuery).Scan(&checkInsToday)
	response.TotalCheckInsToday = checkInsToday

	// System health (basic implementation)
	response.SystemHealth = &domain.SystemHealth{
		DatabaseStatus: "healthy",
		LastCheckIn:    time.Now(),
		UptimeHours:    0, // Would need separate tracking
	}

	response.GeneratedAt = time.Now()

	// Also populate legacy fields for compatibility
	response.OverallAttendanceRate = avgRate
	response.TotalActiveSessions = activeSessions

	return &response, nil
}

// GetDepartmentStats returns per-department breakdown
func (ar *AnalyticsRepo) GetDepartmentStats() (*domain.DepartmentStatsResponse, error) {
	var departments []domain.DepartmentStat

	// Query based on events table department field since students don't have department
	query := `
		SELECT 
			COALESCE(dept.department, 'Unknown') as department,
			COUNT(DISTINCT s.id) as total_students,
			COUNT(DISTINCT l.id) as total_lecturers,
			COUNT(DISTINCT e.id) as total_events,
			COALESCE(ROUND(CAST(
				AVG(
					CASE 
						WHEN dept_attendance.total > 0 
						THEN (dept_attendance.present::float / dept_attendance.total) * 100
						ELSE 0 
					END
				)
			AS NUMERIC), 2), 0) as average_attendance_rate,
			COUNT(ua.id) as total_check_ins
		FROM (SELECT DISTINCT department FROM events WHERE department IS NOT NULL AND department != '') dept
		LEFT JOIN events e ON dept.department = e.department
		LEFT JOIN user_attendances ua ON e.id = ua.event_id
		LEFT JOIN students s ON ua.student_id = s.id
		LEFT JOIN lecturers l ON dept.department = l.department
		LEFT JOIN (
			SELECT 
				e2.department,
				COUNT(*) as total,
				COUNT(CASE WHEN ua2.status = 'present' THEN 1 END) as present
			FROM events e2
			LEFT JOIN user_attendances ua2 ON e2.id = ua2.event_id
			WHERE e2.department IS NOT NULL
			GROUP BY e2.department
		) dept_attendance ON dept.department = dept_attendance.department
		GROUP BY dept.department
		ORDER BY total_events DESC
	`

	if err := ar.db.Raw(query).Scan(&departments).Error; err != nil {
		return nil, err
	}

	return &domain.DepartmentStatsResponse{
		Departments: departments,
	}, nil
}
