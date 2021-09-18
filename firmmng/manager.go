package firmmng

import "ugc_test_task/models"

type Manager struct {
	conf Config
}

func New(conf Config) (m Manager) {
	m.conf = conf
	return m
}

//todo: normalize of phone numbers

func (m Manager) GetFirms(query GetQuery, callback func(firm models.Company) error) {
	callback(models.Company{Name: "Test firm"})
}
