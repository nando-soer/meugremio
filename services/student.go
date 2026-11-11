package services

import (
	"context"
	"log"
	"net/url"
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

func (student *StudentService) EmailExists(ctx context.Context, email string) (bool, error) {
	return student.model.EmailExists(ctx, email)
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

	data.Password, err = utils.HashPassword(data.Password)
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
	code, err := student.model.QueryCodeByPendingId(ctx, pendingId)
	return (len(code) > 0) && (err != nil), err
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

func (student *StudentService) generateRecoveryLink(address, token string) (string, error) {
	url, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	url.Path = "/usuario/esqueci"

	query := url.Query()
	query.Set("token", token)

	url.RawQuery = query.Encode()

	return url.String(), err
}

func (student *StudentService) InstantiatePasswordRecovery(ctx context.Context, email string) error {
	studentId, err := student.model.QueryStudentIdByEmail(ctx, email)
	if err != nil {
		return err
	}

	token, err := utils.GenerateRecoveryToken()
	if err != nil {
		return err
	}

	link, err := student.generateRecoveryLink("http://localhost:8080", token) // TODO: trocar pra URL de produção
	if err != nil {
		return err
	}

	// err = student.SendConfirmationEmail(ctx, email, link)
	// if err != nil {
	// 	return err
	// }

	log.Println(link)

	hashedToken := utils.HashToken(token)
	err = student.model.InsertPasswordRecovery(ctx, studentId, hashedToken)

	return err
}

func (student *StudentService) RecoverStudentRecoveryId(ctx context.Context, token string) (int64, error) {
	hash := utils.HashToken(token)
	studentId, err := student.model.QueryStudentRecoveryIdByHash(ctx, hash)
	return studentId, err
}

func (student StudentService) ChangePasswordById(ctx context.Context, id int64, newPassword string) error {
	newPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = student.model.UpdatePasswordById(ctx, id, newPassword)
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

	return dto.StudentResponse{
		Name:        data.Name,
		Email:       data.Email,
		Class:       data.Class,
		Description: data.Description,
		Banner:      data.Banner,
		Avatar:      data.Avatar}, nil
}

func (student *StudentService) UpdateProfile(ctx context.Context, sessionId string, data dto.StudentEditData) error {
	studentId, err := student.model.QueryStudentIdBySession(ctx, sessionId)
	if err != nil {
		return err
	}
	return student.model.Update(ctx, studentId, data)
}

func (student *StudentService) ClassExists(ctx context.Context, classCode string) bool {
	name, err := student.model.QueryClassNameByCode(ctx, classCode)
	return (len(name) > 0) && (err == nil)
}

func (student *StudentService) GetAllClasses(ctx context.Context) ([]dto.ClassData, error) {
	return student.model.QueryAllClasses(ctx)
}

func (student *StudentService) GetClass(ctx context.Context, classCode string) (dto.ClassData, error) {
	class, err := student.model.QueryClassByCode(ctx, classCode)
	if err != nil {
		return dto.ClassData{}, err
	}
	return class, err
}

func (student *StudentService) DeletePasswordRecovery(ctx context.Context, token string) error {
	hashedToken := utils.HashToken(token)
	err := student.model.DeleteRecoveryByHash(ctx, hashedToken)
	return err
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
		Name:        data.Name,
		Email:       data.Email,
		Class:       data.Class,
		Description: data.Description,
		Banner:      data.Banner,
		Avatar:      data.Avatar,
	}

	return response, err
}
