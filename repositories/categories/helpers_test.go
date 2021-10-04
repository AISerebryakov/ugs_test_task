package categories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareSearchByName(t *testing.T) {
	args := PrepareSearchByName(`?±*tea, level_21, Привет, МёдЁ, @()#@#\/[]{}-=+`)
	expectArgs := "tea level_21 привет мёдё"
	assert.Equal(t, expectArgs, args)
}
