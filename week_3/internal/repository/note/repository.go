package note

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note/converter"
	repoModel "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note/model"
)

const (
	tableName = "note"

	idColumn        = "id"
	titleColumn     = "title"
	contentColumn   = "body"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.NoteRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, info *model.NoteInfo) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(titleColumn, contentColumn).
		Values(info.Title, info.Content).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type noteType struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Note, error) {
	builderSelect := sq.Select(idColumn, titleColumn, contentColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	var note repoModel.Note
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&note.ID, &note.Info.Title, &note.Info.Content, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
