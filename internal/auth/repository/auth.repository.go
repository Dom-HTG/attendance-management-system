package auth

import (
	"github.com/Dom-HTG/attendance-management-system/entities"
	auth "github.com/Dom-HTG/attendance-management-system/internal/auth/domain"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

type AuthRepoInterface interface {
	RegisterStudent(student *auth.RegisterStudentDTO) error
	RegisterLecturer(lecturer *auth.RegisterLecturerDTO) error
	FindStudentByEmail(email string) (*auth.StudentResponse, error)
	FindLecturerByEmail(email string) (*auth.LecturerResponse, error)
	GetStudentByEmailWithPassword(email string) (*entities.Student, error)
	GetLecturerByEmailWithPassword(email string) (*entities.Lecturer, error)
}

func NewAuthRepo(dbInstance *gorm.DB) *AuthRepo {
	return &AuthRepo{
		DB: dbInstance,
	}
}

func (ar *AuthRepo) RegisterStudent(student *auth.RegisterStudentDTO) error {
	// Map DTO to Student entity
	studentEntity := &entities.Student{
		FirstName:    student.FirstName,
		LastName:     student.LastName,
		Email:        student.Email,
		Password:     student.Password,
		MatricNumber: student.MatricNumber,
		Role:         "student",
	}

	tx := ar.DB.Create(&studentEntity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (ar *AuthRepo) RegisterLecturer(lecturer *auth.RegisterLecturerDTO) error {
	// Map DTO to Lecturer entity
	lecturerEntity := &entities.Lecturer{
		FirstName:  lecturer.FirstName,
		LastName:   lecturer.LastName,
		Email:      lecturer.Email,
		Password:   lecturer.Password,
		Department: lecturer.Department,
		StaffID:    lecturer.StaffID,
		Role:       "lecturer",
	}

	tx := ar.DB.Create(&lecturerEntity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (ar *AuthRepo) FindStudentByEmail(email string) (*auth.StudentResponse, error) {
	var student entities.Student

	tx := ar.DB.Where("email = ?", email).First(&student)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &auth.StudentResponse{
		ID:           int(student.ID),
		FirstName:    student.FirstName,
		LastName:     student.LastName,
		Email:        student.Email,
		MatricNumber: student.MatricNumber,
		Role:         student.Role,
		CreatedAt:    student.CreatedAt.String(),
	}, nil
}

func (ar *AuthRepo) FindLecturerByEmail(email string) (*auth.LecturerResponse, error) {
	var lecturer entities.Lecturer

	tx := ar.DB.Where("email = ?", email).First(&lecturer)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &auth.LecturerResponse{
		ID:         int(lecturer.ID),
		FirstName:  lecturer.FirstName,
		LastName:   lecturer.LastName,
		Email:      lecturer.Email,
		Department: lecturer.Department,
		StaffID:    lecturer.StaffID,
		Role:       lecturer.Role,
		CreatedAt:  lecturer.CreatedAt.String(),
	}, nil
}

// Helper methods to get full entities with passwords for login
func (ar *AuthRepo) GetStudentByEmailWithPassword(email string) (*entities.Student, error) {
	var student entities.Student
	tx := ar.DB.Where("email = ?", email).First(&student)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &student, nil
}

func (ar *AuthRepo) GetLecturerByEmailWithPassword(email string) (*entities.Lecturer, error) {
	var lecturer entities.Lecturer
	tx := ar.DB.Where("email = ?", email).First(&lecturer)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &lecturer, nil
}
