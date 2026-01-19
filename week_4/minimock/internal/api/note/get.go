package note

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/converter"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
)

func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	info, err := s.noteService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("Get note data: %v", info)

	return &desc.GetResponse{
		Note: converter.ToDescNoteFromService(info),
	}, nil
}
