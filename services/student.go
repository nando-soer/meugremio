package services

import (
	"context"
	"log"
	"strings"
	"tcc/dto"
	"tcc/models"
	"tcc/utils"
	"time"
)

type StudentService struct {
	model models.StudentModel
}

func CreateStudentService(model models.StudentModel) StudentService {
	return StudentService{model: model}
}

func (student *StudentService) IsEmailAlreadyUsed(ctx context.Context, email string) (bool, error) {
	return student.model.EmailExists(ctx, email)
}

func (student *StudentService) IsEmailCorrect(ctx context.Context, email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])
	return domain == "aluno.cps.sp.gov.br"
}

func (student *StudentService) Validate(ctx context.Context, submission dto.StudentLogin) (bool, error) {
	data, err := student.model.QueryByEmail(ctx, submission.Email)
	if err != nil {
		return false, err
	}

	isPasswordValid, err := utils.ComparePasswords(submission.Password, data.Password)
	if err != nil {
		return false, err
	}

	return isPasswordValid, err
}

func (student *StudentService) InstantiateSession(ctx context.Context, email string) (string, error) {
	sessionId, err := utils.GenerateSessionId()
	if err != nil {
		return "", err
	}

	studentId, err := student.model.QueryStudentIdByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = student.model.InsertNewSession(ctx, sessionId, studentId, time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	return sessionId, err
}

// TODO
// func (student *StudentService) SendConfirmationEmail(ctx context.Context, email string, code string) error {
// }

func (student *StudentService) InstantiatePendingRegister(ctx context.Context, data dto.StudentData) (string, error) {
	code, err := utils.GenerateConfirmationCode()
	if err != nil {
		return "", err
	}

	log.Println(code)

	data.Password, err = utils.HashPassowrd(data.Password)
	if err != nil {
		return "", err
	}

	pendingId, err := student.model.InsertPendingRegister(ctx, code, data)
	if err != nil {
		return "", err
	}

	// err = student.SendConfirmationEmail(ctx, data.Email, code)
	// if err != nil {
	// 	return "", err
	// }

	return pendingId, err
}

func (student *StudentService) DoesPendingRegisterExists(ctx context.Context, pendingId string) (bool, error) {
	code, err := student.model.QueryCodeByPendingId(ctx, pendingId);
	return len(code) > 0, err
}

func (student *StudentService) ValidatePendingRegister(ctx context.Context, pendingId string, code string) (bool, error) {
	officialCode, err := student.model.QueryCodeByPendingId(ctx, pendingId)
	if err != nil {
		return false, err
	}
	return officialCode == code, err
}

func (student *StudentService) FulfillPendingRegister(ctx context.Context, pendingId string) error {
	data, err := student.model.QueryPendingRegister(ctx, pendingId)
	if err != nil {
		return err
	}

	err = student.model.Insert(ctx, data)
	return err
}

func (student *StudentService) ValidateSessionId(ctx context.Context, sessionId string) (bool, error) {
	_, err := student.model.QueryStudentIdBySession(ctx, sessionId)
	if err != nil {
		return false, err
	}
	return true, err
}

func (student *StudentService) RecoverSession(ctx context.Context, sessionId string) (dto.StudentResponse, error) {
	studentId, err := student.model.QueryStudentIdBySession(ctx, sessionId)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	data, err := student.model.QueryById(ctx, studentId)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	return dto.StudentResponse{Name: data.Name, Email: data.Email, Class: data.Class}, nil
}

// ------------------- TESTE ----------------------------------------------------------------

func (student *StudentService) QueryAll(ctx context.Context) ([]dto.StudentData, error) {
	return student.model.QueryAll(ctx)
}

func (student *StudentService) QuerySensibleDataByEmail(ctx context.Context, email string) (dto.StudentData, error) {
	return student.model.QueryByEmail(ctx, email)
}

func (student *StudentService) QueryNonSensibleDataByEmail(ctx context.Context, email string) (dto.StudentResponse, error) {
	data, err := student.model.QueryByEmail(ctx, email)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	response := dto.StudentResponse{
		Name:  data.Name,
		Email: data.Email,
		Class: data.Class,
	}

	return response, err
}
