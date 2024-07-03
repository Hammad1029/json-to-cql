package jsontocql

type Condition struct {
	Column  string     `json:"column"`
	Operand string     `json:"operand"`
	Data    Resolvable `json:"resolvable"`
}

type Column struct {
	Column string     `json:"column"`
	Data   Resolvable `json:"resolvable"`
}

type Projection struct {
	Column    string   `json:"column"`
	As        string   `json:"as"`
	Mutations []string `json:"mutations"`
}

type Resolvable struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
