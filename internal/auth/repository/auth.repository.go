package auth

import (
	auth "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	"gorm.io/gorm"
)

type AuthRepoInterface interface {
	RegisterStudent(student *auth.RegisterStudentDT0) error
	RegisterLecturer(lecturer *auth.RegisterLecturerDTO) error
	// Login() error
}

type AuthRepo struct {
	db *gorm.DB
}

func (auth *AuthRepo) RegisterStudent(student *auth.RegisterStudentDT0) error {
	tx := auth.db.Create(&student)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (auth *AuthRepo) RegisterLecturer(lecturer *auth.RegisterLecturerDTO) error {
	tx := auth.db.Create(&lecturer)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
