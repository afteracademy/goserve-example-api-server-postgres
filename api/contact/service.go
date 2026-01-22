package contact

import (
	"context"

	"github.com/afteracademy/goserve-example-api-server-postgres/api/contact/dto"
	"github.com/afteracademy/goserve-example-api-server-postgres/api/contact/model"
	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/afteracademy/goserve/v2/postgres"
	"github.com/google/uuid"
)

type Service interface {
	CreateMessage(d *dto.MessageCreate) (*model.Message, error)
	FetchMessage(id uuid.UUID) (*model.Message, error)
	FetchPaginatedMessage(p *coredto.Pagination) ([]*model.Message, error)
}

type service struct {
	db postgres.Database
}

func NewService(db postgres.Database) Service {
	return &service{
		db: db,
	}
}

func (s *service) CreateMessage(
	dto *dto.MessageCreate,
) (*model.Message, error) {
	ctx := context.Background()
	msg := model.Message{}

	query := `
		INSERT INTO messages (
			type,
			msg
		)
		VALUES ($1, $2)
		RETURNING
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
	`

	err := s.db.Pool().QueryRow(
		ctx,
		query,
		dto.Type,
		dto.Msg,
	).Scan(
		&msg.ID,
		&msg.Type,
		&msg.Msg,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *service) FetchMessage(
	id uuid.UUID,
) (*model.Message, error) {
	ctx := context.Background()
	query := `
		SELECT
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
		FROM messages
		WHERE id = $1
	`

	var m model.Message

	err := s.db.Pool().QueryRow(ctx, query, id).
		Scan(
			&m.ID,
			&m.Type,
			&m.Msg,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *service) FetchPaginatedMessage(
	p *coredto.Pagination,
) ([]*model.Message, error) {
	ctx := context.Background()
	offset := (p.Page - 1) * p.Limit

	query := `
		SELECT
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
		FROM messages
		WHERE status = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Pool().Query(ctx, query, p.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message

	for rows.Next() {
		var m model.Message
		if err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Msg,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
