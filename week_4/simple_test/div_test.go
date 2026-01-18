package main

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestDiv(t *testing.T) {
	a := gofakeit.Float64()
	b := gofakeit.Float64()

	res := Div(a, b)

	t.Run("test 1", func(t *testing.T) {
		require.Equal(t, a/b, res)
	})

	t.Run("test 2", func(t *testing.T) {
		require.Equal(t, 1.0/3.0, 0.5)
	})
}
