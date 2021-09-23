package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"ugc_test_task/companymng"
	"ugc_test_task/models"

	"github.com/gojuno/minimock/v3"

	"github.com/stretchr/testify/assert"
)

func TestHandler_GetCompanies(t *testing.T) {
	ctrl := minimock.NewController(t)
	companyMng := NewCompanyManagerMock(ctrl).
		GetCompaniesMock.Inspect(func(query companymng.GetQuery, clb func(firm models.Company) error) {
		t.Log("Query: ", query)
		clb(models.Company{Name: "Test firm 34"})

	}).Return(nil).
		AddCompanyMock.Expect(companymng.AddQuery{}).Return(models.Company{}, nil)

	api, err := NewApi(Config{CompanyManager: companyMng})
	assert.NoError(t, err)

	srv := httptest.NewServer(handler{api})
	defer srv.Close()

	res, err := http.Get(srv.URL + companiesPath)
	assert.NoError(t, err)
	if res.StatusCode != http.StatusOK {
		//todo: change checker
		t.Errorf("status not OK")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	t.Log("Body: ", string(body))
}
