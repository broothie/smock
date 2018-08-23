package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const content = "content"

func TestBorderify_SingleLine(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("| %s |\n", content), Borderify(content, len(content)))
	})

	t.Run("with whitespace", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("| %s  |\n", content), Borderify(content, len(content)+1))
	})

	t.Run("with buckling", func(t *testing.T) {
		assert.Equal(t, "| conten\n  t      |\n", Borderify(content, len(content)-1))
	})
}

const multiLineContent = `content
on
several
lines`

func TestBorderify_MultiLine(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		assert.Equal(t, "| content |\n| on      |\n| several |\n| lines   |\n", Borderify(multiLineContent, len("content")))
	})

	t.Run("with buckling", func(t *testing.T) {
		assert.Equal(t, "| conten\n  t      |\n| on     |\n| severa\n  l      |\n| lines  |\n", Borderify(multiLineContent, len("conten")))
	})
}

const entrified = `+------+
| head |
+------+
| cont
  ent  |
| on   |
| seve
  ral  |
| line
  s    |
+------+
`

func TestEntrify(t *testing.T) {
	assert.Equal(t, entrified, Entrify(len("head"), "head", multiLineContent))
}
