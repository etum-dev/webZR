package utils

import (
	"errors"
	"testing"
)

func TestFileWrite(t *testing.T) {
	t.Run("jsonwrite", func(t *testing.T) {
		errors.New("asdasda")
	})

	t.Run("Filewrite OOM", func(t *testing.T) {
		// test that simulates full diskspace.
		// TODO: in my program, i should check and warn about low diskspace if target list > 80 hosts or so.
		// TODO: mathematic check on potentially largest file
	})
}
