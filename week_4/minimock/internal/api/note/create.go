package note

import (
	"context"
	"log"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/converter"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.noteService.Create(ctx, converter.ToNoteInfoFromDesc(req.Info))
	if err != nil {
		return nil, err
	}

	log.Printf("Note created! %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
