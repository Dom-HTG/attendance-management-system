package auth

import (
	"github.com/Dom-HTG/attendance-management-system/entities"
	"github.com/gin-gonic/gin"
)

// Repository Interface.
type AuthRepoInterface interface {
	RegisterStudent(student *RegisterStudentDTO) error
	RegisterLecturer(lecturer *RegisterLecturerDTO) error
	FindStudentByEmail(email string) (*StudentResponse, error)
	FindLecturerByEmail(email string) (*LecturerResponse, error)
	GetStudentByEmailWithPassword(email string) (*entities.Student, error)
	GetLecturerByEmailWithPassword(email string) (*entities.Lecturer, error)
}

// Service Interface.
type AuthSvcInterface interface {
	RegisterStudent(ctx *gin.Context)
	RegisterLecturer(ctx *gin.Context)
	LoginStudent(ctx *gin.Context)
	LoginLecturer(ctx *gin.Context)
}

// Request DTOs
type RegisterStudentDTO struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	MatricNumber string `json:"matric_number" binding:"required"`
}

type LoginStudentDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterLecturerDTO struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	Department string `json:"department" binding:"required"`
	StaffID    string `json:"staff_id" binding:"required"`
}

type LoginLecturerDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Response DTOs
type StudentResponse struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	MatricNumber string `json:"matric_number"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at,omitempty"`
}

type LecturerResponse struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	StaffID    string `json:"staff_id"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at,omitempty"`
}

type LoginResponse struct {
	Message     string      `json:"message"`
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type RefreshTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
