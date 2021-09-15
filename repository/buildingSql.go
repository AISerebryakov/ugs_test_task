package repository

import (
	"fmt"
	"strings"
	"ugc_test_task/models"
)

const (
	buildingsTableName = "buildings"
)

var insertBuildingSQL = ""

func init() {
	insertBuildingSQL = initInsertBuildingSQL()
}

func initInsertBuildingSQL() string {
	rows := []string{models.IdKey, models.AddressKey, models.LocationKey}
	return fmt.Sprintf("insert into %s (%s) values (%s);",
		buildingsTableName,
		strings.Join(rows, ","),
		strings.Join(pgSqlArguments(len(rows)), ","))
}
