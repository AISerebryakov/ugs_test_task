package companyrepos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoriesToLtreeArgs(t *testing.T) {
	categories := []string{"level_1.level_2", "", "level_1.level_3"}
	expectedArgs := "level_1.level_2*@|level_1.level_3*@"
	args := categoriesToLtreeArgs(categories)
	assert.Equal(t, expectedArgs, args)
}

func BenchmarkCategoriesToLtreeArgs(b *testing.B) {
	categories := []string{"level_1.level_2", "", "level_1.level_3"}
	for i := 0; i < b.N; i++ {
		args := categoriesToLtreeArgs(categories)
		_ = len(args)
	}
}
