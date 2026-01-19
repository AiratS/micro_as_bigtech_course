package converter

import (
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	repoModel "github.com/AiratS/micro_as_bigtech_course/week_3/internal/repository/note/model"
)

func ToNoteFromRepo(repoNote *repoModel.Note) *model.Note {
	return &model.Note{
		ID:        repoNote.ID,
		Info:      ToNoteInfoFromRepo(repoNote.Info),
		CreatedAt: repoNote.CreatedAt,
		UpdatedAt: repoNote.UpdatedAt,
	}
}

func ToNoteInfoFromRepo(repoNoteInfo repoModel.NoteInfo) model.NoteInfo {
	return model.NoteInfo{
		Title:   repoNoteInfo.Title,
		Content: repoNoteInfo.Content,
	}
}
