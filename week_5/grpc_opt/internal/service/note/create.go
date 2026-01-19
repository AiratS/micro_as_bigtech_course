package note

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
)

func (s *noteService) Create(ctx context.Context, noteInfo *model.NoteInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		id, err := s.noteRepository.Create(ctx, noteInfo)
		if err != nil {
			return err
		}

		note, err := s.noteRepository.Get(ctx, id)
		if err != nil {
			return err
		}

		log.Printf("note data: %v", note)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, err
}
