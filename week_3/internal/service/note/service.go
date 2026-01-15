package note

import (
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
)

type noteService struct {
	noteRepository repository.NoteRepository
}

func NewService(
	noteRepository repository.NoteRepository,
) service.NoteService {
	return &noteService{
		noteRepository: noteRepository,
	}
}
