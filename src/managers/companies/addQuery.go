package companies

import (
	"encoding/json"
	"fmt"
	"ugc_test_task/src/models"

	"github.com/nyaruka/phonenumbers"
)

type AddQuery struct {
	Name         string   `json:"name"`
	BuildingId   string   `json:"building_id"`
	Address      string   `json:"address"`
	PhoneNumbers []string `json:"phone_numbers"`
	Categories   []string `json:"categories"`
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
	if len(query.Address) == 0 {
		return fmt.Errorf("'%s' is empty", models.AddressKey)
	}
	if len(query.PhoneNumbers) == 0 {
		return fmt.Errorf("'%s' is empty", models.PhoneNumbersKey)
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
