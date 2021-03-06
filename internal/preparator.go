package internal

import (
	"strings"
)

func PrepareQuery(query string) string {
	switch GetDBProvider() {
	case "mariadb":
		return query
	case "mysql":
		return query
	case "postgres":
		return preparePostgres(query)
	default:
		return query
	}
}

func preparePostgres(query string) string {
	return strings.Replace(query, "?", "$1", 1)
}
