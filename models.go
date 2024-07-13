package jsontocql

type Condition struct {
	Column  string     `json:"column" mapstructure:"column"`
	Operand string     `json:"operand" mapstructure:"operand"`
	Data    Resolvable `json:"resolvable" mapstructure:"resolvable"`
}

type Column struct {
	Column string     `json:"column" mapstructure:"column"`
	Data   Resolvable `json:"resolvable" mapstructure:"resolvable"`
}

type Projection struct {
	Column    string   `json:"column" mapstructure:"column"`
	As        string   `json:"as" mapstructure:"as"`
	Mutations []string `json:"mutations" mapstructure:"mutations"`
}

type Resolvable struct {
	ResolveType string                 `json:"resolveType" mapstructure:"resolveType"`
	ResolveData map[string]interface{} `json:"resolveData" mapstructure:"resolveData"`
}
