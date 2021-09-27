package categories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamesToLtreeArgs(t *testing.T) {
	names := []string{"level_1.level_2", "", "level_1.level_3"}
	expectedArgs := "level_1.level_2*@|level_1.level_3*@"
	args := NamesToLtreeArgs(names)
	assert.Equal(t, expectedArgs, args)
}

func TestPrepareSearchByName(t *testing.T) {
	args := PrepareSearchByName(`?±*tea, level_21, Привет, МёдЁ, @()#@#\/[]{}-=+`)
	expectArgs := "tea*@|level_21*@|Привет*@|МёдЁ*@"
	assert.Equal(t, expectArgs, args)
}
