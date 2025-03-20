package auth

import (
	"net/http"

	"github.com/Dom-HTG/attendance-management-system/internal/auth"
	"github.com/Dom-HTG/attendance-management-system/pkg/responses"
	"github.com/Dom-HTG/attendance-management-system/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthSvcInterface interface {
	RegisterStudent(ctx *gin.Context)
	RegisterLecturer(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthSvc struct {
	repository auth.AuthRepoInterface
}

func RegisterStudent(ctx *gin.Context) {
	var registerUserData *auth.RegisterStudentDT0

	if e := ctx.ShouldBindJSON(&registerUserData); e != nil {
		responses.ApiFailure(ctx, "Unable to bind request body", http.StatusBadRequest, e)
	}

	// hash student password and store it.
	hash, er := utils.HashPassword(registerUserData.Password)
	if er != nil {
		responses.ApiFailure(ctx, "Unable to hash password", http.StatusInternalServerError, er)
	}

	registerUserData.Password = string(hash)

}
