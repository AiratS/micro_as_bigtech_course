package note

import (
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/client/db"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
)

type noteService struct {
	noteRepository repository.NoteRepository
	txManager      db.TxManager
}

func NewService(
	noteRepository repository.NoteRepository,
	txManager db.TxManager,
) service.NoteService {
	return &noteService{
		noteRepository: noteRepository,
		txManager:      txManager,
	}
}
