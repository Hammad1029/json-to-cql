package jsontocql

var queryTypes = map[string]string{
	"select": "SELECT",
	"insert": "INSERT",
	"update": "UPDATE",
	"delete": "DELETE",
}

const asNotAllowed = "null"

var allowedOperators = map[string]string{
	"eq":     "=",
	"ne":     "!=",
	"gt":     ">",
	"gte":    ">=",
	"lt":     "<",
	"lte":    "<=",
	"in":     "IN",
	"nin":    "NOT IN",
	"con":    "~",
	"notcon": "!~",
}
