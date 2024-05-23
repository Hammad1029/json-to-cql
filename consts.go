package jsontocql

const (
	Select = "SELECT"
	Insert = "INSERT"
	Update = "UPDATE"
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
