package repository

import (
	"fmt"
	"strings"
	"ugc_test_task/models"
)

const (
	firmsTableName = "firms"
)

var insertFirmSQL = ""

func init() {
	insertFirmSQL = initInsertFirmSQL()
}

func initInsertFirmSQL() string {
	rows := []string{models.IdKey, models.NameKey, models.BuildingIdKey, models.PhoneNumbersKey}
	return fmt.Sprintf("insert into %s (%s) values (%s);",
		firmsTableName,
		strings.Join(rows, ","),
		strings.Join(pgSqlArguments(len(rows)), ","))
}

func initCreateFirmsTableSQL() string {
	return fmt.Sprintf("")
}
