package auth

import (
	"net/http"

	auth "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	"github.com/Dom-HTG/attendance-management-system/pkg/responses"
	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
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
	var registerUserData *auth.RegisterStudentDT0

	if e := ctx.ShouldBindJSON(&registerUserData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
	}

	// hash password.
	hash, er := utils.HashPassword(registerUserData.Password)
	if er != nil {
		responses.ApiFailure(ctx, "Unable to hash password", http.StatusInternalServerError, er)
	}

	// replace password with hash.
	registerUserData.Password = string(hash)

	//Save user to database.
	if err := svc.Repository.RegisterStudent(registerUserData); err != nil {
		responses.ApiFailure(ctx, "Failed to register student", http.StatusInternalServerError, err)
	}

	responses.ApiSuccess(ctx, http.StatusCreated, "Student successfully registered", &auth.RegisterStudentDT0{})
}

func (svc *AuthSvc) RegisterLecturer(ctx *gin.Context) {
	var registerUserData *auth.RegisterLecturerDTO

	if e := ctx.ShouldBindJSON(&registerUserData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
	}

	// hash password.
	hash, er := utils.HashPassword(registerUserData.Password)
	if er != nil {
		responses.ApiFailure(ctx, "Unable to hash password", http.StatusInternalServerError, er)
	}

	// replace password with hash.
	registerUserData.Password = string(hash)

	//Save user to database.
	if err := svc.Repository.RegisterLecturer(registerUserData); err != nil {
		responses.ApiFailure(ctx, "Failed to register lecturer", http.StatusInternalServerError, err)
	}

	responses.ApiSuccess(ctx, http.StatusCreated, "Lecturer successfully registered", &auth.RegisterLecturerDTO{})
}
