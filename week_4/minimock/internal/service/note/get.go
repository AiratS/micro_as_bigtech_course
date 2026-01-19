package note

import (
	"context"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
)

func (s *noteService) Get(ctx context.Context, id int64) (*model.Note, error) {
	return s.noteRepository.Get(ctx, id)
}
