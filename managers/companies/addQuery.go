package companies

import (
	"encoding/json"
	"fmt"

	"github.com/nyaruka/phonenumbers"
	"github.com/pretcat/ugc_test_task/models"
)

const (
	CategoryIdsKey = "category_ids"
)

type AddQuery struct {
	TraceId      string   `json:"-"`
	Name         string   `json:"name"`
	BuildingId   string   `json:"building_id"`
	PhoneNumbers []string `json:"phone_numbers"`
	CategoryIds  []string `json:"category_ids"`
}

func NewAddQueryFromJson(data []byte) (query AddQuery, err error) {
	if err = json.Unmarshal(data, &query); err != nil {
		return AddQuery{}, err
	}
	return query, nil
}

func (query AddQuery) Validate() (err error) {
	for _, phoneNumber := range query.PhoneNumbers {
		phoneNumber, err = formatPhoneNumber(phoneNumber)
		if err != nil {
			return err
		}
	}
	if len(query.Name) == 0 {
		return fmt.Errorf("'%s' is empty", models.NameKey)
	}
	if len(query.BuildingId) == 0 {
		return fmt.Errorf("'%s' is empty", models.BuildingIdKey)
	}
	if len(query.PhoneNumbers) == 0 {
		return fmt.Errorf("'%s' is empty", models.PhoneNumbersKey)
	}
	if len(query.CategoryIds) == 0 {
		return fmt.Errorf("'%s' is empty", CategoryIdsKey)
	}
	return nil
}

func formatPhoneNumber(num string) (string, error) {
	pn, err := phonenumbers.Parse(num, "")
	if err != nil {
		return "", fmt.Errorf("parse phone mumber: '%s': %v", num, err)
	}
	reason := phonenumbers.IsPossibleNumberWithReason(pn)
	switch reason {
	case phonenumbers.INVALID_COUNTRY_CODE:
		return "", fmt.Errorf("phone number '%s': invalid country code", num)
	case phonenumbers.INVALID_LENGTH:
		return "", fmt.Errorf("phone number '%s': invalid length", num)
	case phonenumbers.TOO_SHORT:
		return "", fmt.Errorf("phone number '%s' too short", num)
	case phonenumbers.TOO_LONG:
		return "", fmt.Errorf("phone number '%s' too long", num)
	}
	return fmt.Sprintf("+%d%d", *pn.CountryCode, *pn.NationalNumber), nil
}
