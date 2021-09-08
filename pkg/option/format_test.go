package option_test

import (
	"strings"
	"testing"

	"github.com/maargenton/go-testpredicate/pkg/require"
	"github.com/maargenton/go-testpredicate/pkg/subexpr"

	"github.com/maargenton/go-cli/pkg/option"
)

func lines(s string) []string {
	var lines = strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func TestFormatOptionDescription(t *testing.T) {

	t.Run("Given a list of argument usage", func(t *testing.T) {
		var usage = []option.Description{
			{"short", "short description"},
			{"long-long-long-long", "long description"},
			{"no-desc", ""},
		}

		t.Run("when calling FormatOptionDescription()", func(t *testing.T) {
			var lines = lines(option.FormatOptionDescription("__", 80, usage))

			t.Run("then all descriptions are aligned", func(t *testing.T) {
				require.That(t, lines[0]).StartsWith("__short               :")
				require.That(t, lines[1]).StartsWith("__long-long-long-long :")
				require.That(t, lines[2]).StartsWith("__no-desc")
			})
		})
	})

	t.Run("Given a list of argument usage with long descriptions", func(t *testing.T) {
		var usage = []option.Description{
			{
				Option:      "-p, --port <value>",
				Description: "somewhat long-ish description, just long enough to be wrapped",
			},
		}

		t.Run("when calling FormatOptionDescription()", func(t *testing.T) {
			var lines = lines(option.FormatOptionDescription("__", 50, usage))

			t.Run("then descriptions is wrapped and aligned accros multiple lines", func(t *testing.T) {
				require.That(t, lines[0]).StartsWith("__-p, --port <value> : ")
				require.That(t, lines[1:]).All(
					subexpr.Value().StartsWith("                       "))
			})
		})
	})
}

func TestFormatCompletion(t *testing.T) {

	t.Run("Given a list of completion suggestions", func(t *testing.T) {
		var suggestions = []option.Description{
			{"short", "short description"},
			{"long-long-long-long", "long long long long description"},
			{"no-desc", ""},
		}

		t.Run("when calling FormatCompletion()", func(t *testing.T) {
			l := lines(option.FormatCompletion(80, suggestions))
			t.Run("then suggestions are formatted one per line", func(t *testing.T) {
				require.That(t, l).Length().Eq(3)
			})
		})
		t.Run("when calling FormatCompletion() with narrow width", func(t *testing.T) {
			l := lines(option.FormatCompletion(40, suggestions))
			t.Run("then long descriptions are truncated", func(t *testing.T) {
				require.That(t, l[1]).EndsWith("...")
			})
		})
		t.Run("when calling FormatCompletion() with extremely narrow width", func(t *testing.T) {
			l := lines(option.FormatCompletion(20, suggestions))
			t.Run("then descriptions are dropped", func(t *testing.T) {
				require.That(t, l[1]).EndsWith("long-long")
			})
		})
	})
}
