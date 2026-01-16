package note

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/db"
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
	dbClient db.Client
}

func NewRepository(dbClient db.Client) repository.NoteRepository {
	return &repo{
		dbClient: dbClient,
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

	q := db.Query{
		Name:     "create.note",
		QueryRaw: query,
	}

	var id int64
	err = r.dbClient.DB().ScanOneContext(ctx, &id, q, args...)
	if err != nil {
		return 0, err
	}

	return id, nil
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
	q := db.Query{
		Name:     "get.note",
		QueryRaw: query,
	}

	err = r.dbClient.DB().ScanOneContext(ctx, &note, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToNoteFromRepo(&note), nil
}
