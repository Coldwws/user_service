package prettier

import (
	"fmt"
	"strings"
)

const (
	PlaceholderDollar = "$"
)

func Pretty(query string, placeholder string, args ...any) string {
	for i, param := range args {
		var value string

		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}

		query = strings.Replace(
			query,
			fmt.Sprintf("%s%d", placeholder, i+1),
			value,
			1,
		)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.TrimSpace(query)
}