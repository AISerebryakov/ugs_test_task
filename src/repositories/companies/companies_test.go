package companies

import (
	"context"
	"fmt"
	"testing"
	"ugc_test_task/src/common"
	"ugc_test_task/src/models"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertCompany(t *testing.T) {
	fmt.Println("RUN TEST")
	repos, err := testRepos()
	assert.NoError(t, err, "Get repository")
	repos.Stop()

	comp := models.Company{
		Id:           uuid.NewString(),
		Name:         "Test_Firm",
		CreateAt:     common.NewTimestamp(),
		BuildingId:   uuid.NewString(),
		Address:      "Test address 2",
		PhoneNumbers: []string{"+76472834883"},
		Categories:   []string{"Level_11.Level_21.Level_31", "Level_11.Level_21.Level_32"},
	}
	err = repos.Insert(context.Background(), comp)
	assert.NoError(t, err, "Insert comp to repository")

	savedComp, found, err := repos.Select(context.Background()).ById(comp.Id).One()
	assert.NoError(t, err, "Fetch comp from repository")
	assert.Equal(t, true, found, "Company not found")
	assert.Equal(t, comp, savedComp)

	//err = repos.DeleteCompanyById(context.Background(), savedComp.Id)
	//assert.NoError(t, err, "Delete comp from repository")
}

func TestRepository_FetchCompaniesForCategories(t *testing.T) {
	repos, err := testRepos()
	assert.NoError(t, err, "Get repository")

	newComp := models.Company{
		Id:           uuid.NewString(),
		Name:         "Test_Firm",
		CreateAt:     common.NewTimestamp(),
		BuildingId:   uuid.NewString(),
		Address:      "Test address",
		PhoneNumbers: []string{"+76456734235"},
		Categories:   []string{"Top.Transport.Moto", "Top.Transport.Cars"},
	}
	err = repos.Insert(context.Background(), newComp)
	assert.NoError(t, err, "Insert comp to repository")

	found := false
	err = repos.Select(context.Background()).ForCategories(newComp.Categories).Iter(func(company models.Company) error {
		t.Log(company)
		if company.Id == newComp.Id {
			found = true
			assert.Equal(t, newComp, company)
		}
		return nil
	})
	assert.NoError(t, err, "Fetch companies for categories")
	assert.Equal(t, true, found, "Company not found")

	err = repos.DeleteCompanyById(context.Background(), newComp.Id)
	assert.NoError(t, err, "Delete comp from repository")
}
