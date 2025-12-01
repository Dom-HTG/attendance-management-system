package auth

import (
	"errors"
	"net/http"

	auth "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	"github.com/Dom-HTG/attendance-management-system/pkg/responses"
	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthSvc struct {
	Repository auth.AuthRepoInterface
}

// constructor.
func NewAuthSvc(repo auth.AuthRepoInterface) *AuthSvc {
	return &AuthSvc{
		Repository: repo,
	}
}

func (svc *AuthSvc) RegisterStudent(ctx *gin.Context) {
	var registerUserData auth.RegisterStudentDTO

	if e := ctx.ShouldBindJSON(&registerUserData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
		return
	}

	// hash password.
	hash, er := utils.HashPassword(registerUserData.Password)
	if er != nil {
		responses.ApiFailure(ctx, "Unable to hash password", http.StatusInternalServerError, er)
		return
	}

	// replace password with hash.
	registerUserData.Password = string(hash)

	// Save user to database.
	if err := svc.Repository.RegisterStudent(&registerUserData); err != nil {
		responses.ApiFailure(ctx, "Failed to register student", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusCreated, "Student successfully registered", map[string]string{
		"message": "Student successfully registered. Please login with your credentials.",
	})
}

func (svc *AuthSvc) RegisterLecturer(ctx *gin.Context) {
	var registerUserData auth.RegisterLecturerDTO

	if e := ctx.ShouldBindJSON(&registerUserData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
		return
	}

	// hash password.
	hash, er := utils.HashPassword(registerUserData.Password)
	if er != nil {
		responses.ApiFailure(ctx, "Unable to hash password", http.StatusInternalServerError, er)
		return
	}

	// replace password with hash.
	registerUserData.Password = string(hash)

	// Save user to database.
	if err := svc.Repository.RegisterLecturer(&registerUserData); err != nil {
		responses.ApiFailure(ctx, "Failed to register lecturer", http.StatusInternalServerError, err)
		return
	}

	responses.ApiSuccess(ctx, http.StatusCreated, "Lecturer successfully registered", map[string]string{
		"message": "Lecturer successfully registered. Please login with your credentials.",
	})
}

func (svc *AuthSvc) LoginStudent(ctx *gin.Context) {
	var loginData *auth.LoginStudentDTO

	if e := ctx.ShouldBindJSON(&loginData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
		return
	}

	// Get student by email with password for comparison
	studentEntity, err := svc.Repository.GetStudentByEmailWithPassword(loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ApiFailure(ctx, "Invalid email or password", http.StatusUnauthorized, nil)
			return
		}
		responses.ApiFailure(ctx, "Database error", http.StatusInternalServerError, err)
		return
	}

	// Compare passwords
	if !utils.CompareHash(loginData.Password, studentEntity.Password) {
		responses.ApiFailure(ctx, "Invalid email or password", http.StatusUnauthorized, nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(studentEntity.ID), studentEntity.Email, "student", 60)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate token", http.StatusInternalServerError, err)
		return
	}

	// Prepare response
	studentResponse := &auth.StudentResponse{
		ID:           int(studentEntity.ID),
		FirstName:    studentEntity.FirstName,
		LastName:     studentEntity.LastName,
		Email:        studentEntity.Email,
		MatricNumber: studentEntity.MatricNumber,
		Role:         studentEntity.Role,
	}

	loginResponse := &auth.LoginResponse{
		Message:     "Student login successful",
		AccessToken: token,
		User:        studentResponse,
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Login successful", loginResponse)
}

func (svc *AuthSvc) LoginLecturer(ctx *gin.Context) {
	var loginData *auth.LoginLecturerDTO

	if e := ctx.ShouldBindJSON(&loginData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
		return
	}

	// Get lecturer by email with password for comparison
	lecturerEntity, err := svc.Repository.GetLecturerByEmailWithPassword(loginData.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			responses.ApiFailure(ctx, "Invalid email or password", http.StatusUnauthorized, nil)
			return
		}
		responses.ApiFailure(ctx, "Database error", http.StatusInternalServerError, err)
		return
	}

	// Compare passwords
	if !utils.CompareHash(loginData.Password, lecturerEntity.Password) {
		responses.ApiFailure(ctx, "Invalid email or password", http.StatusUnauthorized, nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(lecturerEntity.ID), lecturerEntity.Email, "lecturer", 60)
	if err != nil {
		responses.ApiFailure(ctx, "Failed to generate token", http.StatusInternalServerError, err)
		return
	}

	// Prepare response
	lecturerResponse := &auth.LecturerResponse{
		ID:         int(lecturerEntity.ID),
		FirstName:  lecturerEntity.FirstName,
		LastName:   lecturerEntity.LastName,
		Email:      lecturerEntity.Email,
		Department: lecturerEntity.Department,
		StaffID:    lecturerEntity.StaffID,
		Role:       lecturerEntity.Role,
	}

	loginResponse := &auth.LoginResponse{
		Message:     "Lecturer login successful",
		AccessToken: token,
		User:        lecturerResponse,
	}

	responses.ApiSuccess(ctx, http.StatusOK, "Login successful", loginResponse)
}
