package auth

import (
	auth "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(dbInstance *gorm.DB) *AuthRepo {
	return &AuthRepo{
		DB: dbInstance,
	}
}

func (auth *AuthRepo) RegisterStudent(student *auth.RegisterStudentDT0) error {
	tx := auth.DB.Create(&student)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (auth *AuthRepo) RegisterLecturer(lecturer *auth.RegisterLecturerDTO) error {
	tx := auth.DB.Create(&lecturer)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
