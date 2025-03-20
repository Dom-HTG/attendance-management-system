package auth

import "github.com/gin-gonic/gin"

// Repository Interface.
type AuthRepoInterface interface {
	RegisterStudent(student *RegisterStudentDT0) error
	RegisterLecturer(lecturer *RegisterLecturerDTO) error
	// Login() error
}

// Service Interface.
type AuthSvcInterface interface {
	RegisterStudent(ctx *gin.Context)
	RegisterLecturer(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type RegisterStudentDT0 struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	MatricNumber string `json:"matric_number"`
}

type LoginStudentDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterLecturerDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Department string `json:"department"`
	StaffID    string `json:"staff_id"`
}

type LoginLecturerDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordDTO struct {
	Email string `json:"email"`
}

type RefreshTokenDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
