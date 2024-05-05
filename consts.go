package jsontocql

const (
	Select = "SELECT"
	Update = "INSERT"
	Insert = "UPDATE"
	Delete = "DELETE"
)

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
