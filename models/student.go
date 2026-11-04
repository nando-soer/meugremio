package models

import (
	"context"
	"encoding/json"
	"tcc/dto"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentModel struct {
	pool *pgxpool.Pool
}

func CreateStudentModel(pool *pgxpool.Pool) StudentModel {
	return StudentModel{pool: pool}
}

func (model *StudentModel) Insert(ctx context.Context, data dto.StudentData) (err error) {
	_, err = model.pool.Exec(ctx, "INSERT INTO students (name, email, class, password) VALUES ($1, $2, $3, $4)", data.Name, data.Email, data.Class, data.Password)
	return
}

func (model *StudentModel) InsertNewSession(ctx context.Context, sessionId string, studentId int64, expiration time.Time) error {
	_, err := model.pool.Exec(ctx, "INSERT INTO sessions (id, student_id, expires_at) VALUES ($1, $2, $3)", sessionId, studentId, expiration)
	return err
}

func (model *StudentModel) InsertPendingRegister(ctx context.Context, code string, data dto.StudentData) (string, error) {
	var pendingId string
	encodedStudentData, err := json.Marshal(data)
	err = model.pool.QueryRow(ctx, "INSERT INTO pending_registers (code, data, attempts, expires_at) VALUES ($1, $2, $3, NOW() + INTERVAL '10 minutes') RETURNING id", code, encodedStudentData, 0).Scan(&pendingId)
	return pendingId, err
}

func (model *StudentModel) QueryPendingRegister(ctx context.Context, pendingId string) (dto.StudentData, error) {
	var rawData []byte

	err := model.pool.QueryRow(ctx, "SELECT data FROM pending_registers WHERE id=$1 AND expires_at > NOW()", pendingId).Scan(&rawData)
	if err != nil {
		return dto.StudentData{}, err
	}

	var data dto.StudentData
	err = json.Unmarshal(rawData, &data)

	return data, err
}

func (model *StudentModel) QueryCodeByPendingId(ctx context.Context, pendingId string) (string, error) {
	var code string
	err := model.pool.QueryRow(ctx, "SELECT code FROM pending_registers WHERE id=$1 AND expires_at > NOW()", pendingId).Scan(&code)
	return code, err
}

func (model *StudentModel) QueryStudentIdBySession(ctx context.Context, sessionId string) (int64, error) {
	var id int64
	err := model.pool.QueryRow(ctx, "SELECT student_id FROM sessions WHERE id=$1", sessionId).Scan(&id)
	return id, err
}

func (model *StudentModel) QueryStudentIdByEmail(ctx context.Context, email string) (int64, error) {
	var id int64
	err := model.pool.QueryRow(ctx, "SELECT id FROM students WHERE email=$1", email).Scan(&id)
	return id, err
}

func (model *StudentModel) QueryByEmail(ctx context.Context, email string) (data dto.StudentData, err error) {
	data.Email = email
	err = model.pool.QueryRow(ctx, "SELECT name, class, password FROM students WHERE email=$1", email).Scan(&data.Name, &data.Class, &data.Password)
	return
}

func (model *StudentModel) QueryById(ctx context.Context, id int64) (dto.StudentData, error) {
	var data dto.StudentData
	err := model.pool.QueryRow(ctx, "SELECT name, email, class, password FROM students WHERE id=$1", id).Scan(&data.Name, &data.Email, &data.Class, &data.Password)
	return data, err
}

func (model *StudentModel) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := model.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM students WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

// ----------------------------------------- TESTE -----------------------------------------------------

func (model *StudentModel) QueryAll(ctx context.Context) (students []dto.StudentData, err error) {
	rows, err := model.pool.Query(ctx, "SELECT name, email, class, password FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var current dto.StudentData

		err := rows.Scan(&current.Name, &current.Email, &current.Class, &current.Password)
		if err != nil {
			return nil, err
		}

		students = append(students, current)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return
}
