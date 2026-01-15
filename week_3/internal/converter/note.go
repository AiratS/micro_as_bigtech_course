package converter

import (
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteInfoFromDesc(descNoteInfo *desc.NoteInfo) *model.NoteInfo {
	return &model.NoteInfo{
		Title:   descNoteInfo.Title,
		Content: descNoteInfo.Content,
	}
}

func ToDescNoteFromService(serviceNote *model.Note) *desc.Note {
	var updatedAt *timestamppb.Timestamp
	if serviceNote.UpdatedAt.Valid {
		updatedAt = timestamppb.New(serviceNote.UpdatedAt.Time)
	}

	return &desc.Note{
		Id:        serviceNote.ID,
		Info:      ToDescNoteInfoFromService(&serviceNote.Info),
		CreatedAt: timestamppb.New(serviceNote.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToDescNoteInfoFromService(serviceNoteInfo *model.NoteInfo) *desc.NoteInfo {
	return &desc.NoteInfo{
		Title:   serviceNoteInfo.Title,
		Content: serviceNoteInfo.Content,
	}
}
