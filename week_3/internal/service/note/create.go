package note

import (
	"context"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
)

func (s *noteService) Create(ctx context.Context, noteInfo *model.NoteInfo) (int64, error) {
	return s.noteRepository.Create(ctx, noteInfo)
}
