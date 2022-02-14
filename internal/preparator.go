package internal

import (
	"amigo/pkg"
	"strings"
)

func PrepareQuery(query string) string {
	switch *pkg.DbProvider {
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
