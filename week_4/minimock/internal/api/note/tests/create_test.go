package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/note"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/model"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service"
	"github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/mocks"
	desc "github.com/AiratS/micro_as_bigtech_course/week_3/pkg/note_v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type noteServiceMockFunc func(mc *minimock.Controller) service.NoteService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	ctx := context.Background()
	mc := minimock.NewController(t)

	id := gofakeit.Int64()
	title := gofakeit.City()
	content := gofakeit.Animal()

	req := &desc.CreateRequest{
		Info: &desc.NoteInfo{
			Title:   title,
			Content: content,
		},
	}

	res := &desc.CreateResponse{
		Id: id,
	}

	info := &model.NoteInfo{
		Title:   title,
		Content: content,
	}

	serviceErr := fmt.Errorf("service error")

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "Success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := mocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "Error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := mocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock := tt.noteServiceMock(mc)
			api := note.NewImplementation(mock)
			newID, err := api.Create(ctx, req)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
