package repository

import (
	"context"
	"testing"
	"time"
	"ugc_test_task/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertCategory(t *testing.T) {
	category := models.Category{
		Id:       uuid.NewString(),
		Name:     "Top.Transport.Electro.Bus",
		CreateAt: time.Now().UnixNano() / 1e6,
	}
	err := repos.InsertCategory(context.Background(), category)
	assert.NoError(t, err, "Insert category to repository")

	savedCategory, found, err := repos.FetchCategoryById(context.Background(), category.Id)
	assert.NoError(t, err, "Fetch category from repository")
	assert.Equal(t, true, found, "Category not found")
	assert.Equal(t, category, savedCategory)

	err = repos.DeleteCategoryById(context.Background(), savedCategory.Id)
	assert.NoError(t, err, "Delete category from repository")
}

func TestRepository_FetchByNames(t *testing.T) {
	err := repos.fetchCategoryIdsByNames(context.Background(), []string{"Top.Transport.Moto", "Top.Transport.Bus"}, func(id, name string) error {
		return nil
	})
	assert.NoError(t, err, "Fetch categories by names")
}
