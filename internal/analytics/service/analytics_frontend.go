package service

// This file contains additional service methods for frontend-required analytics endpoints
// Added to support lecturer and admin dashboards

import (
	domain "github.com/Dom-HTG/attendance-management-system/internal/analytics/domain"
)

// ===== NEW Lecturer Analytics Service Methods =====

// GetLecturerEvents returns all events created by a lecturer
func (as *AnalyticsService) GetLecturerEvents(lecturerID int) (*domain.LecturerEventsResponse, error) {
	return as.repo.GetLecturerEvents(lecturerID)
}

// GetLecturerSummary returns aggregated statistics for lecturer dashboard
func (as *AnalyticsService) GetLecturerSummary(lecturerID int) (*domain.LecturerSummaryResponse, error) {
	return as.repo.GetLecturerSummary(lecturerID)
}

// ===== NEW Admin Analytics Service Methods =====

// GetAdminOverviewNew returns university-wide statistics for admin dashboard
func (as *AnalyticsService) GetAdminOverviewNew() (*domain.AdminOverviewResponse, error) {
	return as.repo.GetAdminOverviewNew()
}

// GetDepartmentStats returns per-department breakdown
func (as *AnalyticsService) GetDepartmentStats() (*domain.DepartmentStatsResponse, error) {
	return as.repo.GetDepartmentStats()
}
