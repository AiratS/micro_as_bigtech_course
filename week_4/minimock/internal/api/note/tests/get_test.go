package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/mocks"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService
	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	mc := minimock.NewController(t)

	ctx := context.Background()

	id := gofakeit.Int64()
	title := gofakeit.Animal()
	content := gofakeit.Name()
	createdAt := gofakeit.Date()
	updatedAt := gofakeit.Date()

	noteModel := &model.Note{
		ID: id,
		Info: model.NoteInfo{
			Title:   title,
			Content: content,
		},
		CreatedAt: createdAt,
		UpdatedAt: sql.NullTime{
			Time:  updatedAt,
			Valid: true,
		},
	}

	req := &desc.GetRequest{
		Id: id,
	}

	res := &desc.GetResponse{
		Note: &desc.Note{
			Id: id,
			Info: &desc.NoteInfo{
				Title:   title,
				Content: content,
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
		},
	}

	serviceErr := fmt.Errorf("service error")

	tests := []struct {
		name            string
		args            args
		err             error
		want            *desc.GetResponse
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "Success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err:  nil,
			want: res,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := mocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(noteModel, nil)
				return mock
			},
		},
		{
			name: "Error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err:  serviceErr,
			want: nil,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := mocks.NewNoteServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			noteServiceMock := tt.noteServiceMock(mc)

			api := note.NewImplementation(noteServiceMock)
			info, err := api.Get(ctx, req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, info)
		})
	}
}
