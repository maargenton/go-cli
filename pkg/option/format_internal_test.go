package option

import (
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/subexpr"
	"github.com/maargenton/go-testpredicate/pkg/verify"
)

func TestLineWrap(t *testing.T) {
	t.Run("Given a piece of text wider than allowed", func(t *testing.T) {
		text := "generate a bash script used to setup completion for this command; " +
			"to use, add the following line to your .bash_profile or equivalent: " +
			"`eval $(" + "<cmd>" + " --bash-completion-script)`"

		t.Run("when calling lineWrap()", func(t *testing.T) {
			width := 40
			lines := lineWrap(text, width)

			t.Run("then all lines are at most the width", func(t *testing.T) {
				verify.That(t, lines).Length().Eq(5)
				verify.That(t, lines).All(
					subexpr.Value().Length().Le(width))
			})
		})
		t.Run("when calling lineWrap() with narrow width", func(t *testing.T) {
			width := 20
			lines := lineWrap(text, width)

			t.Run("then lines longer than width have no space", func(t *testing.T) {
				verify.That(t, lines).Length().Eq(10)
				verify.That(t, lines[:9]).All(
					subexpr.Value().Length().Le(width))
				verify.That(t, lines[9]).Matches(`^[^\s]+$`)
			})
		})
	})

	t.Run("Given a piece of text with line breaks", func(t *testing.T) {
		text := "generate a bash script used to setup completion for this command; " +
			"to use, add the following line to your .bash_profile or equivalent:\n" +
			"`eval $(" + "<cmd>" + " --bash-completion-script)`"

		t.Run("when calling lineWrap()", func(t *testing.T) {
			width := 50
			lines := lineWrap(text, width)

			t.Run("then line breaks are preserved", func(t *testing.T) {
				verify.That(t, lines).Length().Eq(4)
				verify.That(t, lines).All(
					subexpr.Value().Length().Le(width))
				verify.That(t, lines).Any(
					subexpr.Value().StartsWith("`eval $("))
			})
		})
	})
}
