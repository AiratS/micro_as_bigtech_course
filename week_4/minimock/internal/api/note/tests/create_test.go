package tests

import (
	"context"
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

	req := &desc.CreateRequest{
		Info: &desc.NoteInfo{
			Title:    gofakeit.City(),
			Content:  gofakeit.City(),
			Author:   gofakeit.Name(),
			IsPublic: gofakeit.Bool(),
		},
	}

	res := &desc.CreateResponse{
		Id: id,
	}

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		noteServiceMock noteServiceMockFunc
	}{
		{
			name: "No error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			noteServiceMock: func(mc *minimock.Controller) service.NoteService {
				mock := mocks.NewNoteServiceMock(mc)
				mock.CreateMock.Expect(ctx, &model.NoteInfo{}).Return(id, nil)
				return mock
			},
		},
	}

	for _, test := range tests {
		t.Run("First", func(t *testing.T) {
			mock := test.noteServiceMock(mc)
			api := note.NewImplementation(mock)
			api.Create(ctx, req)
			require.Equal(t, res, req)
		})
	}
}
