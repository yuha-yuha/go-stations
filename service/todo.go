package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	var createdAT, UpdatedAT time.Time

	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, err
	}

	sqlResult := s.db.QueryRowContext(ctx, confirm, id)
	var todo model.TODO
	sqlResult.Scan(&subject, &description, &createdAT, &UpdatedAT)
	todo.ID = id
	todo.Subject = subject
	todo.Description = description
	todo.CreatedAt = createdAT
	todo.UpdatedAt = UpdatedAT

	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	todos := []*model.TODO{}

	var query string
	var rows *sql.Rows
	var err error

	if prevID == 0 {
		query = read
		rows, err = s.db.QueryContext(ctx, query, size)

		if err != nil {
			return nil, err
		}

	} else {
		query = readWithID
		rows, err = s.db.QueryContext(ctx, query, prevID, size)

		if err != nil {
			return nil, err
		}

	}
	for rows.Next() {
		var created_at time.Time
		var updated_at time.Time
		var id int64
		var subject, description string

		err = rows.Scan(&id, &subject, &description, &created_at, &updated_at)

		todo := model.TODO{
			ID:          id,
			Subject:     subject,
			Description: description,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
		}
		log.Print("sdfsdfsdfsdfsdfsf")
		todos = append(todos, &todo)
		if err != nil {
			return nil, err
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	var created_at, updated_at time.Time

	if id == 0 {
		return nil, &model.ErrNotFound{}
	}

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(subject, description, id)

	if err != nil {
		return nil, err
	}

	row := s.db.QueryRowContext(ctx, confirm, id)

	row.Scan(&subject, &description, &created_at, &updated_at)

	todo := model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   created_at,
		UpdatedAt:   updated_at,
	}

	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return nil
	}

	stmt, err := s.db.PrepareContext(ctx, fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1)))
	var idInterfaces []interface{}

	for _, id := range ids {
		idInterfaces = append(idInterfaces, id)
	}

	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, idInterfaces...)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
