package handler

// This file contains additional handler methods for frontend-required analytics endpoints
// Added to support lecturer and admin dashboards

import (
	"net/http"

	"github.com/Dom-HTG/attendance-management-system/pkg/middleware"
	"github.com/Dom-HTG/attendance-management-system/pkg/responses"
	"github.com/gin-gonic/gin"
)

// ===== NEW Lecturer Analytics Handlers =====

// GetLecturerEvents handles GET /api/events/lecturer
func (ah *AnalyticsHandler) GetLecturerEvents(ctx *gin.Context) {
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	events, err := ah.service.GetLecturerEvents(lecturerID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve lecturer events", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Events retrieved successfully", events)
}

// GetLecturerSummary handles GET /api/analytics/lecturer/summary
func (ah *AnalyticsHandler) GetLecturerSummary(ctx *gin.Context) {
	lecturerID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		responses.ApiFailure(ctx, "User not found in context", http.StatusUnauthorized, nil)
		return
	}

	summary, err := ah.service.GetLecturerSummary(lecturerID)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve lecturer summary", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Lecturer summary retrieved successfully", summary)
}

// ===== NEW Admin Analytics Handlers =====

// GetAdminOverviewNew handles GET /api/analytics/admin/overview
func (ah *AnalyticsHandler) GetAdminOverviewNew(ctx *gin.Context) {
	overview, err := ah.service.GetAdminOverviewNew()
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve admin overview", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Admin overview retrieved successfully", overview)
}

// GetDepartmentStats handles GET /api/analytics/admin/departments
func (ah *AnalyticsHandler) GetDepartmentStats(ctx *gin.Context) {
	stats, err := ah.service.GetDepartmentStats()
	if err != nil {
		responses.ApiFailure(ctx, "Failed to retrieve department statistics", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Department statistics retrieved successfully", stats)
}
