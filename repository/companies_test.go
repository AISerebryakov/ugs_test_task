package repository

import (
	"context"
	"testing"
	"ugc_test_task/models"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertCompany(t *testing.T) {
	comp := models.Company{
		Id:           uuid.NewString(),
		Name:         "Test_Firm",
		BuildingId:   uuid.NewString(),
		Address:      "Test address",
		PhoneNumbers: []string{"+76456734235"},
		Categories:   []string{"Top.Transport.Moto", "Top.Transport.Cars"},
	}
	err := repos.InsertCompany(context.Background(), comp)
	assert.NoError(t, err, "Insert comp to repository")

	savedComp, found, err := repos.FetchCompanyById(context.Background(), comp.Id)
	assert.NoError(t, err, "Fetch comp from repository")
	assert.Equal(t, true, found, "Company not found")
	assert.Equal(t, comp, savedComp)

	err = repos.DeleteCompanyById(context.Background(), savedComp.Id)
	assert.NoError(t, err, "Delete comp from repository")
}

func TestRepository_FetchCompaniesForCategories(t *testing.T) {
	newComp := models.Company{
		Id:           uuid.NewString(),
		Name:         "Test_Firm",
		BuildingId:   uuid.NewString(),
		Address:      "Test address",
		PhoneNumbers: []string{"+76456734235"},
		Categories:   []string{"Top.Transport.Moto", "Top.Transport.Cars"},
	}
	err := repos.InsertCompany(context.Background(), newComp)
	assert.NoError(t, err, "Insert comp to repository")

	found := false
	err = repos.FetchCompaniesForCategories(context.Background(), newComp.Categories, func(comp models.Company) error {
		t.Log(comp)
		if comp.Id == newComp.Id {
			found = true
			assert.Equal(t, newComp, comp)
		}
		return nil
	})
	assert.NoError(t, err, "Fetch companies for categories")

	if !found {
		t.Error("Comp not found")
	}

	err = repos.DeleteCompanyById(context.Background(), newComp.Id)
	assert.NoError(t, err, "Delete comp from repository")
}
